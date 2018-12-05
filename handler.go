package main

import (
	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls"
)

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
