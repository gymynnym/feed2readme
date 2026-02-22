package model

import "time"

type Post struct {
	Title   string
	Link    string
	PubDate time.Time
}

type Options struct {
	ConfigPath   string
	MarkdownPath string
	Limit        int
}
