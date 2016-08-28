# discord-image-downloader-go

This is a simple tool which downloads pictures posted in discord channels of your choice to a local folder. The tool only handles links posted while it is running, you have to keep it running the whole time. It handles various sources like twitter differently to make sure to download the best quality available. It is written in go and the code is very ugly.

## How to use?
When you run the tool for the first time it creates an `config.ini` file with example values. Edit these values and run the tool for a second time. It should connect to discords api and wait for new messages.

In case you are using two-factor authentication you have to login using your token. Remove the the email and password lines under the auth section in the config file and instead put in `token = <your token>`. You can acquire your token from the developer tools in your browser (`localStorage.token`).

### Where do I get the channel id?
Open discord in your browser and go to the channel you want to monitor. In your adress bar should be an URL like `https://discordapp.com/channels/1234/5678`. The number after the last slash is the channel id, in this case `5678`.
