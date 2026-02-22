package markdown

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gymynnym/feed2readme/internal/model"
)

const (
	feedStartMarker = "<!-- FEED START -->"
	feedEndMarker   = "<!-- FEED END -->"
)

func BuildFeedSection(loghubPosts, zennPosts []model.Post) []string {
	lines := make([]string, 0, 32)
	lines = append(lines, feedStartMarker)

	lines = append(lines, "## Latest Posts")

	if len(loghubPosts) > 0 {
		loghubTable := buildFeedTable(loghubPosts, model.ProviderLogHub)
		lines = append(lines, "")
		lines = append(lines, loghubTable)
	}

	if len(zennPosts) > 0 {
		zennTable := buildFeedTable(zennPosts, model.ProviderZenn)
		lines = append(lines, "")
		lines = append(lines, zennTable)
	}

	lines = append(lines, feedEndMarker)
	return lines
}

func buildFeedTable(posts []model.Post, provider model.Provider) string {
	lines := make([]string, 0, len(posts)+2)

	heading := fmt.Sprintf("|%s (%s)|date|", provider.String(), provider.Language())
	lines = append(lines, heading)
	lines = append(lines, "|:---|:---|")

	for _, post := range posts {
		column := fmt.Sprintf("|[%s](%s)|**`%s`**|", post.Title, post.Link, post.PubDate.Format("2006-01-02"))
		lines = append(lines, column)
	}

	return strings.Join(lines, "\n")
}

func Read(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("read markdown %q: %w", path, err)
	}
	return string(content), nil
}

func ReplaceOrAppend(markdown, feedBlock string) (string, error) {
	start := strings.Index(markdown, feedStartMarker)
	end := strings.Index(markdown, feedEndMarker)

	switch {
	case start == -1 && end == -1:
		trimmed := strings.TrimRight(markdown, "\n")
		if trimmed == "" {
			return feedBlock + "\n", nil
		}
		return trimmed + "\n\n" + feedBlock + "\n", nil
	case start == -1 || end == -1:
		return "", errors.New("found only one feed marker; both FEED START and FEED END are required")
	case end < start:
		return "", errors.New("FEED END appears before FEED START")
	}

	before := strings.TrimRight(markdown[:start], "\n")
	after := strings.TrimLeft(markdown[end+len(feedEndMarker):], "\n")

	var out strings.Builder
	if before != "" {
		out.WriteString(before)
		out.WriteString("\n\n")
	}
	out.WriteString(feedBlock)
	if after != "" {
		out.WriteString("\n\n")
		out.WriteString(after)
	}

	result := out.String()
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return result, nil
}
