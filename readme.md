# PROJECT ABANDONED, [PLEASE USE FORK.](https://github.com/get-got/discord-downloader-go/)

---

# discord-image-downloader-go
[<img src="https://img.shields.io/badge/Support-me!-orange.svg">](https://www.paypal.me/swk)
[![Go Report Card](https://goreportcard.com/badge/github.com/Seklfreak/discord-image-downloader-go)](https://goreportcard.com/report/github.com/Seklfreak/discord-image-downloader-go)
[![Build Status](https://travis-ci.com/Seklfreak/discord-image-downloader-go.svg?branch=master)](https://travis-ci.com/Seklfreak/discord-image-downloader-go)

[**DOWNLOAD THE LATEST RELEASE BUILD**](https://github.com/Seklfreak/discord-image-downloader-go/releases/latest)

This project is not often maintained. For an actively maintained fork that implements features such as extensive JSON settings with channel-specific configurations, see [**get-got/discord-downloader-go**](https://github.com/get-got/discord-downloader-go)

## Discord SelfBots are forbidden!
[Official Statement](https://support.discordapp.com/hc/en-us/articles/115002192352-Automated-user-accounts-self-bots-)
### You have been warned.

This is a simple tool which downloads media posted in Discord channels of your choice to a local folder. It handles various sources like Twitter differently to make sure to download the best quality available.

## Websites currently supported
- Discord Attachments
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

If you are using a normal user account **without two-factor authentication (2FA)**, simply enter your email and password into the corresponding lines in `config.ini`, under the auth section.

If you are using **two-factor authentication (2FA) you have to login using your token.** Remove the email and password lines under the auth section in the config file and instead put in `token = <your token>`. You can acquire your token from the developer tools in your browser (`localStorage.token`) or discord client (`Ctrl+Shift+I` (Windows) or `Cmd+Option+I` (Mac), click Application, click Local Storage, click `https://discordapp.com`, and find "token" and paste the value).

If you wish to use a **bot account (not a user account)**, go to https://discord.com/developers/applications and create an application, then create a bot in the `Bot` tab in application settings. The bot tab will show you your token. You can invite to your server(s) by going to the `OAuth2` tab in application settings, check `bot`, and copy+paste the url into your browser. **In the `config.ini`, add "Bot " before your token. (example: `token = Bot mytokenhere`)**

## How to download old files?
By default, the tool only downloads new links posted while the tool is running. You can also set up the tool to download the complete history of a channel. To do this you have to run this tool with a separate discord account. Send your second account a dm on your primary account and get the channel id from the direct message channel. Now add this channel id to the config by adding the following lines:
```
[interactive channels]
<your channel id> = <some valid path>
```
After this is done restart the tool and send `history` as a DM to your second account. The bot will ask for the channel id of the channel you want to download and start the downloads. You can view all available commands by sending `help`.

### Where do I get the Channel ID?
Enable Developer Mode (in Discord Appearance settings) and right click the channel you need, and click Copy ID.

**OR,** Open discord in your browser and go to the channel you want to monitor. In your address bar should be a URL like `https://discordapp.com/channels/1234/5678`. The number after the last slash is the channel ID, in this case, `5678`.

### Where do I get the Channel ID for Direct Messages?
1. Inspect Element in the Discord client (`Ctrl+Shift+I` for Windows or `Cmd+Option+I` for Mac)
1. Go to the `Elements` tab on the left.
1. Click this icon ![arrow going into box](https://i.imgur.com/PkDOCyZ.png) (the arrow going into a box) and then click on the avatar for the persons DMs you want to grab the ID for.
1. Somewhere slightly above the HTML it takes you to, there should be a line that looks like this ![](https://i.imgur.com/614rZnX.png)
1. Copy the number after the `/@me/`. That is your Channel ID to use.

**OR,** Open discord in your browser and go to the channel you want to monitor. In your address bar should be a URL like `https://discordapp.com/channels/@me/5678`. The number after the last slash is the channel ID, in this case, `5678`.
