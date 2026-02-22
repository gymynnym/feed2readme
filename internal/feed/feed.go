package feed

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gymynnym/feed2readme/internal/model"
)

type rss struct {
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

func GetPosts(client *http.Client, provider model.Provider, usernames []string, limit int) []model.Post {
	combined := make([]model.Post, 0)

	for _, username := range usernames {
		username = strings.TrimSpace(username)
		if username == "" {
			continue
		}

		url := provider.FeedURL(username)
		feed, err := fetchRSSFeed(client, url)
		if err != nil {
			panic(fmt.Sprintf("failed to fetch feed: %v", err))
		}

		posts, err := parsePosts(feed, limit)
		if err != nil {
			panic(fmt.Sprintf("failed to parse feed: %v", err))
		}

		combined = append(combined, posts...)
	}

	sort.Slice(combined, func(i, j int) bool {
		if combined[i].PubDate.Equal(combined[j].PubDate) {
			return combined[i].Title < combined[j].Title
		}
		return combined[i].PubDate.After(combined[j].PubDate)
	})

	if len(combined) > limit {
		combined = combined[:limit]
	}

	return combined
}

func fetchRSSFeed(client *http.Client, url string) (*rss, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed rss
	if err := xml.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decode rss %s: %w", url, err)
	}

	return &parsed, nil
}

func parsePosts(feed *rss, limit int) ([]model.Post, error) {
	posts := make([]model.Post, 0, limit)
	for _, item := range feed.Channel.Items {
		title := normalizeText(item.Title)
		link := strings.TrimSpace(item.Link)
		if title == "" || link == "" {
			continue
		}

		pubDate, err := parsePubDate(item.PubDate)
		if err != nil {
			pubDate = time.Time{}
		}

		posts = append(posts, model.Post{
			Title:   title,
			Link:    link,
			PubDate: pubDate,
		})

		if len(posts) >= limit {
			break
		}
	}

	return posts, nil
}

func parsePubDate(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, errors.New("empty pubDate")
	}

	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, raw)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported pubDate format %q", raw)
}

func normalizeText(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	return strings.Join(strings.Fields(value), " ")
}
