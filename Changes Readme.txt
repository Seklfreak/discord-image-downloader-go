//This same piece of code can be found at https://pastebin.com/MBhFQfHQ
//Edited by EnanoFurtivo: https://github.com/EnanoFurtivo/discord-image-downloader-go
//Sourcecode by Selfreak: https://github.com/Seklfreak/discord-image-downloader-go

//Edited file main.go in 'Case == history'

//this code adds the functinality to query more than one channel id for history donwload at a time

//IMPORTANT: in the main.go that is uploaded in my github i had to delete a bunch of code due to a to a problem with mvdan.cc library. Still i was able to compile and test the functionality of the added feature.
//Something to have in mind is that the piece of code below is complete and nothing was substracted from it only it was added with the for function wich allows this modification to work (nothing else is needed).
//If you replace this into the original main.go from Selfreak and you manage to fully compile, all of the original functionalities should still be working just fine plus the added for functionality.

case message == "history", historyCommandIsActive:
				i := 0
				_, historyCommandIsSet := historyCommandActive[m.ChannelID]
				if !historyCommandIsSet || historyCommandActive[m.ChannelID] == "" {
					historyCommandActive[m.ChannelID] = ""
					
					//dg.ChannelMessageSend(m.ChannelID, fmt.Sprintf(message))
					//dg.ChannelMessageSend(m.ChannelID, fmt.Sprintf(array[x]))
					
					idArray := strings.Split(message, ", ")
					for index, chanelValue := range idArray {
						fmt.Sprintf(chanelValue,index)
						//dg.ChannelMessageSend(m.ChannelID, fmt.Sprintf(chanelValue,index))
						if folder, ok := ChannelWhitelist[chanelValue]; ok {
							dg.ChannelMessageSend(m.ChannelID, fmt.Sprintf("downloading to `%s`", folder))
							historyCommandActive[m.ChannelID] = "downloading"
							lastBefore := ""
							var lastBeforeTime time.Time
						MessageRequestingLoop:
							for true {
								if lastBeforeTime != (time.Time{}) {
									fmt.Printf("[%s] Requesting 100 more messages, (before %s)\n", time.Now().Format(time.Stamp), lastBeforeTime)
									dg.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Requesting 100 more messages, (before %s)\n", lastBeforeTime))
								}
								messages, err := dg.ChannelMessages(chanelValue, 100, lastBefore, "", "")
								if err == nil {
									if len(messages) <= 0 {
										delete(historyCommandActive, m.ChannelID)
										break MessageRequestingLoop
									}
									lastBefore = messages[len(messages)-1].ID
									lastBeforeTime, err = messages[len(messages)-1].Timestamp.Parse()
									if err != nil {
										fmt.Println(err)
									}
									for _, message := range messages {
										fileTime := time.Now()
										if m.Timestamp != "" {
											fileTime, err = message.Timestamp.Parse()
											if err != nil {
												fmt.Println(err)
											}
										}
										if historyCommandActive[m.ChannelID] == "cancel" {
											delete(historyCommandActive, m.ChannelID)
											break MessageRequestingLoop
										}
										for _, iAttachment := range message.Attachments {
											if len(findDownloadedImageByUrl(iAttachment.URL)) == 0 {
												i++
												startDownload(iAttachment.URL, iAttachment.Filename, folder, message.ChannelID, message.Author.ID, fileTime)
											}
										}
									}
								} else {
									dg.ChannelMessageSend(m.ChannelID, err.Error())
									fmt.Println(err)
									delete(historyCommandActive, m.ChannelID)
									break MessageRequestingLoop
								}
							}
							dg.ChannelMessageSend(m.ChannelID, fmt.Sprintf("done, %d download links started!", i))
						} else {
							dg.ChannelMessageSend(m.ChannelID, "please send me a channel id or various channel ID's splitted by commas ej.: 426332433793941508, 926832233496945508, ... (from the whitelist)")
						}
					}
				} else if historyCommandActive[m.ChannelID] == "downloading" && message == "cancel" {
					historyCommandActive[m.ChannelID] = "cancel"
				}
