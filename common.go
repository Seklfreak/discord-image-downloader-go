package main

import (
	"strings"
)

// isDiscordEmoji matches https://cdn.discordapp.com/emojis/503141595860959243.gif, and similar URLs/filenames
func isDiscordEmoji(link string) bool {
	// always match discord emoji URLs, eg https://cdn.discordapp.com/emojis/340989430460317707.png
	if strings.HasPrefix(link, discordEmojiBaseUrl) {
		return true
	}

	return false
}
