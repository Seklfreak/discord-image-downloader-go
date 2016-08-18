package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mvdan/xurls"
	"gopkg.in/ini.v1"
)

var (
	ChannelWhitelist map[string]string
	BaseDownloadPath string
	RegexpUrlTwitter *regexp.Regexp
	RegexpUrlTistory *regexp.Regexp
	RegexpUrlGfycat  *regexp.Regexp
)

type GfycatObject struct {
	GfyItem map[string]string
}

func main() {
	var err error
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("unable to read config file", err)
		cfg = ini.Empty()
	}

	if !cfg.Section("auth").HasKey("email") ||
		!cfg.Section("auth").HasKey("password") {
		cfg.Section("auth").NewKey("email", "your@email.com")
		cfg.Section("auth").NewKey("password", "yourpassword")
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
		`^http(s?):\/\/pbs\.twimg\.com\/media\/[a-zA-Z0-9]+\.jpg((\:[a-z]+)?)$`)
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

	dg, err := discordgo.New(
		cfg.Section("auth").Key("email").String(),
		cfg.Section("auth").Key("password").String())
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Client is now connected. Press CTRL-C to exit.")
	// keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
			} else {
				// Any other url
				downloadFromUrl(iFoundUrl,
					getContentDispositionFilename(iFoundUrl), downloadPath)
			}
		}
	}
}

func handleTwitterUrl(url string, folder string) error {
	parts := strings.Split(url, ":")
	if len(parts) < 2 {
		return errors.New("unable to parse twitter url")
	} else {
		downloadFromUrl("https:"+parts[1]+":orig", path.Base(parts[1]), folder)
	}
	return nil
}

func handleTistoryUrl(url string, folder string) error {
	url = strings.Replace(url, "/image/", "/original/", -1)
	downloadFromUrl(url, getContentDispositionFilename(url), folder)
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
		if url == "" {
			return errors.New("failed to read response from gfycat")
		} else {
			downloadFromUrl(
				gfycatUrl, getContentDispositionFilename(gfycatUrl), folder)
		}
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

func getContentDispositionFilename(dUrl string) string {
	resp, err := http.Head(dUrl)
	if err != nil {
		return path.Base(dUrl)
	}
	for key, iHeader := range resp.Header {
		if key == "Content-Disposition" {
			parts := strings.Split(iHeader[0], "\"")
			if len(parts) == 3 {
				filename, err := url.QueryUnescape(parts[1])
				if err != nil {
					return parts[1]
				} else {
					return filename
				}
			}
		}
	}
	return path.Base(dUrl)
}

func downloadFromUrl(url string, filename string, path string) {
	err := os.MkdirAll(path, 755)
	if err != nil {
		fmt.Println("Error while creating folder", path, "-", err)
		return
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
	}

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	bodyOfResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Could not read response", url, "-", err)
		return
	}
	contentType := http.DetectContentType(bodyOfResp)
	contentTypeParts := strings.Split(contentType, "/")
	if contentTypeParts[0] != "image" {
		fmt.Println("No image found at", url)
		return
	}

	err = ioutil.WriteFile(completePath, bodyOfResp, 0644)
	if err != nil {
		fmt.Println("Error while writing to disk", url, "-", err)
		return
	}

	fmt.Printf("[%s] Downloaded url: %s to %s\n", time.Now().Format(time.Stamp), url, completePath)
}
