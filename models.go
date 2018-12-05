package main

import "time"

type DownloadItem struct {
	Link     string
	Filename string
	Time     time.Time
}
