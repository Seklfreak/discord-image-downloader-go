package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/bwmarrin/discordgo"
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

func updateDiscordStatus() {
	dg.UpdateStatusComplex(discordgo.UpdateStatusData{
		Game: &discordgo.Game{
			Name: fmt.Sprintf("%d downloaded pictures", countDownloadedImages()),
			Type: discordgo.GameTypeWatching,
		},
		Status: "online",
	})
}

func Pagify(text string, delimiter string) []string {
	result := make([]string, 0)
	textParts := strings.Split(text, delimiter)
	currentOutputPart := ""
	for _, textPart := range textParts {
		if len(currentOutputPart)+len(textPart)+len(delimiter) <= 1992 {
			if len(currentOutputPart) > 0 || len(result) > 0 {
				currentOutputPart += delimiter + textPart
			} else {
				currentOutputPart += textPart
			}
		} else {
			result = append(result, currentOutputPart)
			currentOutputPart = ""
			if len(textPart) <= 1992 {
				currentOutputPart = textPart
			}
		}
	}
	if currentOutputPart != "" {
		result = append(result, currentOutputPart)
	}
	return result
}

func filepathExtension(filepath string) string {
	if strings.Contains(filepath, "?") {
		filepath = strings.Split(filepath, "?")[0]
	}
	filepath = path.Ext(filepath)
	return filepath
}
