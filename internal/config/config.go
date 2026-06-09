package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server  ServerConfig  `toml:"server"`
	Storage StorageConfig `toml:"storage"`
	GitHub  GitHubConfig  `toml:"github"`
	Cache   CacheConfig   `toml:"cache"`
	Themes  ThemesConfig  `toml:"themes"`
}

type ServerConfig struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type StorageConfig struct {
	Path string `toml:"path"`
}

type GitHubConfig struct {
	DefaultToken string `toml:"default_token"`
	RateLimitTTL string `toml:"rate_limit_ttl"`
}

type CacheConfig struct {
	TTL string `toml:"ttl"`
}

type ThemesConfig struct {
	Default string `toml:"default"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Server:  ServerConfig{Port: 8080, Host: "0.0.0.0"},
		Storage: StorageConfig{Path: "womm.db"},
		GitHub:  GitHubConfig{RateLimitTTL: "1h"},
		Cache:   CacheConfig{TTL: "1h"},
		Themes:  ThemesConfig{Default: "pixel"},
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, nil
	}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
