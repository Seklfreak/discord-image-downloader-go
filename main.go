package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/bwmarrin/discordgo"
	"github.com/mvdan/xurls"
	"gopkg.in/ini.v1"
)

var (
	ChannelWhitelist     map[string]string
	BaseDownloadPath     string
	RegexpUrlTwitter     *regexp.Regexp
	RegexpUrlTistory     *regexp.Regexp
	RegexpUrlGfycat      *regexp.Regexp
	RegexpUrlInstagram   *regexp.Regexp
	RegexpUrlImgurGifv   *regexp.Regexp
	RegexpUrlImgurAlbum  *regexp.Regexp
	RegexpUrlGoogleDrive *regexp.Regexp
	ImagesDownloaded     int
	dg                   *discordgo.Session
)

const (
	VERSION         string = "1.9"
	RELEASE_URL     string = "https://github.com/Seklfreak/discord-image-downloader-go/releases/latest"
	IMGUR_CLIENT_ID string = "a39473314df3f59"
)

type GfycatObject struct {
	GfyItem map[string]string
}

type ImgurAlbumObject struct {
	Data []struct {
		Link string
	}
}

func main() {
	fmt.Printf("discord-image-downloader-go version %s\n", VERSION)
	fmt.Printf("Go to %s to get the latest release.\n", RELEASE_URL)

	var err error
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("unable to read config file", err)
		cfg = ini.Empty()
	}

	if (!cfg.Section("auth").HasKey("email") ||
		!cfg.Section("auth").HasKey("password")) &&
		!cfg.Section("auth").HasKey("token") {
		cfg.Section("auth").NewKey("email", "your@email.com")
		cfg.Section("auth").NewKey("password", "yourpassword")
		cfg.Section("general").NewKey("skip edits", "true")
		cfg.Section("channels").NewKey("channelid1", "C:\\full\\path\\1")
		cfg.Section("channels").NewKey("channelid2", "C:\\full\\path\\2")
		cfg.Section("channels").NewKey("channelid3", "C:\\full\\path\\3")
		err = cfg.SaveTo("config.ini")

		if err != nil {
			fmt.Println("unable to write config file", err)
			return
		}
		fmt.Println("Wrote config file, please fill out and restart the program")
		return
	}

	ChannelWhitelist = cfg.Section("channels").KeysHash()

	RegexpUrlTwitter, err = regexp.Compile(
		`^http(s?):\/\/pbs\.twimg\.com\/media\/[^\./]+\.jpg((\:[a-z]+)?)$`)
	if err != nil {
		fmt.Println("Regexp error", err)
		return
	}
	RegexpUrlTistory, err = regexp.Compile(
		`^http(s?):\/\/[a-z0-9]+\.uf\.tistory\.com\/(image|original)\/[A-Z0-9]+$`)
	if err != nil {
		fmt.Println("Regexp error", err)
		return
	}
	RegexpUrlGfycat, err = regexp.Compile(
		`^http(s?):\/\/gfycat\.com\/[A-Za-z]+$`)
	if err != nil {
		fmt.Println("Regexp error", err)
		return
	}
	RegexpUrlInstagram, err = regexp.Compile(
		`^http(s?):\/\/(www\.)?instagram\.com\/p\/[^/]+\/(\?[^/]+)?$`)
	if err != nil {
		fmt.Println("Regexp error", err)
		return
	}
	RegexpUrlImgurGifv, err = regexp.Compile(
		`^http(s?):\/\/i\.imgur\.com\/[A-Za-z0-9]+\.gifv$`)
	if err != nil {
		fmt.Println("Regexp error", err)
		return
	}
	RegexpUrlImgurAlbum, err = regexp.Compile(
		`^http(s?):\/\/imgur\.com\/a\/[A-Za-z0-9]+$`)
	if err != nil {
		fmt.Println("Regexp error", err)
		return
	}
	RegexpUrlGoogleDrive, err = regexp.Compile(
		`^http(s?):\/\/drive\.google\.com\/file\/d\/[^/]+\/view$`)
	if err != nil {
		fmt.Println("Regexp error", err)
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
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	if cfg.Section("general").HasKey("skip edits") {
		if cfg.Section("general").Key("skip edits").MustBool() == false {
			dg.AddHandler(messageUpdate)
		}
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	fmt.Printf("Client is now connected as %s (ID: %s). Press CTRL-C to exit.\n",
		u.Username, u.ID)

	err = dg.UpdateStatus(1, "")
	if err != nil {
		fmt.Println("error setting idle status,", err)
	}

	// keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	handleDiscordMessage(m.Message)
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	handleDiscordMessage(m.Message)
}

func handleDiscordMessage(m *discordgo.Message) {
	if folderName, ok := ChannelWhitelist[m.ChannelID]; ok {
		downloadPath := folderName
		for _, iAttachment := range m.Attachments {
			downloadFromUrl(iAttachment.URL, iAttachment.Filename, downloadPath)
		}
		foundUrls := xurls.Strict.FindAllString(m.Content, -1)
		for _, iFoundUrl := range foundUrls {
			// Twitter url?
			if RegexpUrlTwitter.MatchString(iFoundUrl) {
				err := handleTwitterUrl(iFoundUrl, downloadPath)
				if err != nil {
					fmt.Println("twitter url failed,", iFoundUrl, ",", err)
					continue
				}
				// Tistory url?
			} else if RegexpUrlTistory.MatchString(iFoundUrl) {
				err := handleTistoryUrl(iFoundUrl, downloadPath)
				if err != nil {
					fmt.Println("tistory url failed,", iFoundUrl, ",", err)
					continue
				}
			} else if RegexpUrlGfycat.MatchString(iFoundUrl) {
				err := handleGfycatUrl(iFoundUrl, downloadPath)
				if err != nil {
					fmt.Println("gfycat url failed,", iFoundUrl, ",", err)
					continue
				}
			} else if RegexpUrlInstagram.MatchString(iFoundUrl) {
				err := handleInstagramUrl(iFoundUrl, downloadPath)
				if err != nil {
					fmt.Println("instagram url failed,", iFoundUrl, ",", err)
					continue
				}
			} else if RegexpUrlImgurGifv.MatchString(iFoundUrl) {
				err := handleImgurGifvUrl(iFoundUrl, downloadPath)
				if err != nil {
					fmt.Println("imgur gifv url failed, ", iFoundUrl, ",", err)
					continue
				}
			} else if RegexpUrlImgurAlbum.MatchString(iFoundUrl) {
				err := handleImgurAlbumUrl(iFoundUrl, downloadPath)
				if err != nil {
					fmt.Println("imgur album url failed, ", iFoundUrl, ",", err)
					continue
				}
			} else if RegexpUrlGoogleDrive.MatchString(iFoundUrl) {
				err := handleGoogleDriveUrl(iFoundUrl, downloadPath)
				if err != nil {
					fmt.Println("google drive album url failed, ", iFoundUrl, ",", err)
					continue
				}
			} else {
				// Any other url
				downloadFromUrl(iFoundUrl,
					"", downloadPath)
			}
		}
	}
}

func handleTwitterUrl(url string, folder string) error {
	parts := strings.Split(url, ":")
	if len(parts) < 2 {
		return errors.New("unable to parse twitter url")
	} else {
		downloadFromUrl("https:"+parts[1]+":orig", filenameFromUrl(parts[1]), folder)
	}
	return nil
}

func handleTistoryUrl(url string, folder string) error {
	url = strings.Replace(url, "/image/", "/original/", -1)
	downloadFromUrl(url, "", folder)
	return nil
}

func handleGfycatUrl(url string, folder string) error {
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		return errors.New("unable to parse gfycat url")
	} else {
		gfycatId := parts[len(parts)-1]
		gfycatObject := new(GfycatObject)
		getJson("https://gfycat.com/cajax/get/"+gfycatId, gfycatObject)
		gfycatUrl := gfycatObject.GfyItem["gifUrl"]
		if gfycatUrl == "" {
			gfycatUrl = gfycatObject.GfyItem["mp4Url"]
			fmt.Println("fallback to gfycat mp4")
		}
		if url == "" {
			return errors.New("failed to read response from gfycat")
		} else {
			downloadFromUrl(
				gfycatUrl, "", folder)
		}
	}
	return nil
}

func getInstagramVideoUrl(url string) string {
	resp, _ := http.Get(url)

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

func handleInstagramUrl(url string, folder string) error {
	// if instagram video
	videoUrl := getInstagramVideoUrl(url)
	if videoUrl != "" {
		downloadFromUrl(videoUrl, "", folder)
		return nil
	}

	// if instagram picture
	afterLastSlash := strings.LastIndex(url, "/")
	mediaUrl := url[:afterLastSlash] + strings.Replace(url[afterLastSlash:], "/", "/media/?size=l", -1)
	mediaUrl = strings.Replace(mediaUrl, "?taken-by=", "&taken-by", -1)
	downloadFromUrl(mediaUrl, "", folder)
	return nil
}

func handleImgurGifvUrl(url string, folder string) error {
	url = strings.Replace(url, "i.imgur.com/", "imgur.com/download/", -1)
	url = strings.Replace(url, ".gifv", "", -1)
	downloadFromUrl(url, "", folder)
	return nil
}

func handleImgurAlbumUrl(url string, folder string) error {
	afterLastSlash := strings.LastIndex(url, "/")
	albumId := url[afterLastSlash+1:]
	headers := make(map[string]string)
	headers["Authorization"] = "Client-ID " + IMGUR_CLIENT_ID
	imgurAlbumObject := new(ImgurAlbumObject)
	getJsonWithHeaders("https://api.imgur.com/3/album/"+albumId+"/images", imgurAlbumObject, headers)
	fmt.Printf("[%s] Found imgur album url: %s\n", time.Now().Format(time.Stamp), url)
	for _, v := range imgurAlbumObject.Data {
		downloadFromUrl(v.Link, "", folder)
	}
	return nil
}

func handleGoogleDriveUrl(url string, folder string) error {
	parts := strings.Split(url, "/")
	if len(parts) != 7 {
		return errors.New("unable to parse google drive url")
	} else {
		fileId := parts[len(parts)-2]
		downloadFromUrl("https://drive.google.com/uc?export=download&id="+fileId, "", folder)
	}
	return nil
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

func filenameFromUrl(dUrl string) string {
	base := path.Base(dUrl)
	parts := strings.Split(base, "?")
	return parts[0]
}

func downloadFromUrl(dUrl string, filename string, path string) {
	err := os.MkdirAll(path, 755)
	if err != nil {
		fmt.Println("Error while creating folder", path, "-", err)
		return
	}

	response, err := http.Get(dUrl)
	if err != nil {
		fmt.Println("Error while downloading", dUrl, "-", err)
		return
	}
	defer response.Body.Close()

	if filename == "" {
		filename = filenameFromUrl(response.Request.URL.String())
		for key, iHeader := range response.Header {
			if key == "Content-Disposition" {
				_, params, err := mime.ParseMediaType(iHeader[0])
				newFilename := params["filename"]
				if err == nil && newFilename != "" {
					filename = newFilename
				}
			}
		}
	}

	completePath := path + string(os.PathSeparator) + filename
	if _, err := os.Stat(completePath); err == nil {
		tmpPath := completePath
		i := 1
		for {
			completePath = tmpPath[0:len(tmpPath)-len(filepath.Ext(tmpPath))] +
				"-" + strconv.Itoa(i) + filepath.Ext(tmpPath)
			if _, err := os.Stat(completePath); os.IsNotExist(err) {
				break
			}
			i = i + 1
		}
		fmt.Printf("[%s] Saving possible duplicate (filenames match): %s to %s\n", time.Now().Format(time.Stamp), tmpPath, completePath)
	}

	bodyOfResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Could not read response", dUrl, "-", err)
		return
	}
	contentType := http.DetectContentType(bodyOfResp)
	contentTypeParts := strings.Split(contentType, "/")
	if contentTypeParts[0] != "image" && contentTypeParts[0] != "video" {
		fmt.Println("No image or video found at", dUrl)
		return
	}

	err = ioutil.WriteFile(completePath, bodyOfResp, 0644)
	if err != nil {
		fmt.Println("Error while writing to disk", dUrl, "-", err)
		return
	}

	fmt.Printf("[%s] Downloaded url: %s to %s\n", time.Now().Format(time.Stamp), dUrl, completePath)
	updateDiscordStatus()
}

func updateDiscordStatus() {
	ImagesDownloaded++
	dg.UpdateStatus(0, fmt.Sprintf("%d pictures downloaded", ImagesDownloaded))
}
