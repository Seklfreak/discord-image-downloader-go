# discord-image-downloader-go

This is a simple tool which downloads pictures posted in discord channels of your choice to a local folder. It handles various sources like twitter differently to make sure to download the best quality available. It is written in go and the code is very ugly.

## How to use?
When you run the tool for the first time it creates an `config.ini` file with example values. Edit these values and run the tool for a second time. It should connect to discords api and wait for new messages.

### Where do I get the channel id?
Open discord in your browser and go to the channel you want to monitor. In your adress bar should be an URL like `https://discordapp.com/channels/1234/5678`. The number after the last slash is the channel id, in this case `5678`.
