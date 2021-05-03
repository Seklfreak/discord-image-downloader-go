package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/Jeffail/gabs"
	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/go-version"
	"golang.org/x/net/context"
	"golang.org/x/net/html"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"gopkg.in/ini.v1"
)

var (
	ChannelWhitelist                 map[string]string
	InteractiveChannelWhitelist      map[string]string
	dg                               *discordgo.Session
	DownloadTistorySites             bool
	interactiveChannelLinkTemp       map[string]string
	DiscordUserId                    string
	myDB                             *db.DB
	historyCommandActive             map[string]string
	MaxDownloadRetries               int
	flickrApiKey                     string
	twitterClient                    *anaconda.TwitterApi
	DownloadTimeout                  int
	SendNoticesToInteractiveChannels bool
	StatusEnabled                    bool
	StatusType                       string
	StatusLabel                      discordgo.GameType
	StatusSuffix                     string
	clientCredentialsJson            string
	DriveService                     *drive.Service
	RegexpFilename                   *regexp.Regexp
)

type GfycatObject struct {
	GfyItem struct {
		Mp4URL string `json:"mp4Url"`
	} `json:"gfyItem"`
}

type ImgurAlbumObject struct {
	Data []struct {
		Link string
	}
}

func main() {
	fmt.Printf("> discord-image-downloader-go v%s -- discordgo v%s\n", VERSION, discordgo.VERSION)
	if !isLatestRelease() {
		fmt.Printf("update available on %s !\n", RELEASE_URL)
	}

	var err error

	var configFile string
	flag.StringVar(&configFile, "config", DEFAULT_CONFIG_FILE, "config file to read from")

	flag.Parse()

	if configFile == "" {
		configFile = DEFAULT_CONFIG_FILE
	}

	fmt.Printf("reading config from %s\n", configFile)
	cfg, err := ini.Load(configFile)
	if err != nil {
		fmt.Println("unable to read config file", err)
		cfg = ini.Empty()
	}

	if (!cfg.Section("auth").HasKey("email") ||
		!cfg.Section("auth").HasKey("password")) &&
		!cfg.Section("auth").HasKey("token") {
		cfg.Section("auth").NewKey("email", "your@email.com")
		cfg.Section("auth").NewKey("password", "your password")
		cfg.Section("general").NewKey("skip edits", "true")
		cfg.Section("general").NewKey("download tistory sites", "false")
		cfg.Section("general").NewKey("max download retries", "5")
		cfg.Section("general").NewKey("download timeout", "60")
		cfg.Section("general").NewKey("send notices to interactive channels", "false")
		cfg.Section("status").NewKey("status enabled", "true")
		cfg.Section("status").NewKey("status type", "online")
		cfg.Section("status").NewKey("status label", fmt.Sprint(discordgo.GameTypeWatching))
		cfg.Section("status").NewKey("status suffix", "downloaded pictures")
		cfg.Section("channels").NewKey("channelid1", "C:\\full\\path\\1")
		cfg.Section("channels").NewKey("channelid2", "C:\\full\\path\\2")
		cfg.Section("channels").NewKey("channelid3", "C:\\full\\path\\3")
		cfg.Section("flickr").NewKey("api key", "your flickr api key")
		cfg.Section("twitter").NewKey("consumer key", "your consumer key")
		cfg.Section("twitter").NewKey("consumer secret", "your consumer secret")
		cfg.Section("twitter").NewKey("access token", "your access token")
		cfg.Section("twitter").NewKey("access token secret", "your access token secret")
		err = cfg.SaveTo("config.ini")

		if err != nil {
			fmt.Println("unable to write config file", err)
			return
		}
		fmt.Println("Wrote config file, please fill out and restart the program")
		return
	}

	myDB, err = db.OpenDB(DATABASE_DIR)
	if err != nil {
		fmt.Println("unable to create db", err)
		return
	}
	if myDB.Use("Downloads") == nil {
		fmt.Println("Creating new database...")
		if err := myDB.Create("Downloads"); err != nil {
			fmt.Println("unable to create db", err)
			return
		}
		if err := myDB.Use("Downloads").Index([]string{"Url"}); err != nil {
			fmt.Println("unable to create index", err)
			return
		}
	}

	ChannelWhitelist = cfg.Section("channels").KeysHash()
	InteractiveChannelWhitelist = cfg.Section("interactive channels").KeysHash()
	interactiveChannelLinkTemp = make(map[string]string)
	historyCommandActive = make(map[string]string)
	flickrApiKey = cfg.Section("flickr").Key("api key").MustString("yourflickrapikey")
	twitterConsumerKey := cfg.Section("twitter").Key("consumer key").MustString("your consumer key")
	twitterConsumerSecret := cfg.Section("twitter").Key("consumer secret").MustString("your consumer secret")
	twitterAccessToken := cfg.Section("twitter").Key("access token").MustString("your access token")
	twitterAccessTokenSecret := cfg.Section("twitter").Key("access token secret").MustString("your access token secret")
	if twitterAccessToken != "" &&
		twitterAccessTokenSecret != "" &&
		twitterConsumerKey != "" &&
		twitterConsumerSecret != "" {
		twitterClient = anaconda.NewTwitterApiWithCredentials(
			twitterAccessToken,
			twitterAccessTokenSecret,
			twitterConsumerKey,
			twitterConsumerSecret,
		)
	}

	err = initRegex()
	if err != nil {
		fmt.Println("error initialising regex,", err.Error())
		return
	}

	if cfg.Section("auth").HasKey("token") {
		dg, err = discordgo.New(cfg.Section("auth").Key("token").String())
	} else {
		dg, err = discordgo.New(
			cfg.Section("auth").Key("email").String(),
			cfg.Section("auth").Key("password").String())
	}
	if err != nil {
		// Newer discordgo throws this error for some reason with Email/Password login
		if err.Error() != "Unable to fetch discord authentication token. <nil>" {
			fmt.Println("error creating Discord session,", err)
			return
		}
	}

	dg.AddHandler(messageCreate)

	if cfg.Section("general").HasKey("skip edits") {
		if cfg.Section("general").Key("skip edits").MustBool() == false {
			dg.AddHandler(messageUpdate)
		}
	}

	DownloadTistorySites = cfg.Section("general").Key("download tistory sites").MustBool()
	MaxDownloadRetries = cfg.Section("general").Key("max download retries").MustInt(3)
	DownloadTimeout = cfg.Section("general").Key("download timeout").MustInt(60)
	SendNoticesToInteractiveChannels = cfg.Section("general").Key("send notices to interactive channels").MustBool(false)

	StatusEnabled = cfg.Section("status").Key("status enabled").MustBool(true)
	StatusType = cfg.Section("status").Key("status type").MustString("online")
	StatusLabel = discordgo.GameType(cfg.Section("status").Key("status label").MustInt(int(discordgo.GameTypeWatching)))
	StatusSuffix = cfg.Section("status").Key("status suffix").MustString("downloaded pictures")

	// setup google drive client
	clientCredentialsJson = cfg.Section("google").Key("client credentials json").MustString("")
	if clientCredentialsJson != "" {
		ctx := context.Background()
		authJson, err := ioutil.ReadFile(clientCredentialsJson)
		if err != nil {
			fmt.Println("error opening google credentials json,", err)
		} else {
			config, err := google.JWTConfigFromJSON(authJson, drive.DriveReadonlyScope)
			if err != nil {
				fmt.Println("error parsing google credentials json,", err)
			} else {
				client := config.Client(ctx)
				DriveService, err = drive.New(client)
				if err != nil {
					fmt.Println("error setting up google drive client,", err)
				}
			}
		}
	}
	dg.LogLevel = -1 // to ignore dumb wsapi error
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	dg.LogLevel = 0 // reset

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	fmt.Printf("Client is now connected as %s. Press CTRL-C to exit.\n",
		u.Username)
	DiscordUserId = u.ID

	updateDiscordStatus()

	// keep program running until CTRL-C is pressed.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Closing database...")
	myDB.Close()
	fmt.Println("Logging out of Discord...")
	dg.Close()
	fmt.Println("Exiting...")
	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	handleDiscordMessage(m.Message)
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.EditedTimestamp != discordgo.Timestamp("") {
		handleDiscordMessage(m.Message)
	}
}

func skipDuplicateLinks(linkList map[string]string, channelID string, interactive bool) map[string]string {
	if interactive == false {
		newList := make(map[string]string, 0)
		for link, filename := range linkList {
			downloadedImages := findDownloadedImageByUrl(link)
			isMatched := false
			for _, downloadedImage := range downloadedImages {
				if downloadedImage.ChannelId == channelID {
					isMatched = true
				}
			}
			if isMatched == false {
				newList[link] = filename
			} else {
				fmt.Println("url already downloaded in this channel:", link)
			}
		}
		return newList
	}
	return linkList
}

func handleDiscordMessage(m *discordgo.Message) {
	// If message content is empty (likely due to userbot/selfbot)
	if m.Content == "" && len(m.Attachments) == 0 {
		nms, err := dg.ChannelMessages(m.ChannelID, 10, "", "", "")
		if err == nil {
			if len(nms) > 0 {
				for _, nm := range nms {
					if nm.ID == m.ID {
						m = nm
					}
				}
			}
		}
	}

	if folderName, ok := ChannelWhitelist[m.ChannelID]; ok {
		// download from whitelisted channels
		downloadItems := getDownloadItemsOfMessage(m)

		for _, downloadItem := range downloadItems {
			startDownload(
				downloadItem.Link,
				downloadItem.Filename,
				folderName,
				m.ChannelID,
				m.Author.ID,
				downloadItem.Time,
			)
		}

	} else if _, ok := InteractiveChannelWhitelist[m.ChannelID]; ok {
		// handle interactive channel

		// skip messages from the Bot
		if m.Author == nil || m.Author.ID == DiscordUserId {
			return
		}

		dg.ChannelTyping(m.ChannelID)

		args := strings.Fields(m.Content)
		if len(args) <= 0 {
			return
		}

		_, historyCommandIsActive := historyCommandActive[m.ChannelID]

		switch strings.ToLower(args[0]) {
		case "help":
			helpHandler(m)
		case "version":
			versionHandler(m)
		case "channels":
			channelsHandler(m)
		case "stats":
			statsHandler(m)
		case "history":
			historyHandler(m)
		default:
			if historyCommandIsActive {
				historyHandler(m)
				return
			}

			defaultHandler(m)
		}
	}
}

type GithubReleaseApiObject struct {
	TagName string `json:"tag_name"`
}

func isLatestRelease() bool {
	githubReleaseApiObject := new(GithubReleaseApiObject)
	getJson(RELEASE_API_URL, githubReleaseApiObject)
	currentVer, err := version.NewVersion(VERSION)
	if err != nil {
		fmt.Println(err)
		return true
	}
	lastVer, err := version.NewVersion(githubReleaseApiObject.TagName)
	if err != nil {
		fmt.Println(err)
		return true
	}
	if lastVer.GreaterThan(currentVer) {
		return false
	}
	return true
}

// http://stackoverflow.com/a/35240286/1443726
func IsValid(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}

func getTwitterUrls(url string) (map[string]string, error) {
	parts := strings.Split(url, ":")
	if len(parts) < 2 {
		return nil, errors.New("unable to parse twitter url")
	}
	return map[string]string{"https:" + parts[1] + ":orig": filenameFromUrl(parts[1])}, nil
}

func getTwitterParamsUrls(url string) (map[string]string, error) {
	matches := RegexpUrlTwitterParams.FindStringSubmatch(url)

	return map[string]string{
		"https://pbs.twimg.com/media/" + matches[3] + "." + matches[4] + ":orig": matches[3] + "." + matches[4],
	}, nil
}

func getTwitterStatusUrls(url string, channelID string) (map[string]string, error) {
	if twitterClient == nil {
		return nil, errors.New("invalid twitter api keys set")
	}

	matches := RegexpUrlTwitterStatus.FindStringSubmatch(url)
	statusId, err := strconv.ParseInt(matches[4], 10, 64)
	if err != nil {
		return nil, err
	}

	tweet, err := twitterClient.GetTweet(statusId, nil)
	if err != nil {
		return nil, err
	}

	links := make(map[string]string)
	for _, tweetMedia := range tweet.ExtendedEntities.Media {
		if len(tweetMedia.VideoInfo.Variants) > 0 {
			var lastVideoVariant anaconda.Variant
			for _, videoVariant := range tweetMedia.VideoInfo.Variants {
				if videoVariant.Bitrate >= lastVideoVariant.Bitrate {
					lastVideoVariant = videoVariant
				}
			}
			if lastVideoVariant.Url != "" {
				links[lastVideoVariant.Url] = ""
			}
		} else {
			foundUrls, _ := getDownloadLinks(tweetMedia.Media_url_https, channelID, false)
			for foundUrlKey, foundUrlValue := range foundUrls {
				links[foundUrlKey] = foundUrlValue
			}
		}
	}
	for _, tweetUrl := range tweet.Entities.Urls {
		foundUrls, _ := getDownloadLinks(tweetUrl.Expanded_url, channelID, false)
		for foundUrlKey, foundUrlValue := range foundUrls {
			links[foundUrlKey] = foundUrlValue
		}
	}

	return links, nil
}

func getGfycatUrls(url string) (map[string]string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		return nil, errors.New("unable to parse gfycat url")
	} else {
		gfycatId := parts[len(parts)-1]
		gfycatObject := new(GfycatObject)
		getJson("https://api.gfycat.com/v1/gfycats/"+gfycatId, gfycatObject)
		gfycatUrl := gfycatObject.GfyItem.Mp4URL
		if url == "" {
			return nil, errors.New("failed to read response from gfycat")
		}
		return map[string]string{gfycatUrl: ""}, nil
	}
}

func getInstagramUrls(url string) (map[string]string, error) {
	username, shortcode := getInstagramInfo(url)
	filename := fmt.Sprintf("instagram %s - %s", username, shortcode)
	// if instagram video
	videoUrl := getInstagramVideoUrl(url)
	if videoUrl != "" {
		return map[string]string{videoUrl: filename + filepathExtension(videoUrl)}, nil
	}
	// if instagram album
	albumUrls := getInstagramAlbumUrls(url)
	if len(albumUrls) > 0 {
		//fmt.Println("is instagram album")
		links := make(map[string]string)
		for i, albumUrl := range albumUrls {
			links[albumUrl] = filename + " " + strconv.Itoa(i+1) + filepathExtension(albumUrl)
		}
		return links, nil
	}
	// if instagram picture
	afterLastSlash := strings.LastIndex(url, "/")
	mediaUrl := url[:afterLastSlash]
	mediaUrl += strings.Replace(strings.Replace(url[afterLastSlash:], "?", "&", -1), "/", "/media/?size=l", -1)
	return map[string]string{mediaUrl: filename + ".jpg"}, nil
}

func getInstagramInfo(url string) (string, string) {
	resp, err := http.Get(url)

	if err != nil {
		return "N/A", "N/A"
	}

	defer resp.Body.Close()
	z := html.NewTokenizer(resp.Body)

ParseLoop:
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			break ParseLoop
		}
		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			t := z.Token()
			for _, a := range t.Attr {
				if a.Key == "type" {
					if a.Val == "text/javascript" {
						z.Next()
						content := string(z.Text())
						if strings.Contains(content, "window._sharedData = ") {
							content = strings.Replace(content, "window._sharedData = ", "", 1)
							content = content[:len(content)-1]
							jsonParsed, err := gabs.ParseJSON([]byte(content))
							if err != nil {
								fmt.Println("error parsing instagram json: ", err)
								continue ParseLoop
							}
							entryChildren, err := jsonParsed.Path("entry_data.PostPage").Children()
							if err != nil {
								fmt.Println("unable to find entries children: ", err)
								continue ParseLoop
							}
							for _, entryChild := range entryChildren {
								shortcode := entryChild.Path("graphql.shortcode_media.shortcode").Data().(string)
								username := entryChild.Path("graphql.shortcode_media.owner.username").Data().(string)
								return username, shortcode
							}
						}
					}
				}
			}
		}
	}
	return "N/A", "N/A"
}

func getImgurSingleUrls(url string) (map[string]string, error) {
	url = regexp.MustCompile(`(r\/[^\/]+\/)`).ReplaceAllString(url, "") // remove subreddit url
	fmt.Println(url)
	url = strings.Replace(url, "imgur.com/", "imgur.com/download/", -1)
	url = strings.Replace(url, ".gifv", "", -1)
	return map[string]string{url: ""}, nil
}

func getImgurAlbumUrls(url string) (map[string]string, error) {
	url = regexp.MustCompile(`(#[A-Za-z0-9]+)?$`).ReplaceAllString(url, "") // remove anchor
	afterLastSlash := strings.LastIndex(url, "/")
	albumId := url[afterLastSlash+1:]
	headers := make(map[string]string)
	headers["Authorization"] = "Client-ID " + IMGUR_CLIENT_ID
	imgurAlbumObject := new(ImgurAlbumObject)
	getJsonWithHeaders("https://api.imgur.com/3/album/"+albumId+"/images", imgurAlbumObject, headers)
	links := make(map[string]string)
	for _, v := range imgurAlbumObject.Data {
		links[v.Link] = ""
	}
	if len(links) <= 0 {
		return getImgurSingleUrls(url)
	}
	fmt.Printf("[%s] Found imgur album with %d images (url: %s)\n", time.Now().Format(time.Stamp), len(links), url)
	return links, nil
}

func getGoogleDriveUrls(url string) (map[string]string, error) {
	parts := strings.Split(url, "/")
	if len(parts) != 7 {
		return nil, errors.New("unable to parse google drive url")
	}
	fileId := parts[len(parts)-2]
	return map[string]string{"https://drive.google.com/uc?export=download&id=" + fileId: ""}, nil
}

func getGoogleDriveFolderUrls(url string) (map[string]string, error) {
	matches := RegexpUrlGoogleDriveFolder.FindStringSubmatch(url)
	if len(matches) < 4 || matches[3] == "" {
		return nil, errors.New("unable to find google drive folder ID in link")
	}
	if DriveService.BasePath == "" {
		return nil, errors.New("please set up google credentials")
	}
	googleDriveFolderID := matches[3]

	links := make(map[string]string)

	driveQuery := fmt.Sprintf("\"%s\" in parents", googleDriveFolderID)
	driveFields := "nextPageToken, files(id)"
	result, err := DriveService.Files.List().Q(driveQuery).Fields(googleapi.Field(driveFields)).PageSize(1000).Do()
	if err != nil {
		fmt.Println("driveQuery:", driveQuery)
		fmt.Println("driveFields:", driveFields)
		fmt.Println("err:", err)
		return nil, err
	}
	for _, file := range result.Files {
		fileUrl := "https://drive.google.com/uc?export=download&id=" + file.Id
		links[fileUrl] = ""
	}

	for {
		if result.NextPageToken == "" {
			break
		}
		result, err = DriveService.Files.List().Q(driveQuery).Fields(googleapi.Field(driveFields)).PageSize(1000).PageToken(result.NextPageToken).Do()
		if err != nil {
			return nil, err
		}
		for _, file := range result.Files {
			links[file.Id] = ""
		}
	}
	return links, nil
}

type FlickrPhotoSizeObject struct {
	Label  string `json:"label"`
	Width  int    `json:"width,int,string"`
	Height int    `json:"height,int,string"`
	Source string `json:"source"`
	URL    string `json:"url"`
	Media  string `json:"media"`
}

type FlickrPhotoObject struct {
	Sizes struct {
		Canblog     int                     `json:"canblog"`
		Canprint    int                     `json:"canprint"`
		Candownload int                     `json:"candownload"`
		Size        []FlickrPhotoSizeObject `json:"size"`
	} `json:"sizes"`
	Stat string `json:"stat"`
}

func getFlickrUrlFromPhotoId(photoId string) string {
	reqUrl := fmt.Sprintf("https://www.flickr.com/services/rest/?format=json&nojsoncallback=1&method=%s&api_key=%s&photo_id=%s",
		"flickr.photos.getSizes", flickrApiKey, photoId)
	flickrPhoto := new(FlickrPhotoObject)
	getJson(reqUrl, flickrPhoto)
	var bestSize FlickrPhotoSizeObject
	for _, size := range flickrPhoto.Sizes.Size {
		if bestSize.Label == "" {
			bestSize = size
		} else {
			if size.Width > bestSize.Width || size.Height > bestSize.Height {
				bestSize = size
			}
		}
	}
	return bestSize.Source
}

func getFlickrPhotoUrls(url string) (map[string]string, error) {
	if flickrApiKey == "" || flickrApiKey == "yourflickrapikey" || flickrApiKey == "your flickr api key" {
		return nil, errors.New("invalid flickr api key set")
	}
	matches := RegexpUrlFlickrPhoto.FindStringSubmatch(url)
	photoId := matches[5]
	if photoId == "" {
		return nil, errors.New("unable to get photo id from url")
	}
	return map[string]string{getFlickrUrlFromPhotoId(photoId): ""}, nil
}

type FlickrAlbumObject struct {
	Photoset struct {
		ID        string `json:"id"`
		Primary   string `json:"primary"`
		Owner     string `json:"owner"`
		Ownername string `json:"ownername"`
		Photo     []struct {
			ID        string `json:"id"`
			Secret    string `json:"secret"`
			Server    string `json:"server"`
			Farm      int    `json:"farm"`
			Title     string `json:"title"`
			Isprimary string `json:"isprimary"`
			Ispublic  int    `json:"ispublic"`
			Isfriend  int    `json:"isfriend"`
			Isfamily  int    `json:"isfamily"`
		} `json:"photo"`
		Page    int    `json:"page"`
		PerPage int    `json:"per_page"`
		Perpage int    `json:"perpage"`
		Pages   int    `json:"pages"`
		Total   string `json:"total"`
		Title   string `json:"title"`
	} `json:"photoset"`
	Stat string `json:"stat"`
}

func getFlickrAlbumUrls(url string) (map[string]string, error) {
	if flickrApiKey == "" || flickrApiKey == "yourflickrapikey" {
		return nil, errors.New("invalid flickr api key set")
	}
	matches := RegexpUrlFlickrAlbum.FindStringSubmatch(url)
	if len(matches) < 10 || matches[9] == "" {
		return nil, errors.New("unable to find flickr album ID in link")
	}
	albumId := matches[9]
	if albumId == "" {
		return nil, errors.New("unable to get album id from url")
	}
	reqUrl := fmt.Sprintf("https://www.flickr.com/services/rest/?format=json&nojsoncallback=1&method=%s&api_key=%s&photoset_id=%s&per_page=500",
		"flickr.photosets.getPhotos", flickrApiKey, albumId)
	flickrAlbum := new(FlickrAlbumObject)
	getJson(reqUrl, flickrAlbum)
	links := make(map[string]string)
	for _, photo := range flickrAlbum.Photoset.Photo {
		links[getFlickrUrlFromPhotoId(photo.ID)] = ""
	}
	fmt.Printf("[%s] Found flickr album with %d images (url: %s)\n", time.Now().Format(time.Stamp), len(links), url)
	return links, nil
}

func getFlickrAlbumShortUrls(url string) (map[string]string, error) {
	result, err := http.Get(url)
	if err != nil {
		return nil, errors.New("error getting long url from shortened flickr album url: " + err.Error())
	}
	if RegexpUrlFlickrAlbum.MatchString(result.Request.URL.String()) {
		return getFlickrAlbumUrls(result.Request.URL.String())
	}
	return nil, errors.New("got invalid url while trying to get long url from short flickr album url")
}

type StreamableObject struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Files  struct {
		Mp4 struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"mp4"`
		Mp4Mobile struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"mp4-mobile"`
	} `json:"files"`
	URL          string      `json:"url"`
	ThumbnailURL string      `json:"thumbnail_url"`
	Message      interface{} `json:"message"`
}

func getStreamableUrls(url string) (map[string]string, error) {
	matches := RegexpUrlStreamable.FindStringSubmatch(url)
	shortcode := matches[3]
	if shortcode == "" {
		return nil, errors.New("unable to get shortcode from url")
	}
	reqUrl := fmt.Sprintf("https://api.streamable.com/videos/%s", shortcode)
	streamable := new(StreamableObject)
	getJson(reqUrl, streamable)
	if streamable.Status != 2 || streamable.Files.Mp4.URL == "" {
		return nil, errors.New("streamable object has no download candidate")
	}
	link := streamable.Files.Mp4.URL
	if !strings.HasPrefix(link, "http") {
		link = "https:" + link
	}
	links := make(map[string]string)
	links[link] = ""
	return links, nil
}

func getPossibleTistorySiteUrls(url string) (map[string]string, error) {
	client := new(http.Client)
	request, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept-Encoding", "identity")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	respHead, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	contentType := ""
	for headerKey, headerValue := range respHead.Header {
		if headerKey == "Content-Type" {
			contentType = headerValue[0]
		}
	}
	if !strings.Contains(contentType, "text/html") {
		return nil, nil
	}

	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept-Encoding", "identity")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	var links = make(map[string]string)

	doc.Find(".article img, #content img, div[role=main] img, .section_blogview img").Each(func(i int, s *goquery.Selection) {
		foundUrl, exists := s.Attr("src")
		if exists {
			isTistoryCdnUrl := RegexpUrlTistoryLegacyWithCDN.MatchString(foundUrl)
			isTistoryUrl := RegexpUrlTistoryLegacy.MatchString(foundUrl)
			if isTistoryCdnUrl == true {
				finalTistoryUrls, _ := getTistoryWithCDNUrls(foundUrl)
				if len(finalTistoryUrls) > 0 {
					for finalTistoryUrl := range finalTistoryUrls {
						foundFilename := s.AttrOr("filename", "")
						links[finalTistoryUrl] = foundFilename
					}
				}
			} else if isTistoryUrl == true {
				finalTistoryUrls, _ := getLegacyTistoryUrls(foundUrl)
				if len(finalTistoryUrls) > 0 {
					for finalTistoryUrl := range finalTistoryUrls {
						foundFilename := s.AttrOr("filename", "")
						links[finalTistoryUrl] = foundFilename
					}
				}
			}
		}
	})

	if len(links) > 0 {
		fmt.Printf("[%s] Found tistory album with %d images (url: %s)\n", time.Now().Format(time.Stamp), len(links), url)
	}
	return links, nil
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getJsonWithHeaders(url string, target interface{}, headers map[string]string) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getInstagramVideoUrl(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		return ""
	}

	defer resp.Body.Close()
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return ""
		}
		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			t := z.Token()
			if t.Data == "meta" {
				for _, a := range t.Attr {
					if a.Key == "property" {
						if a.Val == "og:video" || a.Val == "og:video:secure_url" {
							for _, at := range t.Attr {
								if at.Key == "content" {
									return at.Val
								}
							}
						}
					}
				}
			}
		}
	}
}

func getInstagramAlbumUrls(url string) []string {
	var links []string
	resp, err := http.Get(url)

	if err != nil {
		return links
	}

	defer resp.Body.Close()
	z := html.NewTokenizer(resp.Body)

ParseLoop:
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			break ParseLoop
		}
		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			t := z.Token()
			for _, a := range t.Attr {
				if a.Key == "type" {
					if a.Val == "text/javascript" {
						z.Next()
						content := string(z.Text())
						if strings.Contains(content, "window._sharedData = ") {
							content = strings.Replace(content, "window._sharedData = ", "", 1)
							content = content[:len(content)-1]
							jsonParsed, err := gabs.ParseJSON([]byte(content))
							if err != nil {
								fmt.Println("error parsing instagram json: ", err)
								continue ParseLoop
							}
							entryChildren, err := jsonParsed.Path("entry_data.PostPage").Children()
							if err != nil {
								fmt.Println("unable to find entries children: ", err)
								continue ParseLoop
							}
							for _, entryChild := range entryChildren {
								albumChildren, err := entryChild.Path("graphql.shortcode_media.edge_sidecar_to_children.edges").Children()
								if err != nil {
									continue ParseLoop
								}
								for _, albumChild := range albumChildren {
									link, ok := albumChild.Path("node.display_url").Data().(string)
									if ok {
										links = append(links, link)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	if len(links) > 0 {
		fmt.Printf("[%s] Found instagram album with %d images (url: %s)\n", time.Now().Format(time.Stamp), len(links), url)
	}

	return links
}

func filenameFromUrl(dUrl string) string {
	base := path.Base(dUrl)
	parts := strings.Split(base, "?")
	return parts[0]
}

func startDownload(dUrl string, filename string, path string, channelId string, userId string, fileTime time.Time) {
	success := false
	for i := 0; i < MaxDownloadRetries; i++ {
		success = downloadFromUrl(dUrl, filename, path, channelId, userId, fileTime)
		if success == true {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	if success == false {
		fmt.Println("Gave up on downloading", dUrl)
		if SendNoticesToInteractiveChannels == true {
			for channelId := range InteractiveChannelWhitelist {
				content := fmt.Sprintf("Gave up on downloading %s, no success after %d retries", dUrl, MaxDownloadRetries)
				_, err := dg.ChannelMessageSend(channelId, content)
				if err != nil {
					fmt.Println("Failed to send notice to", channelId, "-", err)
				}
			}
		}
	}
}

func downloadFromUrl(dUrl string, filename string, path string, channelId string, userId string, fileTime time.Time) bool {

	var err error

	// Source validation
	_, err = url.ParseRequestURI(dUrl)
	if err != nil {
		fmt.Println("Error while parsing url", dUrl, "-", err)
		return false
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("Error while creating folder", path, "-", err)
		return false
	}

	timeout := time.Duration(time.Duration(DownloadTimeout) * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("GET", dUrl, nil)
	if err != nil {
		fmt.Println("Error while downloading", dUrl, "-", err)
		return false
	}
	request.Header.Add("Accept-Encoding", "identity")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error while downloading", dUrl, "-", err)
		return false
	}
	defer response.Body.Close()

	if filename == "" {
		filename = filenameFromUrl(response.Request.URL.String())
		for key, iHeader := range response.Header {
			if key == "Content-Disposition" {
				_, params, err := mime.ParseMediaType(iHeader[0])
				if err == nil {
					newFilename, err := url.QueryUnescape(params["filename"])
					if err != nil {
						newFilename = params["filename"]
					}
					if newFilename != "" {
						filename = newFilename
					}
				}
			}
		}
	}

	bodyOfResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Could not read response", dUrl, "-", err)
		return false
	}

	contentType := http.DetectContentType(bodyOfResp)

	// check for valid filename, if not, replace with generic filename
	if !RegexpFilename.MatchString(filename) {
		filename = time.Now().Format("2006-01-02 15-04-05")
		possibleExtension, _ := mime.ExtensionsByType(contentType)
		if len(possibleExtension) > 0 {
			filename += possibleExtension[0]
		}
	}

	completePath := path + string(os.PathSeparator) + filename
	if _, err := os.Stat(completePath); err == nil {
		tmpPath := completePath
		i := 1
		for {
			completePath = tmpPath[0:len(tmpPath)-len(filepathExtension(tmpPath))] +
				"-" + strconv.Itoa(i) + filepathExtension(tmpPath)
			if _, err := os.Stat(completePath); os.IsNotExist(err) {
				break
			}
			i = i + 1
		}
		fmt.Printf("[%s] Saving possible duplicate (filenames match): %s to %s\n", time.Now().Format(time.Stamp), tmpPath, completePath)
	}

	extension := filepath.Ext(filename)
	contentTypeParts := strings.Split(contentType, "/")
	if t := contentTypeParts[0]; t != "image" && t != "video" && t != "audio" &&
		!(t == "application" && isAudioFile(filename)) &&
		strings.ToLower(extension) != ".mov" &&
		strings.ToLower(extension) != ".mp4" &&
		strings.ToLower(extension) != ".webm" {
		fmt.Println("No image, video, or audio found at", dUrl)
		return true
	}

	err = ioutil.WriteFile(completePath, bodyOfResp, 0644)
	if err != nil {
		fmt.Println("Error while writing to disk", dUrl, "-", err)
		return false
	}

	err = os.Chtimes(completePath, fileTime, fileTime)
	if err != nil {
		fmt.Println("Error while changing date", dUrl, "-", err)
	}

	sourceChannelName := channelId
	sourceGuildName := "N/A"
	sourceChannel, _ := dg.State.Channel(channelId)
	if sourceChannel != nil && sourceChannel.Name != "" {
		sourceChannelName = sourceChannel.Name
		sourceGuild, _ := dg.State.Guild(sourceChannel.GuildID)
		if sourceGuild != nil && sourceGuild.Name != "" {
			sourceGuildName = sourceGuild.Name
		}
	}

	fmt.Printf("[%s] Saved URL %s to %s from #%s/%s\n",
		time.Now().Format(time.Stamp), dUrl, completePath, sourceChannelName, sourceGuildName)
	err = insertDownloadedImage(&DownloadedImage{Url: dUrl, Time: time.Now(), Destination: completePath, ChannelId: channelId, UserId: userId})
	if err != nil {
		fmt.Println("Error while writing to database", err)
	}

	updateDiscordStatus()
	return true
}

func isAudioFile(f string) bool {
	switch strings.ToLower(path.Ext(f)) {
	case ".mp3", ".wav", ".aif":
		return true
	default:
		return false
	}
}

type DownloadedImage struct {
	Url         string
	Time        time.Time
	Destination string
	ChannelId   string
	UserId      string
}

func insertDownloadedImage(downloadedImage *DownloadedImage) error {
	_, err := myDB.Use("Downloads").Insert(map[string]interface{}{
		"Url":         downloadedImage.Url,
		"Time":        downloadedImage.Time.String(),
		"Destination": downloadedImage.Destination,
		"ChannelId":   downloadedImage.ChannelId,
		"UserId":      downloadedImage.UserId,
	})
	return err
}

func findDownloadedImageById(id int) *DownloadedImage {
	downloads := myDB.Use("Downloads")
	//var query interface{}
	//json.Unmarshal([]byte(fmt.Sprintf(`[{"eq": "%d", "in": ["Id"]}]`, id)), &query)
	//queryResult := make(map[int]struct{})
	//db.EvalQuery(query, myDB.Use("Downloads"), &queryResult)

	readBack, err := downloads.Read(id)
	if err != nil {
		fmt.Println(err)
	}
	timeT, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", readBack["Time"].(string))
	return &DownloadedImage{
		Url:         readBack["Url"].(string),
		Time:        timeT,
		Destination: readBack["Destination"].(string),
		ChannelId:   readBack["ChannelId"].(string),
		UserId:      readBack["UserId"].(string),
	}
}

func findDownloadedImageByUrl(url string) []*DownloadedImage {
	var query interface{}
	json.Unmarshal([]byte(fmt.Sprintf(`[{"eq": "%s", "in": ["Url"]}]`, url)), &query)
	queryResult := make(map[int]struct{})
	db.EvalQuery(query, myDB.Use("Downloads"), &queryResult)

	downloadedImages := make([]*DownloadedImage, 0)
	for id := range queryResult {
		downloadedImages = append(downloadedImages, findDownloadedImageById(id))
	}
	return downloadedImages
}

func countDownloadedImages() int {
	i := 0
	myDB.Use("Downloads").ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		// fmt.Printf("%v\n", findDownloadedImageById(id))
		i++
		return true
	})
	return i
	// fmt.Println(myDB.Use("Downloads").ApproxDocCount()) TODO?
}

// http://stackoverflow.com/a/18695740/1443726
func sortStringIntMapByValue(m map[string]int) PairList {
	pl := make(PairList, len(m))
	i := 0
	for k, v := range m {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
