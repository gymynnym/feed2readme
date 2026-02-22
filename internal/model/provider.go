package model

import "fmt"

type Provider int

const (
	ProviderLogHub Provider = iota
	ProviderZenn
)

func (p Provider) FeedURL(username string) string {
	switch p {
	case ProviderLogHub:
		return fmt.Sprintf("https://loghub.me/%s/feed", username)
	case ProviderZenn:
		return fmt.Sprintf("https://zenn.dev/%s/feed", username)
	default:
		return ""
	}
}

func (p Provider) Language() string {
	switch p {
	case ProviderLogHub:
		return "한국어"
	case ProviderZenn:
		return "日本語"
	default:
		return "unknown"
	}
}

func (p Provider) String() string {
	switch p {
	case ProviderLogHub:
		return "LogHub"
	case ProviderZenn:
		return "Zenn"
	default:
		return "unknown"
	}
}
