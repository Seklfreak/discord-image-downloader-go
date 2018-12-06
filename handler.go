package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls"
)

// helpHandler displays the help message
func helpHandler(message *discordgo.Message) {
	dg.ChannelMessageSend(
		message.ChannelID,
		"**<link>** to download a link\n"+
			"**version** to find out the version\n"+
			"**stats** to view stats\n"+
			"**channels** to list active channels\n"+
			"**history** to download a full channel history\n"+
			"**help** to open this help\n",
	)
}

// versionHandler displays the current version
func versionHandler(message *discordgo.Message) {
	result := fmt.Sprintf("discord-image-downloder-go **v%s**\n", VERSION)

	if !isLatestRelease() {
		result += fmt.Sprintf("**update available on <%s>**", RELEASE_URL)
	} else {
		result += "version is up to date"
	}

	dg.ChannelMessageSend(message.ChannelID, result)
}

// channelsHandler displays enabled channels
func channelsHandler(message *discordgo.Message) {
	var err error
	var channel *discordgo.Channel
	var channelRecipientUsername string

	result := "**Channels**\n"
	for channelID, channelFolder := range ChannelWhitelist {
		channel, err = dg.State.Channel(channelID)
		if err != nil {
			continue
		}

		switch channel.Type {
		case discordgo.ChannelTypeDM:
			channelRecipientUsername = ""
			for _, recipient := range channel.Recipients {
				channelRecipientUsername = recipient.Username + ", "
			}
			channelRecipientUsername = strings.TrimRight(channelRecipientUsername, ", ")
			if channelRecipientUsername == "" {
				channelRecipientUsername = "N/A"
			}
			result += fmt.Sprintf("@%s (`#%s`): `%s`\n", channelRecipientUsername, channel.ID, channelFolder)
		default:
			result += fmt.Sprintf("<#%s> (`#%s`): `%s`\n", channel.ID, channel.ID, channelFolder)
		}

	}
	result += "**Interactive Channels**\n"
	for channelID, channelFolder := range InteractiveChannelWhitelist {
		channel, err = dg.State.Channel(channelID)
		if err != nil {
			continue
		}

		switch channel.Type {
		case discordgo.ChannelTypeDM:
			channelRecipientUsername = ""
			for _, recipient := range channel.Recipients {
				channelRecipientUsername = recipient.Username + ", "
			}
			channelRecipientUsername = strings.TrimRight(channelRecipientUsername, ", ")
			if channelRecipientUsername == "" {
				channelRecipientUsername = "N/A"
			}
			result += fmt.Sprintf("@%s (`#%s`): `%s`\n", channelRecipientUsername, channel.ID, channelFolder)
		default:
			result += fmt.Sprintf("<#%s> (`#%s`): `%s`\n", channel.ID, channel.ID, channelFolder)
		}
	}

	for _, page := range Pagify(result, "\n") {
		dg.ChannelMessageSend(message.ChannelID, page)
	}
}

func statsHandler(message *discordgo.Message) {
	channelStats := make(map[string]int)
	userStats := make(map[string]int)
	userGuilds := make(map[string]string)

	var i int
	myDB.Use("Downloads").ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		downloadedImage := findDownloadedImageById(id)
		channelStats[downloadedImage.ChannelId] += 1
		userStats[downloadedImage.UserId] += 1
		if _, ok := userGuilds[downloadedImage.UserId]; !ok {
			channel, err := dg.State.Channel(downloadedImage.ChannelId)
			if err == nil && channel.GuildID != "" {
				userGuilds[downloadedImage.UserId] = channel.GuildID
			}
		}
		i++
		return true
	})
	channelStatsSorted := sortStringIntMapByValue(channelStats)
	userStatsSorted := sortStringIntMapByValue(userStats)
	replyMessage := fmt.Sprintf("I downloaded **%d** pictures in **%d** channels by **%d** users\n", i, len(channelStats), len(userStats))

	replyMessage += "**channel breakdown**\n"
	for _, downloads := range channelStatsSorted {
		channel, err := dg.State.Channel(downloads.Key)
		if err == nil {
			if channel.Type == discordgo.ChannelTypeDM {
				channelRecipientUsername := "N/A"
				for _, recipient := range channel.Recipients {
					channelRecipientUsername = recipient.Username
				}
				replyMessage += fmt.Sprintf("@%s (`#%s`): **%d** downloads\n", channelRecipientUsername, downloads.Key, downloads.Value)
			} else {
				guild, err := dg.State.Guild(channel.GuildID)
				if err == nil {
					replyMessage += fmt.Sprintf("#%s/%s (`#%s`): **%d** downloads\n", guild.Name, channel.Name, downloads.Key, downloads.Value)
				} else {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println(err)
		}
	}
	replyMessage += "**user breakdown**\n"
	var userI int
	for _, downloads := range userStatsSorted {
		userI++
		if userI > 10 {
			replyMessage += "_only the top 10 users get shown_\n"
			break
		}
		if guildId, ok := userGuilds[downloads.Key]; ok {
			user, err := dg.State.Member(guildId, downloads.Key)
			if err == nil {
				replyMessage += fmt.Sprintf("@%s: **%d** downloads\n", user.User.Username, downloads.Value)
			} else {
				replyMessage += fmt.Sprintf("@`%s`: **%d** downloads\n", downloads.Key, downloads.Value)
			}
		} else {
			replyMessage += fmt.Sprintf("@`%s`: **%d** downloads\n", downloads.Key, downloads.Value)
		}
	}

	for _, page := range Pagify(replyMessage, "\n") {
		dg.ChannelMessageSend(message.ChannelID, page)
	}
}

func historyHandler(message *discordgo.Message) {
	i := 0
	_, historyCommandIsSet := historyCommandActive[message.ChannelID]
	if !historyCommandIsSet || historyCommandActive[message.ChannelID] == "" {
		historyCommandActive[message.ChannelID] = ""

		idArray := strings.Split(message.Content, ",")
		for _, channelValue := range idArray {
			channelValue = strings.TrimSpace(channelValue)
			if folder, ok := ChannelWhitelist[channelValue]; ok {
				dg.ChannelMessageSend(message.ChannelID, fmt.Sprintf("downloading to `%s`", folder))
				historyCommandActive[message.ChannelID] = "downloading"
				lastBefore := ""
				var lastBeforeTime time.Time
			MessageRequestingLoop:
				for true {
					if lastBeforeTime != (time.Time{}) {
						fmt.Printf("[%s] Requesting 100 more messages, (before %s)\n", time.Now().Format(time.Stamp), lastBeforeTime)
						dg.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Requesting 100 more messages, (before %s)\n", lastBeforeTime))
					}
					messages, err := dg.ChannelMessages(channelValue, 100, lastBefore, "", "")
					if err == nil {
						if len(messages) <= 0 {
							delete(historyCommandActive, message.ChannelID)
							break MessageRequestingLoop
						}
						lastBefore = messages[len(messages)-1].ID
						lastBeforeTime, err = messages[len(messages)-1].Timestamp.Parse()
						if err != nil {
							fmt.Println(err)
						}
						for _, message := range messages {
							fileTime := time.Now()
							if message.Timestamp != "" {
								fileTime, err = message.Timestamp.Parse()
								if err != nil {
									fmt.Println(err)
								}
							}
							if historyCommandActive[message.ChannelID] == "cancel" {
								delete(historyCommandActive, message.ChannelID)
								break MessageRequestingLoop
							}
							for _, iAttachment := range message.Attachments {
								if len(findDownloadedImageByUrl(iAttachment.URL)) == 0 {
									i++
									startDownload(iAttachment.URL, iAttachment.Filename, folder, message.ChannelID, message.Author.ID, fileTime)
								}
							}
							foundUrls := xurls.Strict.FindAllString(message.Content, -1)
							for _, iFoundUrl := range foundUrls {
								links := getDownloadLinks(iFoundUrl, message.ChannelID, false)
								for link, filename := range links {
									if len(findDownloadedImageByUrl(link)) == 0 {
										i++
										startDownload(link, filename, folder, message.ChannelID, message.Author.ID, fileTime)
									}
								}
							}
						}
					} else {
						dg.ChannelMessageSend(message.ChannelID, err.Error())
						fmt.Println(err)
						delete(historyCommandActive, message.ChannelID)
						break MessageRequestingLoop
					}
				}
				dg.ChannelMessageSend(message.ChannelID, fmt.Sprintf("done, %d download links started!", i))
			} else {
				dg.ChannelMessageSend(message.ChannelID, "Please tell me one or multiple Channel IDs (separated by commas)\nPlease make sure the channels have been whitelisted before submitting.")
			}
		}
	} else if historyCommandActive[message.ChannelID] == "downloading" && strings.ToLower(message.Content) == "cancel" {
		historyCommandActive[message.ChannelID] = "cancel"
	}
}

func defaultHandler(message *discordgo.Message) {
	folderName := InteractiveChannelWhitelist[message.ChannelID]

	if link, ok := interactiveChannelLinkTemp[message.ChannelID]; ok {
		fileTime := time.Now()
		var err error
		if message.Timestamp != "" {
			fileTime, err = message.Timestamp.Parse()
			if err != nil {
				fmt.Println(err)
			}
		}
		if message.Content == "." {
			dg.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Download of <%s> started", link))
			dg.ChannelTyping(message.ChannelID)
			delete(interactiveChannelLinkTemp, message.ChannelID)
			links := getDownloadLinks(link, message.ChannelID, true)
			for linkR, filename := range links {
				startDownload(linkR, filename, folderName, message.ChannelID, message.Author.ID, fileTime)
			}
			dg.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Download of <%s> finished", link))
		} else if strings.ToLower(message.Content) == "cancel" {
			delete(interactiveChannelLinkTemp, message.ChannelID)
			dg.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Download of <%s> cancelled", link))
		} else if IsValid(message.Content) {
			dg.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Download of <%s> started", link))
			dg.ChannelTyping(message.ChannelID)
			delete(interactiveChannelLinkTemp, message.ChannelID)
			links := getDownloadLinks(link, message.ChannelID, true)
			for linkR, filename := range links {
				startDownload(linkR, filename, message.Content, message.ChannelID, message.Author.ID, fileTime)
			}
			dg.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Download of <%s> finished", link))
		} else {
			dg.ChannelMessageSend(message.ChannelID, "invalid path")
		}
	} else {
		_ = folderName
		foundLinks := false
		for _, iAttachment := range message.Attachments {
			dg.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Where do you want to save <%s>?\nType **.** for default path or **cancel** to cancel the download %s", iAttachment.URL, folderName))
			interactiveChannelLinkTemp[message.ChannelID] = iAttachment.URL
			foundLinks = true
		}
		foundUrls := xurls.Strict.FindAllString(message.Content, -1)
		for _, iFoundUrl := range foundUrls {
			dg.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Where do you want to save <%s>?\nType **.** for default path or **cancel** to cancel the download %s", iFoundUrl, folderName))
			interactiveChannelLinkTemp[message.ChannelID] = iFoundUrl
			foundLinks = true
		}
		if foundLinks == false {
			dg.ChannelMessageSend(message.ChannelID, "unable to find valid link")
		}
	}
}
