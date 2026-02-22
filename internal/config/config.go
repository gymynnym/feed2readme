package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	Feed feedConfig `toml:"feed"`
}

type feedConfig struct {
	LogHub []string `toml:"loghub"`
	Zenn   []string `toml:"zenn"`
}

func Load(path string) (config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	cfg := config{}
	if err := toml.Unmarshal(content, &cfg); err != nil {
		return config{}, fmt.Errorf("parse config %q: %w", path, err)
	}

	return cfg, nil
}
