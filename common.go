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

// deduplicateDownloadItems removes duplicates from a slice of *DownloadItem s identified by the Link
func deduplicateDownloadItems(DownloadItems []*DownloadItem) []*DownloadItem {
	var result []*DownloadItem
	seen := map[string]bool{}

	for _, item := range DownloadItems {
		if seen[item.Link] {
			continue
		}

		seen[item.Link] = true
		result = append(result, item)
	}

	return result
}
