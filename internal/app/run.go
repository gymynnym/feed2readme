package app

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gymynnym/feed2readme/internal/config"
	"github.com/gymynnym/feed2readme/internal/feed"
	"github.com/gymynnym/feed2readme/internal/markdown"
	"github.com/gymynnym/feed2readme/internal/model"
)

func Run(opts model.Options) error {
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	loghubPosts := feed.GetPosts(client, model.ProviderLogHub, cfg.Feed.LogHub, opts.Limit)
	zennPosts := feed.GetPosts(client, model.ProviderZenn, cfg.Feed.Zenn, opts.Limit)
	feedSection := markdown.BuildFeedSection(loghubPosts, zennPosts)

	existingMarkdown, err := markdown.Read(opts.MarkdownPath)
	if err != nil {
		return err
	}

	appendSection := strings.Join(feedSection, "\n")
	updatedMarkdown, err := markdown.ReplaceOrAppend(existingMarkdown, appendSection)
	if err != nil {
		return err
	}

	if err := os.WriteFile(opts.MarkdownPath, []byte(updatedMarkdown), 0o644); err != nil {
		return fmt.Errorf("write markdown %q: %w", opts.MarkdownPath, err)
	}

	return nil
}
