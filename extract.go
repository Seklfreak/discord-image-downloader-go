package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls"
)

func getDownloadLinks(inputURL string, channelID string, interactive bool) map[string]string {
	if RegexpUrlTwitter.MatchString(inputURL) {
		links, err := getTwitterUrls(inputURL)
		if err != nil {
			fmt.Println("twitter URL failed,", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlTwitterStatus.MatchString(inputURL) {
		links, err := getTwitterStatusUrls(inputURL, channelID)
		if err != nil {
			fmt.Println("twitter status URL failed,", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlTistory.MatchString(inputURL) {
		links, err := getTistoryUrls(inputURL)
		if err != nil {
			fmt.Println("tistory URL failed,", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlTistoryLegacy.MatchString(inputURL) {
		links, err := getLegacyTistoryUrls(inputURL)
		if err != nil {
			fmt.Println("legacy tistory URL failed,", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlGfycat.MatchString(inputURL) {
		links, err := getGfycatUrls(inputURL)
		if err != nil {
			fmt.Println("gfycat URL failed,", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlInstagram.MatchString(inputURL) {
		links, err := getInstagramUrls(inputURL)
		if err != nil {
			fmt.Println("instagram URL failed,", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlImgurSingle.MatchString(inputURL) {
		links, err := getImgurSingleUrls(inputURL)
		if err != nil {
			fmt.Println("imgur single URL failed, ", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlImgurAlbum.MatchString(inputURL) {
		links, err := getImgurAlbumUrls(inputURL)
		if err != nil {
			fmt.Println("imgur album URL failed, ", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlGoogleDrive.MatchString(inputURL) {
		links, err := getGoogleDriveUrls(inputURL)
		if err != nil {
			fmt.Println("google drive album URL failed, ", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlFlickrPhoto.MatchString(inputURL) {
		links, err := getFlickrPhotoUrls(inputURL)
		if err != nil {
			fmt.Println("flickr photo URL failed, ", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlFlickrAlbum.MatchString(inputURL) {
		links, err := getFlickrAlbumUrls(inputURL)
		if err != nil {
			fmt.Println("flickr album URL failed, ", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlFlickrAlbumShort.MatchString(inputURL) {
		links, err := getFlickrAlbumShortUrls(inputURL)
		if err != nil {
			fmt.Println("flickr album short URL failed, ", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if RegexpUrlStreamable.MatchString(inputURL) {
		links, err := getStreamableUrls(inputURL)
		if err != nil {
			fmt.Println("streamable URL failed, ", inputURL, ",", err)
		} else if len(links) > 0 {
			return skipDuplicateLinks(links, channelID, interactive)
		}
	}
	if DownloadTistorySites {
		if RegexpUrlPossibleTistorySite.MatchString(inputURL) {
			links, err := getPossibleTistorySiteUrls(inputURL)
			if err != nil {
				fmt.Println("checking for tistory site failed, ", inputURL, ",", err)
			} else if len(links) > 0 {
				return skipDuplicateLinks(links, channelID, interactive)
			}
		}
	}
	if RegexpUrlGoogleDriveFolder.MatchString(inputURL) {
		if interactive {
			links, err := getGoogleDriveFolderUrls(inputURL)
			if err != nil {
				fmt.Println("google drive folder URL failed, ", inputURL, ",", err)
			} else if len(links) > 0 {
				return skipDuplicateLinks(links, channelID, interactive)
			}
		} else {
			fmt.Println("google drive folder only accepted in interactive channels")
		}
	}

	if !interactive && isDiscordEmoji(inputURL) {
		fmt.Printf("skipped %s as it is a Discord emoji\n", inputURL)
		return nil
	}

	// try without queries
	parsedURL, err := url.Parse(inputURL)
	if err == nil {
		parsedURL.RawQuery = ""
		inputURLWithoutQueries := parsedURL.String()
		if inputURLWithoutQueries != inputURL {
			return skipDuplicateLinks(getDownloadLinks(inputURLWithoutQueries, channelID, interactive), channelID, interactive)
		}
	}

	return skipDuplicateLinks(map[string]string{inputURL: ""}, channelID, interactive)
}

// getDownloadItemsOfMessage will extract all unique download links out of a message
func getDownloadItemsOfMessage(message *discordgo.Message) []*DownloadItem {
	var downloadItems []*DownloadItem

	linkTime, err := message.Timestamp.Parse()
	if err != nil {
		linkTime = time.Now()
	}

	rawLinks := getRawLinksOfMessage(message)
	for _, rawLink := range rawLinks {
		downloadLinks := getDownloadLinks(
			rawLink.Link,
			message.ChannelID,
			false,
		)
		for link, filename := range downloadLinks {
			if rawLink.Filename != "" {
				filename = rawLink.Filename
			}

			downloadItems = append(downloadItems, &DownloadItem{
				Link:     link,
				Filename: filename,
				Time:     linkTime,
			})
		}
	}

	downloadItems = deduplicateDownloadItems(downloadItems)

	return downloadItems
}

// getRawLinksOfMessage will extract all raw links of a message
func getRawLinksOfMessage(message *discordgo.Message) []*DownloadItem {
	var links []*DownloadItem

	if message.Author == nil {
		message.Author = new(discordgo.User)
	}

	for _, attachment := range message.Attachments {
		links = append(links, &DownloadItem{
			Link:     attachment.URL,
			Filename: attachment.Filename,
		})
	}

	foundLinks := xurls.Strict.FindAllString(message.Content, -1)
	for _, foundLink := range foundLinks {
		links = append(links, &DownloadItem{
			Link: foundLink,
		})
	}

	for _, embed := range message.Embeds {
		if embed.URL != "" {
			links = append(links, &DownloadItem{
				Link: embed.URL,
			})
		}

		if embed.Description != "" {
			foundLinks = xurls.Strict.FindAllString(embed.Description, -1)
			for _, foundLink := range foundLinks {
				links = append(links, &DownloadItem{
					Link: foundLink,
				})
			}
		}

		if embed.Image != nil && embed.Image.URL != "" {
			links = append(links, &DownloadItem{
				Link: embed.Image.URL,
			})
		}

		if embed.Video != nil && embed.Video.URL != "" {
			links = append(links, &DownloadItem{
				Link: embed.Video.URL,
			})
		}
	}

	return links
}
