package config

import (
	"os"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "womm-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := `
[server]
port = 9090
host = "127.0.0.1"

[storage]
path = "/tmp/test.db"

[github]
default_token = "ghp_test123"
rate_limit_ttl = "2h"

[cache]
ttl = "30m"

[themes]
default = "cyberpunk"
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", cfg.Server.Host)
	}
	if cfg.Storage.Path != "/tmp/test.db" {
		t.Errorf("expected path /tmp/test.db, got %s", cfg.Storage.Path)
	}
	if cfg.GitHub.DefaultToken != "ghp_test123" {
		t.Errorf("expected token, got %s", cfg.GitHub.DefaultToken)
	}
	if cfg.Themes.Default != "cyberpunk" {
		t.Errorf("expected cyberpunk, got %s", cfg.Themes.Default)
	}
	if cfg.GitHub.RateLimitTTL != "2h" {
		t.Errorf("expected rate_limit_ttl 2h, got %s", cfg.GitHub.RateLimitTTL)
	}
	if cfg.Cache.TTL != "30m" {
		t.Errorf("expected cache ttl 30m, got %s", cfg.Cache.TTL)
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load("/nonexistent/womm.toml")
	if err != nil {
		t.Fatalf("should not error on missing file, got: %v", err)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Themes.Default != "pixel" {
		t.Errorf("expected default theme pixel, got %s", cfg.Themes.Default)
	}
}
