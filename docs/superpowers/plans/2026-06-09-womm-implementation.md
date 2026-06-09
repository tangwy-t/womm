# WOMM Badge Generator Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a self-hosted GitHub achievement badge service with 25 sarcastic/humorous badges, SVG generation in 4 themes, GitHub API-based certification, HTTP API, and CLI — all in a single Go binary.

**Architecture:** Monolithic Go service with internal module separation: Badge Registry (definitions), Render Engine (SVG generation via Go templates), Certify Engine (GitHub API validation), Store (SQLite persistence), HTTP Server (chi router), and CLI (cobra).

**Tech Stack:** Go 1.22+, go-chi/chi (HTTP router), spf13/cobra (CLI), google/go-github (GitHub API), modernc.org/sqlite (pure Go SQLite), BurntSushi/toml (config), html/template (SVG rendering).

---

## File Structure

```
/data/womm/
├── main.go
├── go.mod
├── womm.toml
├── Dockerfile
├── cmd/
│   └── root.go                       # Cobra root + all subcommands
├── internal/
│   ├── app/
│   │   └── app.go                    # Wire all modules together
│   ├── config/
│   │   ├── config.go                 # Config struct + TOML loading
│   │   └── config_test.go
│   ├── badge/
│   │   ├── types.go                  # Badge, BadgeType, Rarity
│   │   ├── registry.go              # Registry: lookup, list, i18n, RegisterAll
│   │   ├── registry_test.go
│   │   ├── declarative.go           # 10 declarative badge defs
│   │   └── certified.go             # 15 certified/legendary badge defs
│   ├── store/
│   │   ├── store.go                  # SQLite schema + CRUD
│   │   └── store_test.go
│   ├── certify/
│   │   ├── github.go                 # GitHubClient interface + domain types
│   │   ├── mock.go                   # Mock for testing
│   │   ├── functions.go             # 15 certify* functions + dispatch map
│   │   ├── functions_test.go
│   │   ├── engine.go                 # Certify engine with cache
│   │   └── engine_test.go
│   ├── render/
│   │   ├── render.go                 # Renderer: compose badge+theme+template->SVG
│   │   ├── theme.go                  # 4 theme definitions
│   │   ├── icons.go                  # SVG icon path registry
│   │   ├── render_test.go
│   │   └── templates/
│   │       ├── badge.svg.tmpl
│   │       ├── wide.svg.tmpl
│   │       ├── terminal.svg.tmpl
│   │       └── stamp.svg.tmpl
│   └── server/
│       ├── server.go                 # Chi router + middleware
│       ├── handler.go                # HTTP handlers
│       └── handler_test.go
└── docs/superpowers/...
```

---

## Task 1: Project Scaffolding & Config

**Files:**
- Create: `go.mod`, `main.go`, `womm.toml`
- Create: `cmd/root.go`
- Create: `internal/config/config.go`, `internal/config/config_test.go`

- [ ] **Step 1: Initialize Go module**

```bash
cd /data/womm
git init
go mod init github.com/womm/womm
go get github.com/go-chi/chi/v5@latest
go get github.com/spf13/cobra@latest
go get github.com/google/go-github/v68@latest
go get github.com/BurntSushi/toml@latest
go get modernc.org/sqlite@latest
go get golang.org/x/oauth2@latest
```

- [ ] **Step 1b: Create `.gitignore`**

```
womm
womm.db
*.db
.superpowers/
```

- [ ] **Step 2: Create `main.go`**

```go
package main

import "github.com/womm/womm/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 3: Create `cmd/root.go`**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "womm",
	Short: "WOMM - Works On My Machine badge generator",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

- [ ] **Step 4: Create `womm.toml`**

```toml
[server]
port = 8080
host = "0.0.0.0"

[storage]
path = "womm.db"

[github]
default_token = ""
rate_limit_ttl = "1h"

[cache]
ttl = "1h"

[themes]
default = "pixel"
```

- [ ] **Step 5: Write config test in `internal/config/config_test.go`**

```go
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
```

- [ ] **Step 6: Run test, verify FAIL**

```bash
go test ./internal/config/...
```

Expected: FAIL — `Load` undefined.

- [ ] **Step 7: Create `internal/config/config.go`**

```go
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
```

- [ ] **Step 8: Run tests, verify PASS**

```bash
go test ./internal/config/... -v
```

Expected: PASS (2 tests).

- [ ] **Step 9: Verify build**

```bash
go build -o womm .
./womm --help
```

Expected: Shows usage text.

- [ ] **Step 10: Commit**

```bash
git add -A && git commit -m "feat: project scaffolding, config, cobra cli"
```

---

## Task 2: Badge Types & Registry (25 Badges)

**Files:**
- Create: `internal/badge/types.go`
- Create: `internal/badge/registry.go`
- Create: `internal/badge/registry_test.go`
- Create: `internal/badge/declarative.go`
- Create: `internal/badge/certified.go`

- [ ] **Step 1: Write registry tests**

```go
// internal/badge/registry_test.go
package badge

import "testing"

func TestLookup(t *testing.T) {
	reg := NewRegistry()
	RegisterAll(reg)

	b, ok := reg.Lookup("midnight-coder")
	if !ok {
		t.Fatal("expected midnight-coder to exist")
	}
	if b.Type != Certified {
		t.Errorf("expected Certified, got %v", b.Type)
	}
	if b.Rarity != Rare {
		t.Errorf("expected Rare, got %v", b.Rarity)
	}
}

func TestLookupNotFound(t *testing.T) {
	reg := NewRegistry()
	_, ok := reg.Lookup("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestListAll(t *testing.T) {
	reg := NewRegistry()
	RegisterAll(reg)
	all := reg.ListAll()
	if len(all) != 25 {
		t.Errorf("expected 25 badges, got %d", len(all))
	}
}

func TestListByType(t *testing.T) {
	reg := NewRegistry()
	RegisterAll(reg)
	d := reg.ListByType(Declarative)
	if len(d) != 10 {
		t.Errorf("expected 10 declarative, got %d", len(d))
	}
	c := reg.ListByType(Certified)
	if len(c) != 15 {
		t.Errorf("expected 15 certified, got %d", len(c))
	}
}

func TestI18n(t *testing.T) {
	reg := NewRegistry()
	RegisterAll(reg)
	b, _ := reg.Lookup("midnight-coder")
	if b.LocalizedName("zh") != "午夜编码者" {
		t.Errorf("wrong zh name: %s", b.LocalizedName("zh"))
	}
	if b.LocalizedName("en") != "Midnight Coder" {
		t.Errorf("wrong en name: %s", b.LocalizedName("en"))
	}
}

func TestI18nFallback(t *testing.T) {
	reg := NewRegistry()
	RegisterAll(reg)
	b, _ := reg.Lookup("midnight-coder")
	if b.LocalizedName("fr") != "午夜编码者" {
		t.Errorf("expected zh fallback for unknown lang, got: %s", b.LocalizedName("fr"))
	}
}
```

- [ ] **Step 2: Run test, verify FAIL**

```bash
go test ./internal/badge/...
```

- [ ] **Step 3: Create `internal/badge/types.go`**

```go
package badge

type BadgeType string

const (
	Declarative BadgeType = "declarative"
	Certified   BadgeType = "certified"
)

type Rarity string

const (
	Common    Rarity = "common"
	Rare      Rarity = "rare"
	Legendary Rarity = "legendary"
)

type Badge struct {
	ID       string            `json:"id"`
	Name     map[string]string `json:"name"`
	Subtitle map[string]string `json:"subtitle"`
	Icon     string            `json:"icon"`
	Type     BadgeType         `json:"type"`
	Rarity   Rarity            `json:"rarity"`
}

func (b *Badge) LocalizedName(lang string) string {
	if name, ok := b.Name[lang]; ok {
		return name
	}
	return b.Name["zh"]
}

func (b *Badge) LocalizedSubtitle(lang string) string {
	if sub, ok := b.Subtitle[lang]; ok {
		return sub
	}
	return b.Subtitle["zh"]
}
```

- [ ] **Step 4: Create `internal/badge/registry.go`**

```go
package badge

import "sync"

type Registry struct {
	mu     sync.RWMutex
	badges map[string]*Badge
}

func NewRegistry() *Registry {
	return &Registry{badges: make(map[string]*Badge)}
}

func (r *Registry) Register(b *Badge) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.badges[b.ID] = b
}

func (r *Registry) Lookup(id string) (*Badge, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.badges[id]
	return b, ok
}

func (r *Registry) ListAll() []*Badge {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Badge, 0, len(r.badges))
	for _, b := range r.badges {
		result = append(result, b)
	}
	return result
}

func (r *Registry) ListByType(bt BadgeType) []*Badge {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Badge
	for _, b := range r.badges {
		if b.Type == bt {
			result = append(result, b)
		}
	}
	return result
}

func RegisterAll(r *Registry) {
	for _, b := range declarativeBadges {
		r.Register(b)
	}
	for _, b := range certifiedBadges {
		r.Register(b)
	}
}
```

- [ ] **Step 5: Create `internal/badge/declarative.go`**

All 10 declarative badges with bilingual Name/Subtitle:

```go
package badge

var declarativeBadges = []*Badge{
	{ID: "works-on-my-machine", Name: map[string]string{"zh": "在我机器上能运行", "en": "Works On My Machine"}, Subtitle: map[string]string{"zh": "态度即正义", "en": "Attitude is everything"}, Icon: "checkmark", Type: Declarative, Rarity: Common},
	{ID: "read-not-reply", Name: map[string]string{"zh": "已读不回", "en": "Read Not Reply"}, Subtitle: map[string]string{"zh": "Review了你的PR，然后…没有然后了", "en": "Reviewed your PR, then... nothing"}, Icon: "eye", Type: Declarative, Rarity: Common},
	{ID: "stackoverflow-courier", Name: map[string]string{"zh": "Stack Overflow搬运工", "en": "Stack Overflow Courier"}, Subtitle: map[string]string{"zh": "代码从网上来，到网上去", "en": "Code comes from the web, returns to the web"}, Icon: "stack", Type: Declarative, Rarity: Common},
	{ID: "todo-collector", Name: map[string]string{"zh": "TODO收藏家", "en": "TODO Collector"}, Subtitle: map[string]string{"zh": "// TODO: fix this later × 50", "en": "// TODO: fix this later × 50"}, Icon: "list", Type: Declarative, Rarity: Common},
	{ID: "comment-fundamentalist", Name: map[string]string{"zh": "注释原教旨主义者", "en": "Comment Fundamentalist"}, Subtitle: map[string]string{"zh": "每行代码配三行注释，包括i++", "en": "Three comments per line, including i++"}, Icon: "hash", Type: Declarative, Rarity: Common},
	{ID: "copy-paste-engineer", Name: map[string]string{"zh": "复制粘贴工程师", "en": "Copy Paste Engineer"}, Subtitle: map[string]string{"zh": "Ctrl+C / Ctrl+V 是核心技能", "en": "Ctrl+C / Ctrl+V is my core skill"}, Icon: "clipboard", Type: Declarative, Rarity: Common},
	{ID: "rubber-duck-master", Name: map[string]string{"zh": "橡皮鸭调试大师", "en": "Rubber Duck Master"}, Subtitle: map[string]string{"zh": "对着鸭子说话就能修bug", "en": "Talk to a duck, fix every bug"}, Icon: "duck", Type: Declarative, Rarity: Rare},
	{ID: "no-friday-deploy", Name: map[string]string{"zh": "周五不部署", "en": "No Friday Deploy"}, Subtitle: map[string]string{"zh": "血的教训换来的铁律", "en": "An iron rule forged in blood"}, Icon: "calendar", Type: Declarative, Rarity: Rare},
	{ID: "force-push-warrior", Name: map[string]string{"zh": "Git Force Push勇士", "en": "Force Push Warrior"}, Subtitle: map[string]string{"zh": "--force 是我的日常", "en": "--force is my daily routine"}, Icon: "zap", Type: Declarative, Rarity: Rare},
	{ID: "meeting-survivor", Name: map[string]string{"zh": "会议幸存者", "en": "Meeting Survivor"}, Subtitle: map[string]string{"zh": "今天开了6个会，写了0行代码", "en": "6 meetings today, 0 lines of code"}, Icon: "users", Type: Declarative, Rarity: Common},
}
```

- [ ] **Step 6: Create `internal/badge/certified.go`**

All 15 certified/legendary badges:

```go
package badge

var certifiedBadges = []*Badge{
	{ID: "midnight-coder", Name: map[string]string{"zh": "午夜编码者", "en": "Midnight Coder"}, Subtitle: map[string]string{"zh": "月亮不睡我不睡", "en": "The moon doesn't sleep, neither do I"}, Icon: "moon", Type: Certified, Rarity: Rare},
	{ID: "weekend-warrior", Name: map[string]string{"zh": "周末战士", "en": "Weekend Warrior"}, Subtitle: map[string]string{"zh": "工作使我快乐（周末也是）", "en": "Work makes me happy (weekends too)"}, Icon: "sun", Type: Certified, Rarity: Rare},
	{ID: "issue-lord", Name: map[string]string{"zh": "百Issue之主", "en": "Issue Lord"}, Subtitle: map[string]string{"zh": "一切安好……大概", "en": "Everything is fine... probably"}, Icon: "alert", Type: Certified, Rarity: Rare},
	{ID: "docs-master", Name: map[string]string{"zh": "文档仙人", "en": "Docs Master"}, Subtitle: map[string]string{"zh": "代码没写几行，文档写了一本小说", "en": "Few lines of code, a novel of docs"}, Icon: "book", Type: Certified, Rarity: Rare},
	{ID: "pr-bomber", Name: map[string]string{"zh": "PR轰炸机", "en": "PR Bomber"}, Subtitle: map[string]string{"zh": "一天一个PR，医生远离我", "en": "A PR a day keeps the doctor away"}, Icon: "rocket", Type: Certified, Rarity: Rare},
	{ID: "monkey-wrench", Name: map[string]string{"zh": "猴子扳手", "en": "Monkey Wrench"}, Subtitle: map[string]string{"zh": "我来了，CI挂了", "en": "I arrived, CI broke"}, Icon: "wrench", Type: Certified, Rarity: Rare},
	{ID: "archaeologist", Name: map[string]string{"zh": "考古学家", "en": "Archaeologist"}, Subtitle: map[string]string{"zh": "挖出了上古代码", "en": "Unearthed ancient code"}, Icon: "pickaxe", Type: Certified, Rarity: Legendary},
	{ID: "branch-hoarder", Name: map[string]string{"zh": "分支囤积者", "en": "Branch Hoarder"}, Subtitle: map[string]string{"zh": "每个分支都是'马上要合并的'", "en": "Every branch is 'about to be merged'"}, Icon: "git-branch", Type: Certified, Rarity: Rare},
	{ID: "ghost-committer", Name: map[string]string{"zh": "幽灵提交者", "en": "Ghost Committer"}, Subtitle: map[string]string{"zh": "我还活着，只是不想写代码", "en": "I'm alive, just don't want to code"}, Icon: "ghost", Type: Certified, Rarity: Legendary},
	{ID: "polyglot", Name: map[string]string{"zh": "多语言通才", "en": "Polyglot"}, Subtitle: map[string]string{"zh": "什么都会一点，什么都不精", "en": "Jack of all trades, master of none"}, Icon: "globe", Type: Certified, Rarity: Rare},
	{ID: "true-destroyer", Name: map[string]string{"zh": "真·破坏王", "en": "True Destroyer"}, Subtitle: map[string]string{"zh": "连续3次搞挂CI", "en": "Broke CI 3 times in a row"}, Icon: "skull", Type: Certified, Rarity: Legendary},
	{ID: "y2k-hunter", Name: map[string]string{"zh": "千年虫猎人", "en": "Y2K Hunter"}, Subtitle: map[string]string{"zh": "还在跟1999年的代码打交道", "en": "Still dealing with code from 1999"}, Icon: "calendar", Type: Certified, Rarity: Legendary},
	{ID: "life-404", Name: map[string]string{"zh": "404人生", "en": "Life 404"}, Subtitle: map[string]string{"zh": "个人主页？不存在的", "en": "Profile page? Not found"}, Icon: "404", Type: Certified, Rarity: Legendary},
	{ID: "commit-anniversary", Name: map[string]string{"zh": "首次提交纪念日", "en": "Commit Anniversary"}, Subtitle: map[string]string{"zh": "写代码这么多年了", "en": "Years of writing code"}, Icon: "trophy", Type: Certified, Rarity: Legendary},
	{ID: "fullstack-victim", Name: map[string]string{"zh": "全栈受害者", "en": "Fullstack Victim"}, Subtitle: map[string]string{"zh": "前端后端运维全都要我干", "en": "Frontend, backend, DevOps — all on me"}, Icon: "layers", Type: Certified, Rarity: Legendary},
}
```

- [ ] **Step 7: Run tests, verify all PASS**

```bash
go test ./internal/badge/... -v
```

Expected: 6 tests PASS.

- [ ] **Step 8: Commit**

```bash
git add internal/badge/ && git commit -m "feat: badge types, registry, 25 badges (declarative + certified)"
```

---

## Task 3: SQLite Store

**Files:**
- Create: `internal/store/store.go`
- Create: `internal/store/store_test.go`

- [ ] **Step 1: Write store tests**

```go
// internal/store/store_test.go
package store

import (
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	s, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
}

func TestClaimBadge(t *testing.T) {
	s, _ := Open(":memory:")
	defer s.Close()

	if err := s.ClaimBadge("torvalds", "works-on-my-machine"); err != nil {
		t.Fatal(err)
	}
	ok, err := s.IsUnlocked("torvalds", "works-on-my-machine")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected unlocked")
	}
}

func TestCertifyBadge(t *testing.T) {
	s, _ := Open(":memory:")
	defer s.Close()

	if err := s.CertifyBadge("torvalds", "midnight-coder"); err != nil {
		t.Fatal(err)
	}
	states, err := s.GetUserBadges("torvalds")
	if err != nil {
		t.Fatal(err)
	}
	if len(states) != 1 {
		t.Errorf("expected 1, got %d", len(states))
	}
	if states[0].Source != "certified" {
		t.Errorf("expected certified, got %s", states[0].Source)
	}
}

func TestIsUnlocked_False(t *testing.T) {
	s, _ := Open(":memory:")
	defer s.Close()
	ok, _ := s.IsUnlocked("nobody", "midnight-coder")
	if ok {
		t.Error("expected false")
	}
}

func TestCertCache_SetGet(t *testing.T) {
	s, _ := Open(":memory:")
	defer s.Close()

	err := s.SetCertCache("torvalds", "midnight-coder", true, `{"ratio":0.4}`, 1*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	entry, ok, err := s.GetCertCache("torvalds", "midnight-coder")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected cache hit")
	}
	if !entry.Result {
		t.Error("expected result true")
	}
}

func TestCertCache_Expired(t *testing.T) {
	s, _ := Open(":memory:")
	defer s.Close()

	s.SetCertCache("u", "b", true, "", -1*time.Hour)
	_, ok, _ := s.GetCertCache("u", "b")
	if ok {
		t.Error("expected expired cache to return false")
	}
}
```

- [ ] **Step 2: Run test, verify FAIL**

```bash
go test ./internal/store/...
```

- [ ] **Step 3: Create `internal/store/store.go`**

```go
package store

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct{ db *sql.DB }

type BadgeState struct {
	GitHubUser string     `json:"github_user"`
	BadgeID    string     `json:"badge_id"`
	Unlocked   bool       `json:"unlocked"`
	UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
	Source     string     `json:"source"`
}

type CacheEntry struct {
	GitHubUser string
	BadgeID    string
	Result     bool
	RawData    string
	ExpiresAt  time.Time
}

const schema = `
CREATE TABLE IF NOT EXISTS badge_state (
    github_user TEXT NOT NULL,
    badge_id    TEXT NOT NULL,
    unlocked    INTEGER NOT NULL DEFAULT 0,
    unlocked_at TEXT,
    source      TEXT,
    PRIMARY KEY (github_user, badge_id)
);
CREATE TABLE IF NOT EXISTS cert_cache (
    github_user TEXT NOT NULL,
    badge_id    TEXT NOT NULL,
    result      INTEGER NOT NULL,
    raw_data    TEXT,
    expires_at  TEXT NOT NULL,
    PRIMARY KEY (github_user, badge_id)
);`

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) ClaimBadge(user, badgeID string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(
		`INSERT INTO badge_state (github_user,badge_id,unlocked,unlocked_at,source)
		 VALUES(?,?,1,?,'claimed')
		 ON CONFLICT(github_user,badge_id) DO UPDATE SET unlocked=1,unlocked_at=?,source='claimed'`,
		user, badgeID, now, now)
	return err
}

func (s *Store) CertifyBadge(user, badgeID string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(
		`INSERT INTO badge_state (github_user,badge_id,unlocked,unlocked_at,source)
		 VALUES(?,?,1,?,'certified')
		 ON CONFLICT(github_user,badge_id) DO UPDATE SET unlocked=1,unlocked_at=?,source='certified'`,
		user, badgeID, now, now)
	return err
}

func (s *Store) IsUnlocked(user, badgeID string) (bool, error) {
	var unlocked int
	err := s.db.QueryRow(
		`SELECT unlocked FROM badge_state WHERE github_user=? AND badge_id=?`,
		user, badgeID).Scan(&unlocked)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return unlocked == 1, nil
}

func (s *Store) GetUserBadges(user string) ([]BadgeState, error) {
	rows, err := s.db.Query(
		`SELECT github_user,badge_id,unlocked,unlocked_at,source FROM badge_state
		 WHERE github_user=? AND unlocked=1`, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var states []BadgeState
	for rows.Next() {
		var bs BadgeState
		var ua sql.NullString
		if err := rows.Scan(&bs.GitHubUser, &bs.BadgeID, &bs.Unlocked, &ua, &bs.Source); err != nil {
			return nil, err
		}
		if ua.Valid {
			t, _ := time.Parse(time.RFC3339, ua.String)
			bs.UnlockedAt = &t
		}
		states = append(states, bs)
	}
	return states, rows.Err()
}

func (s *Store) SetCertCache(user, badgeID string, result bool, rawData string, ttl time.Duration) error {
	expires := time.Now().UTC().Add(ttl).Format(time.RFC3339)
	ri := 0
	if result {
		ri = 1
	}
	_, err := s.db.Exec(
		`INSERT INTO cert_cache (github_user,badge_id,result,raw_data,expires_at)
		 VALUES(?,?,?,?,?)
		 ON CONFLICT(github_user,badge_id) DO UPDATE SET result=?,raw_data=?,expires_at=?`,
		user, badgeID, ri, rawData, expires, ri, rawData, expires)
	return err
}

func (s *Store) GetCertCache(user, badgeID string) (*CacheEntry, bool, error) {
	var result int
	var rawData, expiresStr string
	err := s.db.QueryRow(
		`SELECT result,raw_data,expires_at FROM cert_cache WHERE github_user=? AND badge_id=?`,
		user, badgeID).Scan(&result, &rawData, &expiresStr)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	expires, _ := time.Parse(time.RFC3339, expiresStr)
	if time.Now().UTC().After(expires) {
		return nil, false, nil
	}
	return &CacheEntry{
		GitHubUser: user, BadgeID: badgeID,
		Result: result == 1, RawData: rawData, ExpiresAt: expires,
	}, true, nil
}
```

- [ ] **Step 4: Run tests, verify PASS**

```bash
go test ./internal/store/... -v
```

Expected: 6 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/store/ && git commit -m "feat: SQLite store with badge_state and cert_cache"
```

---

## Task 4: GitHub Client Interface & Mock

**Files:**
- Create: `internal/certify/github.go`
- Create: `internal/certify/mock.go`
- Create: `internal/certify/real_client.go`

- [ ] **Step 1: Create `internal/certify/github.go`**

Domain types and interface (no external deps):

```go
package certify

import (
	"context"
	"time"
)

type Commit struct {
	SHA       string
	Timestamp time.Time
	Repo      string
}

type PR struct {
	Number    int
	Repo      string
	Title     string
	State     string
	CreatedAt time.Time
	ClosedAt  *time.Time
	Base      string
}

type CIRun struct {
	ID         int64
	Repo       string
	Conclusion string
	CreatedAt  time.Time
	HeadSHA    string
}

type Repo struct {
	Name      string
	FullName  string
	Language  string
	Size      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Login    string
	Name     string
	Bio      string
	Blog     string
	Location string
	Company  string
	JoinedAt time.Time
}

type Event struct {
	Type      string
	Repo      string
	CreatedAt time.Time
}

type GitHubClient interface {
	GetUser(ctx context.Context, username string) (*User, error)
	ListCommits(ctx context.Context, username string) ([]Commit, error)
	ListPRs(ctx context.Context, username string) ([]PR, error)
	GetCIRuns(ctx context.Context, username, repo string) ([]CIRun, error)
	ListRepos(ctx context.Context, username string) ([]Repo, error)
	ListEvents(ctx context.Context, username string) ([]Event, error)
}
```

- [ ] **Step 2: Create `internal/certify/mock.go`**

```go
package certify

import (
	"context"
	"time"
)

type MockGitHubClient struct {
	UserFn  func(ctx context.Context, username string) (*User, error)
	Commits []Commit
	PRs     []PR
	CIRuns  map[string][]CIRun
	Repos   []Repo
	Events  []Event
	Err     error
}

func (m *MockGitHubClient) GetUser(ctx context.Context, username string) (*User, error) {
	if m.UserFn != nil {
		return m.UserFn(ctx, username)
	}
	return &User{Login: username, JoinedAt: time.Now().Add(-5 * 365 * 24 * time.Hour)}, m.Err
}

func (m *MockGitHubClient) ListCommits(ctx context.Context, username string) ([]Commit, error) {
	return m.Commits, m.Err
}

func (m *MockGitHubClient) ListPRs(ctx context.Context, username string) ([]PR, error) {
	return m.PRs, m.Err
}

func (m *MockGitHubClient) GetCIRuns(ctx context.Context, username, repo string) ([]CIRun, error) {
	if runs, ok := m.CIRuns[repo]; ok {
		return runs, m.Err
	}
	return nil, m.Err
}

func (m *MockGitHubClient) ListRepos(ctx context.Context, username string) ([]Repo, error) {
	return m.Repos, m.Err
}

func (m *MockGitHubClient) ListEvents(ctx context.Context, username string) ([]Event, error) {
	return m.Events, m.Err
}
```

- [ ] **Step 3: Create `internal/certify/real_client.go`**

Real implementation wrapping `google/go-github` + `oauth2`:

```go
package certify

import (
	"context"
	"strings"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

type realGitHubClient struct {
	client *github.Client
}

func extractRepoFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return url
}

func NewRealGitHubClient(token string) GitHubClient {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return &realGitHubClient{client: github.NewClient(tc)}
}

func (r *realGitHubClient) GetUser(ctx context.Context, username string) (*User, error) {
	u, _, err := r.client.Users.Get(ctx, username)
	if err != nil {
		return nil, err
	}
	return &User{
		Login:    u.GetLogin(),
		Name:     u.GetName(),
		Bio:      u.GetBio(),
		Blog:     u.GetBlog(),
		Location: u.GetLocation(),
		Company:  u.GetCompany(),
		JoinedAt: u.GetCreatedAt().Time,
	}, nil
}

func (r *realGitHubClient) ListCommits(ctx context.Context, username string) ([]Commit, error) {
	var all []Commit
	repos, err := r.ListRepos(ctx, username)
	if err != nil {
		return nil, err
	}
	for _, repo := range repos {
		if repo.Size == 0 {
			continue
		}
		commits, _, err := r.client.Repositories.ListCommits(ctx, username, repo.Name, &github.CommitsListOptions{ListOptions: github.ListOptions{PerPage: 100}})
		if err != nil {
			continue
		}
		for _, c := range commits {
			if c.Commit != nil && c.Commit.Author != nil {
				all = append(all, Commit{
					SHA:       c.GetSHA(),
					Timestamp: c.Commit.Author.GetDate().Time,
					Repo:      repo.Name,
				})
			}
		}
	}
	return all, nil
}

func (r *realGitHubClient) ListPRs(ctx context.Context, username string) ([]PR, error) {
	var all []PR
	opt := &github.SearchOptions{ListOptions: github.ListOptions{PerPage: 100}}
	query := "author:" + username + " type:pr"
	result, _, err := r.client.Search.Issues(ctx, query, opt)
	if err != nil {
		return nil, err
	}
	for _, issue := range result.Issues {
		if issue.PullRequestLinks == nil {
			continue
		}
		repoName := extractRepoFromURL(issue.GetRepositoryURL())
		pr := PR{
			Number:    issue.GetNumber(),
			Repo:      repoName,
			Title:     issue.GetTitle(),
			State:     issue.GetState(),
			CreatedAt: issue.GetCreatedAt().Time,
		}
		if issue.ClosedAt != nil {
			t := issue.GetClosedAt().Time
			pr.ClosedAt = &t
		}
		all = append(all, pr)
	}
	return all, nil
}

func (r *realGitHubClient) GetCIRuns(ctx context.Context, username, repo string) ([]CIRun, error) {
	runs, _, err := r.client.Actions.ListRepositoryWorkflowRuns(ctx, username, repo, &github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{PerPage: 30}})
	if err != nil {
		return nil, err
	}
	var all []CIRun
	for _, run := range runs.WorkflowRuns {
		all = append(all, CIRun{
			ID:         run.GetID(),
			Repo:       repo,
			Conclusion: run.GetConclusion(),
			CreatedAt:  run.GetCreatedAt().Time,
			HeadSHA:    run.GetHeadSHA(),
		})
	}
	return all, nil
}

func (r *realGitHubClient) ListRepos(ctx context.Context, username string) ([]Repo, error) {
	var all []Repo
	opt := &github.RepositoryListByUserOptions{ListOptions: github.ListOptions{PerPage: 100}}
	repos, _, err := r.client.Repositories.ListByUser(ctx, username, opt)
	if err != nil {
		return nil, err
	}
	for _, repo := range repos {
		if repo.GetFork() {
			continue
		}
		all = append(all, Repo{
			Name:      repo.GetName(),
			FullName:  repo.GetFullName(),
			Language:  repo.GetLanguage(),
			Size:      repo.GetSize(),
			CreatedAt: repo.GetCreatedAt().Time,
			UpdatedAt: repo.GetUpdatedAt().Time,
		})
	}
	return all, nil
}

func (r *realGitHubClient) ListEvents(ctx context.Context, username string) ([]Event, error) {
	events, _, err := r.client.Activity.ListEventsPerformedByUser(ctx, username, false, &github.ListOptions{PerPage: 100})
	if err != nil {
		return nil, err
	}
	var all []Event
	for _, e := range events {
		repoName := ""
		if e.Repo != nil {
			repoName = e.Repo.GetName()
		}
		all = append(all, Event{
			Type:      e.GetType(),
			Repo:      repoName,
			CreatedAt: e.GetCreatedAt().Time,
		})
	}
	return all, nil
}


```

- [ ] **Step 4: Verify compilation**

```bash
go build ./internal/certify/...
```

Expected: No errors.

- [ ] **Step 5: Commit**

```bash
git add internal/certify/ && git commit -m "feat: GitHubClient interface, mock, and real implementation"
```

---

## Task 5: Certify Functions & Engine

**Files:**
- Create: `internal/certify/functions.go`
- Create: `internal/certify/functions_test.go`
- Create: `internal/certify/engine.go`
- Create: `internal/certify/engine_test.go`

- [ ] **Step 1: Write certify function tests**

```go
// internal/certify/functions_test.go
package certify

import (
	"context"
	"testing"
	"time"
)

func makeCommits(count int, hour int) []Commit {
	now := time.Now()
	commits := make([]Commit, count)
	for i := 0; i < count; i++ {
		commits[i] = Commit{
			SHA:       "sha",
			Timestamp: time.Date(now.Year(), now.Month(), now.Day()-i, hour, 0, 0, 0, time.UTC),
		}
	}
	return commits
}

func TestMidnightCoder_Pass(t *testing.T) {
	// 35% night commits (35 of 100 at hour 3)
	var commits []Commit
	for i := 0; i < 35; i++ {
		commits = append(commits, Commit{Timestamp: time.Date(2025, 1, i+1, 3, 0, 0, 0, time.UTC)})
	}
	for i := 0; i < 65; i++ {
		commits = append(commits, Commit{Timestamp: time.Date(2025, 3, i+1, 14, 0, 0, 0, time.UTC)})
	}
	mock := &MockGitHubClient{Commits: commits}
	pass, err := certifyMidnightCoder(context.Background(), mock, "u")
	if err != nil {
		t.Fatal(err)
	}
	if !pass {
		t.Error("expected pass with 35% night commits")
	}
}

func TestMidnightCoder_Fail(t *testing.T) {
	mock := &MockGitHubClient{Commits: makeCommits(100, 14)}
	pass, _ := certifyMidnightCoder(context.Background(), mock, "u")
	if pass {
		t.Error("expected fail with 0% night commits")
	}
}

func TestMidnightCoder_TooFewCommits(t *testing.T) {
	mock := &MockGitHubClient{Commits: makeCommits(10, 3)}
	pass, _ := certifyMidnightCoder(context.Background(), mock, "u")
	if pass {
		t.Error("expected fail with < 50 commits")
	}
}

func TestWeekendWarrior_Pass(t *testing.T) {
	// Find dates that are Saturday/Sunday, build 30 weekend + 70 weekday
	var commits []Commit
	base := time.Date(2025, 6, 2, 10, 0, 0, 0, time.UTC) // a Monday
	for i := 0; i < 100; i++ {
		day := base.AddDate(0, 0, i)
		if i < 30 {
			// Shift to a Saturday
			day = base.AddDate(0, 0, i*7+5)
		}
		commits = append(commits, Commit{Timestamp: day})
	}
	mock := &MockGitHubClient{Commits: commits}
	pass, _ := certifyWeekendWarrior(context.Background(), mock, "u")
	if !pass {
		t.Error("expected pass")
	}
}

func TestPolyglot_Pass(t *testing.T) {
	repos := []Repo{
		{Language: "Go"}, {Language: "Python"}, {Language: "JavaScript"},
		{Language: "Rust"}, {Language: "TypeScript"}, {Language: "Go"},
	}
	mock := &MockGitHubClient{Repos: repos}
	pass, _ := certifyPolyglot(context.Background(), mock, "u")
	if !pass {
		t.Error("expected pass with 5 languages")
	}
}

func TestPolyglot_Fail(t *testing.T) {
	repos := []Repo{{Language: "Go"}, {Language: "Go"}}
	mock := &MockGitHubClient{Repos: repos}
	pass, _ := certifyPolyglot(context.Background(), mock, "u")
	if pass {
		t.Error("expected fail with 1 language")
	}
}

func TestPRBomber_Pass(t *testing.T) {
	now := time.Now()
	prs := make([]PR, 25)
	for i := range prs {
		prs[i] = PR{Number: i, CreatedAt: now.AddDate(0, 0, -i)}
	}
	mock := &MockGitHubClient{PRs: prs}
	pass, _ := certifyPRBomber(context.Background(), mock, "u")
	if !pass {
		t.Error("expected pass with 25 PRs in 30 days")
	}
}

func TestMonkeyWrench_Pass(t *testing.T) {
	mock := &MockGitHubClient{
		Repos: []Repo{{Name: "myrepo"}},
		CIRuns: map[string][]CIRun{
			"myrepo": {{Conclusion: "failure", Repo: "myrepo"}},
		},
	}
	pass, _ := certifyMonkeyWrench(context.Background(), mock, "u")
	if !pass {
		t.Error("expected pass with CI failure")
	}
}

func TestLife404_Pass(t *testing.T) {
	mock := &MockGitHubClient{
		UserFn: func(ctx context.Context, u string) (*User, error) {
			return &User{Login: u, Bio: "", Blog: "", Company: ""}, nil
		},
	}
	pass, _ := certifyLife404(context.Background(), mock, "u")
	if !pass {
		t.Error("expected pass with empty profile")
	}
}

func TestTrueDestroyer_Pass(t *testing.T) {
	mock := &MockGitHubClient{
		Repos: []Repo{{Name: "repo1"}},
		CIRuns: map[string][]CIRun{
			"repo1": {
				{Conclusion: "failure"}, {Conclusion: "failure"}, {Conclusion: "failure"},
			},
		},
	}
	pass, _ := certifyTrueDestroyer(context.Background(), mock, "u")
	if !pass {
		t.Error("expected pass with 3 consecutive failures")
	}
}
```

- [ ] **Step 2: Run tests, verify FAIL**

```bash
go test ./internal/certify/...
```

- [ ] **Step 3: Create `internal/certify/functions.go`**

```go
package certify

import (
	"context"
	"sort"
	"time"
)

type CertFunc func(ctx context.Context, gh GitHubClient, user string) (bool, error)

var CertFuncs = map[string]CertFunc{
	"midnight-coder":     certifyMidnightCoder,
	"weekend-warrior":    certifyWeekendWarrior,
	"issue-lord":         certifyIssueLord,
	"docs-master":        certifyDocsMaster,
	"pr-bomber":          certifyPRBomber,
	"monkey-wrench":      certifyMonkeyWrench,
	"archaeologist":      certifyArchaeologist,
	"branch-hoarder":     certifyBranchHoarder,
	"ghost-committer":    certifyGhostCommitter,
	"polyglot":           certifyPolyglot,
	"true-destroyer":     certifyTrueDestroyer,
	"y2k-hunter":         certifyY2KHunter,
	"life-404":           certifyLife404,
	"commit-anniversary": certifyCommitAnniversary,
	"fullstack-victim":   certifyFullstackVictim,
}

func certifyMidnightCoder(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	commits, err := gh.ListCommits(ctx, user)
	if err != nil {
		return false, err
	}
	if len(commits) < 50 {
		return false, nil
	}
	night := 0
	for _, c := range commits {
		h := c.Timestamp.Hour()
		if h >= 2 && h < 5 {
			night++
		}
	}
	return float64(night)/float64(len(commits)) >= 0.3, nil
}

func certifyWeekendWarrior(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	commits, err := gh.ListCommits(ctx, user)
	if err != nil {
		return false, err
	}
	if len(commits) < 20 {
		return false, nil
	}
	weekend := 0
	for _, c := range commits {
		d := c.Timestamp.Weekday()
		if d == time.Saturday || d == time.Sunday {
			weekend++
		}
	}
	return float64(weekend)/float64(len(commits)) >= 0.2, nil
}

func certifyIssueLord(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	events, err := gh.ListEvents(ctx, user)
	if err != nil {
		return false, err
	}
	count := 0
	for _, e := range events {
		if e.Type == "IssuesEvent" {
			count++
		}
	}
	return count >= 100, nil
}

func certifyDocsMaster(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	// Simplified: check if user has a repo with very large README relative to code
	// Full implementation would use GitHub Contents API to read README files
	repos, err := gh.ListRepos(ctx, user)
	if err != nil {
		return false, err
	}
	for _, r := range repos {
		if r.Size > 0 && r.Size < 100 {
			return true, nil
		}
	}
	return false, nil
}

func certifyPRBomber(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	prs, err := gh.ListPRs(ctx, user)
	if err != nil {
		return false, err
	}
	cutoff := time.Now().AddDate(0, 0, -30)
	count := 0
	for _, pr := range prs {
		if pr.CreatedAt.After(cutoff) {
			count++
		}
	}
	return count >= 20, nil
}

func certifyMonkeyWrench(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	repos, err := gh.ListRepos(ctx, user)
	if err != nil {
		return false, err
	}
	for _, repo := range repos {
		runs, err := gh.GetCIRuns(ctx, user, repo.Name)
		if err != nil {
			continue
		}
		for _, run := range runs {
			if run.Conclusion == "failure" {
				return true, nil
			}
		}
	}
	return false, nil
}

func certifyArchaeologist(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	repos, err := gh.ListRepos(ctx, user)
	if err != nil {
		return false, err
	}
	threeYearsAgo := time.Now().AddDate(-3, 0, 0)
	for _, r := range repos {
		if r.CreatedAt.Before(threeYearsAgo) && r.UpdatedAt.After(threeYearsAgo) {
			return true, nil
		}
	}
	return false, nil
}

func certifyBranchHoarder(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	// Simplified: branch hoarding requires per-repo branch listing
	// For MVP, return false; enhance later with ListBranches API
	return false, nil
}

func certifyGhostCommitter(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	commits, err := gh.ListCommits(ctx, user)
	if err != nil {
		return false, err
	}
	if len(commits) < 2 {
		return false, nil
	}
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Timestamp.Before(commits[j].Timestamp)
	})
	for i := 1; i < len(commits); i++ {
		gap := commits[i].Timestamp.Sub(commits[i-1].Timestamp)
		if gap > 30*24*time.Hour {
			return true, nil
		}
	}
	return false, nil
}

func certifyPolyglot(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	repos, err := gh.ListRepos(ctx, user)
	if err != nil {
		return false, err
	}
	langs := make(map[string]bool)
	for _, r := range repos {
		if r.Language != "" {
			langs[r.Language] = true
		}
	}
	return len(langs) >= 5, nil
}

func certifyTrueDestroyer(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	repos, err := gh.ListRepos(ctx, user)
	if err != nil {
		return false, err
	}
	for _, repo := range repos {
		runs, err := gh.GetCIRuns(ctx, user, repo.Name)
		if err != nil {
			continue
		}
		consecutive := 0
		for _, run := range runs {
			if run.Conclusion == "failure" {
				consecutive++
				if consecutive >= 3 {
					return true, nil
				}
			} else {
				consecutive = 0
			}
		}
	}
	return false, nil
}

func certifyY2KHunter(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	// Requires inspecting file contents/paths for 1999/2000 references
	// MVP placeholder: return false
	return false, nil
}

func certifyLife404(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	u, err := gh.GetUser(ctx, user)
	if err != nil {
		return false, err
	}
	return u.Bio == "" && u.Blog == "" && u.Company == "", nil
}

func certifyCommitAnniversary(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	commits, err := gh.ListCommits(ctx, user)
	if err != nil {
		return false, err
	}
	if len(commits) == 0 {
		return false, nil
	}
	earliest := commits[0].Timestamp
	for _, c := range commits[1:] {
		if c.Timestamp.Before(earliest) {
			earliest = c.Timestamp
		}
	}
	return time.Since(earliest).Hours()/(365.25*24) >= 5, nil
}

func certifyFullstackVictim(ctx context.Context, gh GitHubClient, user string) (bool, error) {
	repos, err := gh.ListRepos(ctx, user)
	if err != nil {
		return false, err
	}
	fe := map[string]bool{"JavaScript": true, "TypeScript": true, "HTML": true, "CSS": true, "Vue": true, "Svelte": true}
	be := map[string]bool{"Go": true, "Python": true, "Java": true, "Rust": true, "C#": true, "Ruby": true, "PHP": true}
	hasFE, hasBE, hasDevOps := false, false, false
	devOps := []string{"docker", "kubernetes", "terraform", "ansible"}
	for _, r := range repos {
		if fe[r.Language] {
			hasFE = true
		}
		if be[r.Language] {
			hasBE = true
		}
		for _, d := range devOps {
			if r.Name == d {
				hasDevOps = true
			}
		}
	}
	return hasFE && hasBE && hasDevOps, nil
}
```

- [ ] **Step 4: Run tests, verify core functions PASS**

```bash
go test ./internal/certify/... -v -run "TestMidnight|TestPolyglot|TestPRBomber|TestMonkey|TestLife404|TestTrueDestroyer"
```

- [ ] **Step 5: Create `internal/certify/engine.go`**

```go
package certify

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/womm/womm/internal/store"
)

type Engine struct {
	gh    GitHubClient
	store *store.Store
}

func NewEngine(gh GitHubClient, s *store.Store) *Engine {
	return &Engine{gh: gh, store: s}
}

type CertResult struct {
	Passed bool   `json:"passed"`
	Source string `json:"source"` // "cached" or "fresh"
}

func (e *Engine) TryCertify(ctx context.Context, githubUser, badgeID string) (*CertResult, error) {
	fn, ok := CertFuncs[badgeID]
	if !ok {
		return nil, fmt.Errorf("no cert function for badge: %s", badgeID)
	}

	if entry, ok, err := e.store.GetCertCache(githubUser, badgeID); err == nil && ok {
		return &CertResult{Passed: entry.Result, Source: "cached"}, nil
	}

	passed, err := fn(ctx, e.gh, githubUser)
	if err != nil {
		return nil, fmt.Errorf("cert function error: %w", err)
	}

	if passed {
		if err := e.store.CertifyBadge(githubUser, badgeID); err != nil {
			return nil, err
		}
	}

	raw, _ := json.Marshal(map[string]bool{"passed": passed})
	_ = e.store.SetCertCache(githubUser, badgeID, passed, string(raw), 1*time.Hour)

	return &CertResult{Passed: passed, Source: "fresh"}, nil
}

func (e *Engine) IsCertifiable(badgeID string) bool {
	_, ok := CertFuncs[badgeID]
	return ok
}
```

- [ ] **Step 6: Create `internal/certify/engine_test.go`**

```go
package certify

import (
	"context"
	"testing"
	"time"

	"github.com/womm/womm/internal/store"
)

func TestEngine_TryCertify_Pass(t *testing.T) {
	s, _ := store.Open(":memory:")
	defer s.Close()

	commits := make([]Commit, 100)
	for i := 0; i < 65; i++ {
		commits[i] = Commit{Timestamp: time.Date(2025, 1, i+1, 14, 0, 0, 0, time.UTC)}
	}
	for i := 65; i < 100; i++ {
		commits[i] = Commit{Timestamp: time.Date(2025, 3, i-64, 3, 0, 0, 0, time.UTC)}
	}

	mock := &MockGitHubClient{Commits: commits}
	eng := NewEngine(mock, s)

	result, err := eng.TryCertify(context.Background(), "testuser", "midnight-coder")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Error("expected pass")
	}
	if result.Source != "fresh" {
		t.Errorf("expected fresh, got %s", result.Source)
	}

	unlocked, _ := s.IsUnlocked("testuser", "midnight-coder")
	if !unlocked {
		t.Error("expected badge unlocked in store")
	}
}

func TestEngine_TryCertify_Cached(t *testing.T) {
	s, _ := store.Open(":memory:")
	defer s.Close()

	s.SetCertCache("testuser", "midnight-coder", true, `{"passed":true}`, 1*time.Hour)

	mock := &MockGitHubClient{}
	eng := NewEngine(mock, s)

	result, err := eng.TryCertify(context.Background(), "testuser", "midnight-coder")
	if err != nil {
		t.Fatal(err)
	}
	if result.Source != "cached" {
		t.Errorf("expected cached, got %s", result.Source)
	}
}

func TestEngine_TryCertify_UnknownBadge(t *testing.T) {
	s, _ := store.Open(":memory:")
	defer s.Close()
	eng := NewEngine(&MockGitHubClient{}, s)

	_, err := eng.TryCertify(context.Background(), "u", "nonexistent")
	if err == nil {
		t.Error("expected error for unknown badge")
	}
}
```

- [ ] **Step 7: Run all certify tests**

```bash
go test ./internal/certify/... -v
```

- [ ] **Step 8: Commit**

```bash
git add internal/certify/ && git commit -m "feat: certify engine with 15 cert functions, caching, and tests"
```

---

## Task 6: Render Engine — Themes, Icons & Templates

**Files:**
- Create: `internal/render/theme.go`
- Create: `internal/render/icons.go`
- Create: `internal/render/render.go`
- Create: `internal/render/render_test.go`
- Create: `internal/render/templates/badge.svg.tmpl`
- Create: `internal/render/templates/wide.svg.tmpl`
- Create: `internal/render/templates/terminal.svg.tmpl`
- Create: `internal/render/templates/stamp.svg.tmpl`

- [ ] **Step 1: Create `internal/render/theme.go`**

```go
package render

type Theme struct {
	Name        string
	BgColor     string
	FgColor     string
	AccentColor string
	BorderColor string
	DimColor    string
	FontFamily  string
	TitleSize   int
	SubSize     int
}

var themes = map[string]Theme{
	"pixel": {
		Name: "pixel", BgColor: "#1a1a2e", FgColor: "#00ff41",
		AccentColor: "#00ff41", BorderColor: "#00ff41", DimColor: "#336633",
		FontFamily: "monospace", TitleSize: 11, SubSize: 8,
	},
	"cyberpunk": {
		Name: "cyberpunk", BgColor: "#0d0221", FgColor: "#05d9e8",
		AccentColor: "#ff2a6d", BorderColor: "#d300c5", DimColor: "#3d1a5e",
		FontFamily: "monospace", TitleSize: 11, SubSize: 8,
	},
	"glitch": {
		Name: "glitch", BgColor: "#111111", FgColor: "#ffffff",
		AccentColor: "#ff3333", BorderColor: "#666666", DimColor: "#888888",
		FontFamily: "monospace", TitleSize: 13, SubSize: 8,
	},
	"clean": {
		Name: "clean", BgColor: "#ffffff", FgColor: "#333333",
		AccentColor: "#e056a0", BorderColor: "#dddddd", DimColor: "#999999",
		FontFamily: "sans-serif", TitleSize: 11, SubSize: 8,
	},
}

func GetTheme(name string) Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return themes["pixel"]
}
```

- [ ] **Step 2: Create `internal/render/icons.go`**

```go
package render

var icons = map[string]string{
	"checkmark":  `<path d="M20 6L9 17l-5-5" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>`,
	"eye":        `<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" stroke="currentColor" stroke-width="1.5" fill="none"/><circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"stack":      `<rect x="4" y="2" width="16" height="4" rx="1" stroke="currentColor" stroke-width="1.5" fill="none"/><rect x="4" y="10" width="16" height="4" rx="1" stroke="currentColor" stroke-width="1.5" fill="none"/><rect x="4" y="18" width="16" height="4" rx="1" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"list":       `<line x1="8" y1="6" x2="21" y2="6" stroke="currentColor" stroke-width="1.5"/><line x1="8" y1="12" x2="21" y2="12" stroke="currentColor" stroke-width="1.5"/><line x1="8" y1="18" x2="21" y2="18" stroke="currentColor" stroke-width="1.5"/>`,
	"hash":       `<line x1="4" y1="9" x2="20" y2="9" stroke="currentColor" stroke-width="1.5"/><line x1="4" y1="15" x2="20" y2="15" stroke="currentColor" stroke-width="1.5"/><line x1="10" y1="3" x2="8" y2="21" stroke="currentColor" stroke-width="1.5"/><line x1="16" y1="3" x2="14" y2="21" stroke="currentColor" stroke-width="1.5"/>`,
	"clipboard":  `<rect x="8" y="2" width="8" height="4" rx="1" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M16 4h2a2 2 0 012 2v14a2 2 0 01-2 2H6a2 2 0 01-2-2V6a2 2 0 012-2h2" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"duck":       `<circle cx="12" cy="8" r="5" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M17 10c2 0 4 1 4 3s-2 3-4 3H8" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"calendar":   `<rect x="3" y="4" width="18" height="18" rx="2" stroke="currentColor" stroke-width="1.5" fill="none"/><line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="1.5"/>`,
	"zap":        `<polygon points="13,2 3,14 12,14 11,22 21,10 12,10" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"users":      `<circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M2 21v-2a4 4 0 014-4h6a4 4 0 014 4v2" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"moon":       `<path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"sun":        `<circle cx="12" cy="12" r="5" stroke="currentColor" stroke-width="1.5" fill="none"/><line x1="12" y1="1" x2="12" y2="3" stroke="currentColor" stroke-width="1.5"/>`,
	"alert":      `<path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"book":       `<path d="M4 19.5A2.5 2.5 0 016.5 17H20" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"rocket":     `<path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 00-2.91-.09z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"wrench":     `<path d="M14.7 6.3a1 1 0 000 1.4l1.6 1.6a1 1 0 001.4 0l3.77-3.77a6 6 0 01-7.94 7.94l-6.91 6.91a2.12 2.12 0 01-3-3l6.91-6.91a6 6 0 017.94-7.94l-3.76 3.76z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"pickaxe":    `<path d="M14.5 2.5L2 15l3 3L17.5 6.5" stroke="currentColor" stroke-width="1.5" fill="none"/><path d="M14.5 2.5l7 7-3 3-7-7" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"git-branch": `<line x1="6" y1="3" x2="6" y2="15" stroke="currentColor" stroke-width="1.5"/><circle cx="18" cy="6" r="3" stroke="currentColor" stroke-width="1.5" fill="none"/><circle cx="6" cy="18" r="3" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"ghost":      `<path d="M12 2a8 8 0 00-8 8v12l3-3 3 3 3-3 3 3 3-3 3 3V10a8 8 0 00-8-8z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"globe":      `<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5" fill="none"/><line x1="2" y1="12" x2="22" y2="12" stroke="currentColor" stroke-width="1"/>`,
	"skull":      `<circle cx="12" cy="10" r="8" stroke="currentColor" stroke-width="1.5" fill="none"/><circle cx="9" cy="9" r="2" fill="currentColor"/><circle cx="15" cy="9" r="2" fill="currentColor"/>`,
	"404":        `<text x="2" y="18" font-family="monospace" font-size="14" font-weight="bold" fill="currentColor">404</text>`,
	"trophy":     `<path d="M8 21h8M12 17v4M7 4h10v4a5 5 0 01-10 0V4z" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
	"layers":     `<polygon points="12,2 2,7 12,12 22,7" stroke="currentColor" stroke-width="1.5" fill="none"/><polyline points="2,17 12,22 22,17" stroke="currentColor" stroke-width="1.5" fill="none"/>`,
}

func GetIcon(name string) string {
	if icon, ok := icons[name]; ok {
		return icon
	}
	return icons["checkmark"]
}
```

- [ ] **Step 3: Create the 4 SVG templates**

`internal/render/templates/badge.svg.tmpl`:
```xml
<svg xmlns="http://www.w3.org/2000/svg" width="{{.Width}}" height="30" viewBox="0 0 {{.Width}} 30">
  <rect width="{{.Width}}" height="30" fill="{{.Theme.BgColor}}" rx="4"/>
  <rect x="1" y="1" width="{{.RectWidth}}" height="28" fill="none" stroke="{{.Theme.BorderColor}}" stroke-width="1" rx="3" stroke-dasharray="4,2"/>
  <g transform="translate(6,5)" color="{{.Theme.AccentColor}}"><svg width="20" height="20" viewBox="0 0 24 24">{{.Icon}}</svg></g>
  <text x="30" y="13" fill="{{.Theme.FgColor}}" font-family="{{.Theme.FontFamily}}" font-size="{{.Theme.TitleSize}}" font-weight="bold" dominant-baseline="central">{{.Name}}</text>
  <text x="30" y="25" fill="{{.Theme.DimColor}}" font-family="{{.Theme.FontFamily}}" font-size="{{.Theme.SubSize}}">{{.Subtitle}}</text>
</svg>
```

`internal/render/templates/wide.svg.tmpl`:
```xml
<svg xmlns="http://www.w3.org/2000/svg" width="280" height="90" viewBox="0 0 280 90">
  <rect width="280" height="90" fill="{{.Theme.BgColor}}" rx="6"/>
  <rect x="1" y="1" width="278" height="88" fill="none" stroke="{{.Theme.BorderColor}}" stroke-width="1" rx="5"/>
  <g transform="translate(16,14)" color="{{.Theme.AccentColor}}"><svg width="40" height="40" viewBox="0 0 24 24">{{.Icon}}</svg></g>
  <text x="70" y="30" fill="{{.Theme.FgColor}}" font-family="{{.Theme.FontFamily}}" font-size="{{.Theme.TitleSize}}" font-weight="bold">{{.Name}}</text>
  <text x="70" y="50" fill="{{.Theme.DimColor}}" font-family="{{.Theme.FontFamily}}" font-size="{{.Theme.SubSize}}">{{.Subtitle}}</text>
  <text x="70" y="72" fill="{{.Theme.AccentColor}}" font-family="{{.Theme.FontFamily}}" font-size="7" text-transform="uppercase">{{.RarityLabel}}</text>
</svg>
```

`internal/render/templates/terminal.svg.tmpl`:
```xml
<svg xmlns="http://www.w3.org/2000/svg" width="320" height="60" viewBox="0 0 320 60">
  <rect width="320" height="60" fill="{{.Theme.BgColor}}" rx="4"/>
  <rect x="0" y="0" width="320" height="16" fill="{{.Theme.BorderColor}}" opacity="0.3" rx="4"/>
  <circle cx="12" cy="8" r="3" fill="#ff5f56"/><circle cx="24" cy="8" r="3" fill="#ffbd2e"/><circle cx="36" cy="8" r="3" fill="#27c93f"/>
  <text x="16" y="36" fill="{{.Theme.FgColor}}" font-family="{{.Theme.FontFamily}}" font-size="10"><tspan fill="{{.Theme.DimColor}}">$</tspan> womm: badge unlocked</text>
  <text x="16" y="52" fill="{{.Theme.AccentColor}}" font-family="{{.Theme.FontFamily}}" font-size="10" font-weight="bold">&#x2713; {{.Name}} — {{.Subtitle}}</text>
</svg>
```

`internal/render/templates/stamp.svg.tmpl`:
```xml
<svg xmlns="http://www.w3.org/2000/svg" width="160" height="160" viewBox="0 0 160 160">
  <circle cx="80" cy="80" r="72" fill="none" stroke="{{.Theme.BorderColor}}" stroke-width="3" stroke-dasharray="6,3"/>
  <circle cx="80" cy="80" r="62" fill="none" stroke="{{.Theme.FgColor}}" stroke-width="1.5"/>
  <g transform="translate(64,40)" color="{{.Theme.AccentColor}}"><svg width="32" height="32" viewBox="0 0 24 24">{{.Icon}}</svg></g>
  <text x="80" y="90" fill="{{.Theme.FgColor}}" font-family="{{.Theme.FontFamily}}" font-size="{{.Theme.TitleSize}}" font-weight="bold" text-anchor="middle">{{.Name}}</text>
  <text x="80" y="106" fill="{{.Theme.DimColor}}" font-family="{{.Theme.FontFamily}}" font-size="{{.Theme.SubSize}}" text-anchor="middle">{{.Subtitle}}</text>
  <text x="80" y="126" fill="{{.Theme.AccentColor}}" font-family="{{.Theme.FontFamily}}" font-size="7" text-anchor="middle" text-transform="uppercase">{{.RarityLabel}}</text>
</svg>
```

- [ ] **Step 4: Create `internal/render/render.go`**

```go
package render

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"unicode/utf8"

	"github.com/womm/womm/internal/badge"
)

//go:embed templates/*.svg.tmpl
var templateFiles embed.FS

type Renderer struct {
	templates map[string]*template.Template
}

type tmplData struct {
	Width       int
	RectWidth   int
	Theme       Theme
	Icon        template.HTML
	Name        string
	Subtitle    string
	RarityLabel string
}

func NewRenderer() *Renderer {
	r := &Renderer{templates: make(map[string]*template.Template)}
	for _, name := range []string{"badge", "wide", "terminal", "stamp"} {
		data, err := templateFiles.ReadFile("templates/" + name + ".svg.tmpl")
		if err != nil {
			continue
		}
		tmpl, err := template.New(name).Parse(string(data))
		if err != nil {
			continue
		}
		r.templates[name] = tmpl
	}
	return r
}

func (r *Renderer) Render(b *badge.Badge, themeName, templateName, lang string) (string, error) {
	tmpl, ok := r.templates[templateName]
	if !ok {
		return "", fmt.Errorf("unknown template: %s", templateName)
	}

	theme := GetTheme(themeName)
	name := b.LocalizedName(lang)
	subtitle := b.LocalizedSubtitle(lang)

	charLen := utf8.RuneCountInString(name)
	width := 30 + charLen*theme.TitleSize + 20
	if width < 120 {
		width = 120
	}

	data := tmplData{
		Width:       width,
		RectWidth:   width - 2,
		Theme:       theme,
		Icon:        template.HTML(GetIcon(b.Icon)),
		Name:        name,
		Subtitle:    subtitle,
		RarityLabel: string(b.Rarity),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execute: %w", err)
	}
	return buf.String(), nil
}
```

Note: The `Icon` field in `tmplData` must be `template.HTML` (not string) to allow raw SVG injection. Update the struct field type accordingly.

- [ ] **Step 5: Create `internal/render/render_test.go`**

```go
package render

import (
	"strings"
	"testing"

	"github.com/womm/womm/internal/badge"
)

func testBadge() *badge.Badge {
	return &badge.Badge{
		ID:       "test",
		Name:     map[string]string{"zh": "测试徽章", "en": "Test Badge"},
		Subtitle: map[string]string{"zh": "测试副标题", "en": "Subtitle"},
		Icon:     "checkmark",
		Type:     badge.Declarative,
		Rarity:   badge.Common,
	}
}

func TestRenderBadgeSVG(t *testing.T) {
	r := NewRenderer()
	svg, err := r.Render(testBadge(), "pixel", "badge", "zh")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing svg tag")
	}
	if !strings.Contains(svg, "测试徽章") {
		t.Error("missing badge name")
	}
}

func TestRenderEnglish(t *testing.T) {
	r := NewRenderer()
	svg, err := r.Render(testBadge(), "pixel", "badge", "en")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "Test Badge") {
		t.Error("missing English name")
	}
}

func TestRenderAllThemes(t *testing.T) {
	r := NewRenderer()
	b := testBadge()
	for _, theme := range []string{"pixel", "cyberpunk", "glitch", "clean"} {
		t.Run(theme, func(t *testing.T) {
			_, err := r.Render(b, theme, "badge", "zh")
			if err != nil {
				t.Errorf("theme %s failed: %v", theme, err)
			}
		})
	}
}

func TestRenderAllTemplates(t *testing.T) {
	r := NewRenderer()
	b := testBadge()
	for _, tmpl := range []string{"badge", "wide", "terminal", "stamp"} {
		t.Run(tmpl, func(t *testing.T) {
			svg, err := r.Render(b, "pixel", tmpl, "zh")
			if err != nil {
				t.Errorf("template %s failed: %v", tmpl, err)
			}
			if !strings.Contains(svg, "<svg") {
				t.Error("no svg output")
			}
		})
	}
}
```

- [ ] **Step 6: Run render tests**

```bash
go test ./internal/render/... -v
```

- [ ] **Step 7: Commit**

```bash
git add internal/render/ && git commit -m "feat: SVG render engine with 4 themes, icons, and templates"
```

---

## Task 7: HTTP Server & API Handlers

**Files:**
- Create: `internal/server/server.go`
- Create: `internal/server/handler.go`
- Create: `internal/server/handler_test.go`

- [ ] **Step 1: Write handler tests**

```go
// internal/server/handler_test.go
package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/womm/womm/internal/badge"
	"github.com/womm/womm/internal/render"
	"github.com/womm/womm/internal/store"
)

func setup() *Server {
	reg := badge.NewRegistry()
	badge.RegisterAll(reg)
	r := render.NewRenderer()
	s, _ := store.Open(":memory:")
	return NewServer(reg, r, nil, s)
}

func TestHealth(t *testing.T) {
	srv := setup()
	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "ok") {
		t.Error("missing ok")
	}
}

func TestBadge_Declarative(t *testing.T) {
	srv := setup()
	req := httptest.NewRequest("GET", "/api/badge/works-on-my-machine?theme=pixel", nil)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "image/svg+xml" {
		t.Error("wrong content type")
	}
}

func TestBadge_NotFound(t *testing.T) {
	srv := setup()
	req := httptest.NewRequest("GET", "/api/badge/nonexistent", nil)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	// Returns error SVG (200) rather than HTTP 404 for badge embed compatibility
	body := w.Body.String()
	if !strings.Contains(body, "svg") {
		t.Error("expected error SVG")
	}
}

func TestBadges_ListAll(t *testing.T) {
	srv := setup()
	req := httptest.NewRequest("GET", "/api/badges", nil)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	if !strings.Contains(w.Body.String(), "midnight-coder") {
		t.Error("missing badge in list")
	}
}

func TestBadges_UserBadges(t *testing.T) {
	srv := setup()
	srv.store.ClaimBadge("user1", "works-on-my-machine")
	req := httptest.NewRequest("GET", "/api/badges?user=user1", nil)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	if !strings.Contains(w.Body.String(), "works-on-my-machine") {
		t.Error("missing user badge")
	}
}
```

- [ ] **Step 2: Run tests, verify FAIL**

```bash
go test ./internal/server/...
```

- [ ] **Step 3: Create `internal/server/server.go`**

```go
package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/womm/womm/internal/badge"
	"github.com/womm/womm/internal/certify"
	"github.com/womm/womm/internal/config"
	"github.com/womm/womm/internal/render"
	"github.com/womm/womm/internal/store"
)

type Server struct {
	router   chi.Router
	registry *badge.Registry
	renderer *render.Renderer
	certEng  *certify.Engine
	store    *store.Store
}

func NewServer(reg *badge.Registry, renderer *render.Renderer, certEng *certify.Engine, s *store.Store) *Server {
	srv := &Server{
		registry: reg,
		renderer: renderer,
		certEng:  certEng,
		store:    s,
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/api/health", srv.handleHealth)
	r.Get("/api/badge/{id}", srv.handleBadge)
	r.Get("/api/badges", srv.handleBadges)
	srv.router = r
	return srv
}

func (s *Server) ListenAndServe(cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return http.ListenAndServe(addr, s.router)
}
```

- [ ] **Step 4: Create `internal/server/handler.go`**

```go
package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/womm/womm/internal/badge"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleBadge(w http.ResponseWriter, r *http.Request) {
	badgeID := chi.URLParam(r, "id")
	theme := r.URL.Query().Get("theme")
	if theme == "" {
		theme = "pixel"
	}
	style := r.URL.Query().Get("style")
	if style == "" {
		style = "badge"
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "zh"
	}
	user := r.URL.Query().Get("user")

	b, ok := s.registry.Lookup(badgeID)
	if !ok {
		s.writeErrorSVG(w, "badge not found")
		return
	}

	if b.Type == badge.Certified && s.certEng != nil && user != "" {
		result, err := s.certEng.TryCertify(r.Context(), user, badgeID)
		if err != nil || !result.Passed {
			s.writeLockedSVG(w)
			return
		}
	}

	svg, err := s.renderer.Render(b, theme, style, lang)
	if err != nil {
		http.Error(w, "render error", 500)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write([]byte(svg))
}

func (s *Server) handleBadges(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.URL.Query().Get("user")
	if user != "" {
		states, err := s.store.GetUserBadges(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(states)
		return
	}

	badges := s.registry.ListAll()
	json.NewEncoder(w).Encode(badges)
}

func (s *Server) writeErrorSVG(w http.ResponseWriter, msg string) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="280" height="30"><rect width="280" height="30" fill="#1a1a2e" rx="4"/><text x="12" y="20" fill="#ff5555" font-family="monospace" font-size="10">error: ` + msg + `</text></svg>`
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(svg))
}

func (s *Server) writeLockedSVG(w http.ResponseWriter) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="180" height="30"><rect width="180" height="30" fill="#1a1a2e" rx="4"/><text x="12" y="20" fill="#666" font-family="monospace" font-size="10">locked - not unlocked yet</text></svg>`
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(svg))
}
```

- [ ] **Step 5: Run tests**

```bash
go test ./internal/server/... -v
```

- [ ] **Step 6: Commit**

```bash
git add internal/server/ && git commit -m "feat: HTTP server with badge/badges/health endpoints"
```

---

## Task 8: CLI Commands & App Wiring

**Files:**
- Create: `internal/app/app.go`
- Modify: `cmd/root.go` (add all subcommands)
- Create: `Dockerfile`

- [ ] **Step 1: Create `internal/app/app.go`**

```go
package app

import (
	"github.com/womm/womm/internal/badge"
	"github.com/womm/womm/internal/certify"
	"github.com/womm/womm/internal/config"
	"github.com/womm/womm/internal/render"
	"github.com/womm/womm/internal/server"
	"github.com/womm/womm/internal/store"
)

type App struct {
	Config   *config.Config
	Registry *badge.Registry
	Renderer *render.Renderer
	CertEng  *certify.Engine
	Store    *store.Store
	Server   *server.Server
}

func New(cfg *config.Config) (*App, error) {
	reg := badge.NewRegistry()
	badge.RegisterAll(reg)

	renderer := render.NewRenderer()

	s, err := store.Open(cfg.Storage.Path)
	if err != nil {
		return nil, err
	}

	var certEng *certify.Engine
	if cfg.GitHub.DefaultToken != "" {
		ghClient := certify.NewRealGitHubClient(cfg.GitHub.DefaultToken)
		certEng = certify.NewEngine(ghClient, s)
	}

	srv := server.NewServer(reg, renderer, certEng, s)

	return &App{
		Config:   cfg,
		Registry: reg,
		Renderer: renderer,
		CertEng:  certEng,
		Store:    s,
		Server:   srv,
	}, nil
}
```

- [ ] **Step 2: Rewrite `cmd/root.go` with all CLI commands**

```go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/womm/womm/internal/app"
	"github.com/womm/womm/internal/config"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "womm",
	Short: "WOMM - Works On My Machine badge generator",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP badge server",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()
		fmt.Printf("WOMM server starting on %s:%d\n", a.Config.Server.Host, a.Config.Server.Port)
		return a.Server.ListenAndServe(a.Config)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available badges",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()
		badges := a.Registry.ListAll()
		for _, b := range badges {
			fmt.Printf("[%-8s] %-30s  %s\n", b.Rarity, b.LocalizedName("zh"), b.ID)
		}
		return nil
	},
}

var claimCmd = &cobra.Command{
	Use:   "claim <badge-id>",
	Short: "Claim a declarative badge",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		user, _ := cmd.Flags().GetString("user")
		if user == "" {
			user = "default"
		}
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()

		if _, ok := a.Registry.Lookup(args[0]); !ok {
			return fmt.Errorf("badge not found: %s", args[0])
		}
		return a.Store.ClaimBadge(user, args[0])
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show your unlocked badges",
	RunE: func(cmd *cobra.Command, args []string) error {
		user, _ := cmd.Flags().GetString("user")
		if user == "" {
			user = "default"
		}
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()
		states, err := a.Store.GetUserBadges(user)
		if err != nil {
			return err
		}
		if len(states) == 0 {
			fmt.Println("No badges unlocked yet.")
			return nil
		}
		for _, s := range states {
			fmt.Printf("  [%.8s] %s\n", s.Source, s.BadgeID)
		}
		return nil
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate <badge-id>",
	Short: "Generate SVG badge file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		theme, _ := cmd.Flags().GetString("theme")
		output, _ := cmd.Flags().GetString("output")
		lang, _ := cmd.Flags().GetString("lang")
		style, _ := cmd.Flags().GetString("style")

		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()

		b, ok := a.Registry.Lookup(args[0])
		if !ok {
			return fmt.Errorf("badge not found: %s", args[0])
		}
		if theme == "" {
			theme = "pixel"
		}
		if output == "" {
			output = args[0] + ".svg"
		}
		if style == "" {
			style = "badge"
		}
		if lang == "" {
			lang = "zh"
		}

		svg, err := a.Renderer.Render(b, theme, style, lang)
		if err != nil {
			return err
		}
		return os.WriteFile(output, []byte(svg), 0644)
	},
}

var tokenCmd = &cobra.Command{
	Use:   "github-token [token]",
	Short: "Set or show GitHub Personal Access Token",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()
		if len(args) == 0 {
			if a.Config.GitHub.DefaultToken != "" {
				fmt.Println("Token is configured.")
			} else {
				fmt.Println("No token configured. Set it via womm.toml or: womm github-token <TOKEN>")
			}
			return nil
		}
		a.Config.GitHub.DefaultToken = args[0]
		fmt.Println("Token updated in memory. Persist it in womm.toml for permanent use.")
		return nil
	},
}

var certifyCmd = &cobra.Command{
	Use:   "certify <badge-id>",
	Short: "Attempt to certify a badge via GitHub API",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		user, _ := cmd.Flags().GetString("user")
		if user == "" {
			return fmt.Errorf("--user flag required")
		}
		a, err := loadApp()
		if err != nil {
			return err
		}
		defer a.Store.Close()

		if a.CertEng == nil {
			return fmt.Errorf("GitHub token not configured")
		}
		result, err := a.CertEng.TryCertify(cmd.Context(), user, args[0])
		if err != nil {
			return err
		}
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

func loadApp() (*app.App, error) {
	cfg, err := config.Load(configFile)
	if err != nil {
		return nil, err
	}
	return app.New(cfg)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "womm.toml", "config file path")

	claimCmd.Flags().String("user", "", "GitHub username")
	statusCmd.Flags().String("user", "", "GitHub username")
	certifyCmd.Flags().String("user", "", "GitHub username to certify")

	generateCmd.Flags().String("theme", "", "visual theme (pixel/cyberpunk/glitch/clean)")
	generateCmd.Flags().String("output", "", "output file path")
	generateCmd.Flags().String("lang", "", "language (zh/en)")
	generateCmd.Flags().String("style", "", "template style (badge/wide/terminal/stamp)")

	rootCmd.AddCommand(serveCmd, listCmd, claimCmd, statusCmd, generateCmd, tokenCmd, certifyCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Verify build and CLI works**

```bash
go build -o womm .
./womm list
./womm claim works-on-my-machine
./womm status
./womm generate works-on-my-machine --output=test.svg
cat test.svg | head -1
rm test.svg
```

- [ ] **Step 4: Create `Dockerfile`**

```dockerfile
FROM alpine:3.19
RUN mkdir -p /data
COPY womm /usr/local/bin/womm
EXPOSE 8080
VOLUME ["/data"]
ENTRYPOINT ["womm", "serve", "-c", "/data/womm.toml"]
```

- [ ] **Step 5: Run all tests**

```bash
go test ./... -v
```

Expected: All tests pass across all packages.

- [ ] **Step 6: Commit**

```bash
git add -A && git commit -m "feat: CLI commands, app wiring, Dockerfile - MVP complete"
```

---

## Task 9: Real GitHub Client Integration & Polish

This task verifies the real GitHub client compiles correctly and adds `go.sum` cleanup.

- [ ] **Step 1: Run `go mod tidy`**

```bash
go mod tidy
```

- [ ] **Step 2: Run full test suite one more time**

```bash
go test ./... -count=1
```

- [ ] **Step 3: Run `go vet`**

```bash
go vet ./...
```

Expected: No issues.

- [ ] **Step 4: Test serve command starts**

```bash
timeout 3 ./womm serve || true
```

Expected: Server starts and logs, then exits after timeout.

- [ ] **Step 5: Final commit**

```bash
git add -A && git commit -m "chore: tidy deps, verify build clean"
```

---

## Summary

| Task | Description | Key Files |
|------|-------------|-----------|
| 1 | Scaffolding + Config | `main.go`, `cmd/root.go`, `internal/config/` |
| 2 | Badge Types + 25 Definitions | `internal/badge/` |
| 3 | SQLite Store | `internal/store/` |
| 4 | GitHub Client Interface + Mock + Real | `internal/certify/github.go`, `mock.go`, `real_client.go` |
| 5 | Certify Functions + Engine | `internal/certify/functions.go`, `engine.go` |
| 6 | Render Engine (Themes + Icons + SVG) | `internal/render/` |
| 7 | HTTP Server + API Handlers | `internal/server/` |
| 8 | CLI Commands + App Wiring + Docker | `cmd/root.go`, `internal/app/`, `Dockerfile` |
| 9 | Integration + Polish | `go mod tidy`, `go vet` |

Each task produces a working, testable, committable increment. Total: ~25 badges, 4 themes, 4 templates, 15 cert functions, HTTP API, CLI with 7 commands.
