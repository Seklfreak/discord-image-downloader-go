# discord-image-downloader-go
[<img src="https://img.shields.io/badge/Support-me!-orange.svg">](https://www.paypal.me/swk) [![Go Report Card](https://goreportcard.com/badge/github.com/Seklfreak/discord-image-downloader-go)](https://goreportcard.com/report/github.com/Seklfreak/discord-image-downloader-go) [![Build Status](https://travis-ci.org/Seklfreak/discord-image-downloader-go.svg?branch=master)](https://travis-ci.org/Seklfreak/discord-image-downloader-go)

[Download the latest release](https://github.com/Seklfreak/discord-image-downloader-go/releases/latest)

## Discord SelfBots are forbidden!
[Official Statement](https://support.discordapp.com/hc/en-us/articles/115002192352-Automated-user-accounts-self-bots-)
### You have been warned.

This is a simple tool which downloads pictures (and instagram videos) posted in discord channels of your choice to a local folder. It handles various sources like twitter differently to make sure to download the best quality available.

## Websites currently supported
- Discord attachments
- Twitter
- Tistory
- Gfycat
- Instagram
- Imgur
- Google Drive Files and Folders
- Flickr
- Streamable
- Any direct link to an image or video

## How to use?
When you run the tool for the first time it creates a `config.ini` file with example values. Edit these values and run the tool for a second time. It should connect to discords api and wait for new messages.

In case you are using two-factor authentication you have to login using your token. Remove the email and password lines under the auth section in the config file and instead put in `token = <your token>`. You can acquire your token from the developer tools in your browser (`localStorage.token`) or discord client (Control+Shift+I (Windows) or Command+Option+I, click Application, click Local Storage, click `https://discordapp.com`, and find "token" and paste the value).

## How to download old pictures?
By default, the tool only downloads new links posted while the tool is running. You can also set up the tool to download the complete history of a channel. To do this you have to run this tool with a separate discord account. Send your second account a dm on your primary account and get the channel id from the direct message channel. Now add this channel id to the config by adding the following lines:
```
[interactive channels]
<your channel id> = <some valid path>
```
After this is done restart the tool and send `history` as a DM to your second account. The bot will ask for the channel id of the channel you want to download and start the downloads. You can view all available commands by sending `help`.

### Where do I get the channel id?
Open discord in your browser and go to the channel you want to monitor. In your address bar should be a URL like `https://discordapp.com/channels/1234/5678`. The number after the last slash is the channel id, in this case, `5678`. Or, enable Developer mode and right click the channel you need, and click Copy ID.
