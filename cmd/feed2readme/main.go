package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gymynnym/feed2readme/internal/app"
	"github.com/gymynnym/feed2readme/internal/model"
)

const (
	defaultConfigPath   = "feed2readme.toml"
	defaultMarkdownPath = "README.md"
	defaultLimit        = 5
)

func main() {
	configPath := flag.String("c", defaultConfigPath, "path to TOML config file")
	markdownPath := flag.String("m", defaultMarkdownPath, "path to markdown file")
	limit := flag.Int("l", defaultLimit, "max number of posts per provider")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "error: config path is required")
		os.Exit(1)
	}
	if *markdownPath == "" {
		fmt.Fprintln(os.Stderr, "error: markdown path is required")
		os.Exit(1)
	}
	if *limit < 1 {
		fmt.Fprintln(os.Stderr, "error: limit must be 1 or greater")
		os.Exit(1)
	}

	opts := model.Options{
		ConfigPath:   *configPath,
		MarkdownPath: *markdownPath,
		Limit:        *limit,
	}

	if err := app.Run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
