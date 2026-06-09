# WOMM — 反内卷与黑色幽默成就徽章生成器 设计文档

## 概述

WOMM ("Works On My Machine") 是一个自托管的 GitHub 成就徽章服务，生成带有黑色幽默风格的 SVG 徽章，供开发者嵌入 README 展示个性。服务同时提供 CLI 工具进行徽章管理与本地生成。

**核心定位：** 官方成就无聊且与实际技能相关性差。WOMM 用"人味儿"的幽默徽章填补这一空白，部分徽章通过 GitHub API 进行真实的"社会实验式"认证，增添趣味。

## 技术决策

| 决策项 | 选择 | 理由 |
|--------|------|------|
| 后端语言 | Go | 高性能、单二进制部署、适合高并发徽章请求 |
| CLI 语言 | Go | 与后端统一技术栈，编译成单二进制 |
| 架构 | 单体服务 | 自托管场景，部署极简，内部分层 |
| 数据存储 | SQLite | 零配置、单文件、嵌入式，完美契合自托管 |
| 默认视觉风格 | 像素风 / Retro Terminal | 开发者辨识度高，SVG 实现简单 |
| MVP 徽章数 | 20-30 枚 | 涵盖声明式 + 认证式，内容丰富 |
| 交付方式 | URL 嵌入 + CLI 生成静态文件 | 两者兼备 |

## 架构设计

```
┌─────────────────────────────────────────────────┐
│                  womm (binary)                   │
├──────────┬──────────┬───────────┬───────────────┤
│  HTTP    │  Badge   │  Certify  │    Store      │
│  Server  │  Render  │  Engine   │   (SQLite)    │
│          │  Engine  │           │               │
│  /api/*  │  SVG gen │  GitHub   │  badge_state  │
│  /badge/*│  themes  │  API      │  user_config  │
│  /health │  styles  │  rules    │  cache        │
├──────────┴──────────┴───────────┴───────────────┤
│                Badge Registry                    │
│     (declarative definitions + cert rules)      │
└─────────────────────────────────────────────────┘
         ▲                        ▲
         │                        │
    HTTP clients             CLI (cobra)
   (README embeds)          (womm certify/list/...)
```

### 核心模块

| 模块 | 职责 |
|------|------|
| Badge Registry | 声明所有可用徽章的元数据（名称、描述、类型、认证规则、默认图标） |
| Render Engine | 根据徽章定义 + 风格主题生成 SVG |
| Certify Engine | 根据徽章的认证规则，调用 GitHub API 或接受声明式激活 |
| Store | SQLite，存储用户徽章解锁状态、缓存 GitHub API 响应 |
| HTTP Server | 提供徽章 URL（实时 SVG）、管理 API、健康检查 |
| CLI | womm certify / list / generate 等命令 |

## 徽章目录

### A. 声明式徽章（想挂就挂，主打态度）

| # | ID | 名称 | 稀有度 | 毒鸡汤副标题 |
|---|----|------|--------|-------------|
| 1 | works-on-my-machine | 在我机器上能运行 | Common | Works on my machine — 态度即正义 |
| 2 | read-not-reply | 已读不回 | Common | Review 了你的 PR，然后…没有然后了 |
| 3 | stackoverflow-courier | Stack Overflow 搬运工 | Common | 代码从网上来，到网上去 |
| 4 | todo-collector | TODO 收藏家 | Common | // TODO: fix this later × 50 |
| 5 | comment-fundamentalist | 注释原教旨主义者 | Common | 每行代码配三行注释，包括 i++ |
| 6 | copy-paste-engineer | 复制粘贴工程师 | Common | Ctrl+C / Ctrl+V 是我的核心技能 |
| 7 | rubber-duck-master | 橡皮鸭调试大师 | Rare | 对着鸭子说话就能修 bug |
| 8 | no-friday-deploy | 周五不部署 | Rare | 血的教训换来的铁律 |
| 9 | force-push-warrior | Git Force Push 勇士 | Rare | --force 是我的日常 |
| 10 | meeting-survivor | 会议幸存者 | Common | 今天开了 6 个会，写了 0 行代码 |

**时区处理：** GitHub Events API 返回 UTC 时间戳。WOMM 根据 GitHub 用户 profile 中的 `Location` 字段尝试推断时区（关键词匹配，如 "Beijing" → Asia/Shanghai），回退为 UTC。无法精确判断时以 UTC 为准。作为娱乐徽章，"大体准确"即可。

### B. 认证式徽章（通过 GitHub API 自动验证）

| # | ID | 名称 | 认证规则 | 毒鸡汤副标题 |
|---|----|------|---------|-------------|
| 11 | midnight-coder | 午夜编码者 | Rare | 30%+ commits 在 02:00-05:00（用户本地时间） | 月亮不睡我不睡 |
| 12 | weekend-warrior | 周末战士 | Rare | 20%+ commits 在周六/周日 | 工作使我快乐（周末也是） |
| 13 | issue-lord | 百 Issue 之主 | Rare | 用户所有仓库累计 100+ open issues | 一切安好……大概 |
| 14 | docs-master | 文档仙人 | Rare | 用户最活跃仓库的 README > 1000 行且该仓库代码 < 500 行 | 代码没写几行，文档写了一本小说 |
| 15 | pr-bomber | PR 轰炸机 | Rare | 30 天内开了 20+ PR（跨所有仓库） | 一天一个 PR，医生远离我 |
| 16 | monkey-wrench | 猴子扳手 | Rare | 最近 PR 导致 CI 失败 | 我来了，CI 挂了 |
| 17 | archaeologist | 考古学家 | Legendary | 修改了 3+ 年前未动过的文件 | 挖出了上古代码 |
| 18 | branch-hoarder | 分支囤积者 | Rare | 活跃分支 > 15 个 | 每个分支都是"马上要合并的" |
| 19 | ghost-committer | 幽灵提交者 | Legendary | 连续消失 30 天后突然回归 | 我还活着，只是不想写代码 |
| 20 | polyglot | 多语言通才 | Rare | 使用 5+ 编程语言 | 什么都会一点，什么都不精 |

### C. 稀有皮肤 / 彩蛋徽章

| # | ID | 名称 | 稀有度 | 触发条件 |
|---|----|------|--------|---------|
| 21 | true-destroyer | 真·破坏王 | Legendary | 连续 3 次 CI 失败 |
| 22 | y2k-hunter | 千年虫猎人 | Legendary | 修改了含 1999/2000 的文件 |
| 23 | life-404 | 404 人生 | Legendary | 个人主页返回 404 或为空 |
| 24 | commit-anniversary | 首次提交纪念日 | Legendary | 首次 commit 距今 5/10/15 年 |
| 25 | fullstack-victim | 全栈受害者 | Legendary | 同时有前端 + 后端 + DevOps 相关仓库 |

### 徽章数据结构

```go
type Badge struct {
    ID         string            // "midnight-coder"
    Name       map[string]string // {"zh": "午夜编码者", "en": "Midnight Coder"}
    Subtitle   map[string]string // {"zh": "月亮不睡我不睡", "en": "The moon doesn't sleep, neither do I"}
    Icon       string            // 内置 SVG 图标名
    Type       BadgeType         // Declarative | Certified
    CertRule   *CertRule         // nil for declarative
    Rarity     Rarity            // Common | Rare | Legendary
}

type CertRule struct {
    APIEndpoints []string           // 需要的 GitHub API
    Threshold    map[string]float64 // 阈值配置
    Logic        string             // rule expression (documented, evaluated by CertFunc)
}
```

## 渲染引擎

### 模板类型（4 种布局）

| 模板 | 用途 |
|------|------|
| badge | 标准矩形徽章（类似 shields.io）：`[图标] 名称 · 副标题` |
| wide | 宽幅展示卡，显示认证详情：大图标 + 多行文案 |
| terminal | 模拟终端输出：`> womm: badge unlocked ✓` |
| stamp | 印章风格，声明式徽章：圆形印章 + 文字环绕 |

### 主题系统

```go
type Theme struct {
    Name        string            // "pixel", "cyberpunk", "glitch", "clean"
    Palette     map[string]string // bg, fg, accent, border, etc.
    Font        string            // "Press Start 2P", "Fira Code", etc.
    FontSize    FontConfig        // title, subtitle, label 的字号
    Background  BackgroundStyle   // solid, gradient, scanlines
    Border      BorderStyle       // none, solid, dashed, glow
    Decorations []Decoration      // scanlines, glitch-bars, pixel-dots
}
```

内置 4 套主题：pixel（默认）、cyberpunk、glitch、clean。通过 URL 参数 `?theme=cyberpunk` 或 CLI `--theme` 切换。

### SVG 生成流程

```
Badge + Theme + Template
        │
        ▼
  Go html/template
  (SVG XML 模板)
        │
        ▼
  Icon 注入（内联 SVG）
        │
        ▼
  中文文本宽度计算
  (go-text 或 fogleman/gg)
        │
        ▼
  输出 SVG
  (Content-Type: image/svg+xml)
```

关键技术点：
- 用 Go `html/template` 渲染 SVG XML（模板化、可缓存）
- 中文字符宽度计算防止溢出
- 图标全部内联 SVG path，不依赖外部资源
- 输出带 `Cache-Control` 头，CDN 友好

## API 设计

### HTTP 端点

```
GET /api/badge/{badge-id}
    ?user=github-username     (认证式徽章必须)
    &theme=pixel              (默认 pixel)
    &style=badge|wide|terminal|stamp
    &lang=zh|en               (默认 zh，见下方多语言说明)
    → 返回 SVG (Content-Type: image/svg+xml)

GET /api/badges
    ?user=github-username
    → 返回该用户所有已解锁徽章的 JSON 列表

GET /api/badges
    → 返回全部可用徽章目录 JSON

GET /api/health
    → 健康检查
```

**多语言支持：** `lang=en` 使用徽章的英文名称和副标题（如 "Midnight Coder" / "The moon doesn't sleep, neither do I"）。MVP 仅内置中英文两套文案，每枚徽章在 Registry 中同时定义中英文 Name + Subtitle 字段。

### CLI 命令

```
womm serve                          # 启动 HTTP 服务
womm list                           # 列出所有可用徽章
womm certify <badge-id> --user=xxx  # 尝试认证某个徽章
womm claim <badge-id>               # 声明获得某个徽章
womm generate <badge-id>            # 生成本地 SVG 文件
    --theme=pixel
    --output=badge.svg
womm status                         # 查看已获得的所有徽章
womm github-token                   # 配置 GitHub Personal Access Token
```

## 认证引擎

### 认证流程

```
用户请求认证 / 服务收到 badge URL
        │
        ▼
  声明式？ ── YES ──→ 直接标记已解锁
        │
        NO
        ▼
  Check Cache (SQLite)
        │
    cache hit 且未过期？ ── YES ──→ 用缓存结果
        │
        NO
        ▼
  GitHub API 调用 (go-github 库)
        │
        ▼
  CertFunc 求值
        │
    pass?
    YES    NO
     │      │
  解锁    未解锁
  写Store  返回 locked SVG
```

### CertFunc 设计

每个认证式徽章对应一个预编译的 Go 函数，不用通用脚本语言（过重）：

```go
type CertFunc func(ctx context.Context, gh GitHubClient, user string) (bool, error)

var certFuncs = map[string]CertFunc{
    "midnight-coder":   certifyMidnightCoder,
    "weekend-warrior":  certifyWeekendWarrior,
    "issue-lord":       certifyIssueLord,
    "docs-master":      certifyDocsMaster,
    "pr-bomber":        certifyPRBomber,
    "monkey-wrench":    certifyMonkeyWrench,
    "archaeologist":    certifyArchaeologist,
    "branch-hoarder":   certifyBranchHoarder,
    "ghost-committer":  certifyGhostCommitter,
    "polyglot":         certifyPolyglot,
    "true-destroyer":   certifyTrueDestroyer,
    "y2k-hunter":       certifyY2KHunter,
    "life-404":         certifyLife404,
    "commit-anniversary": certifyCommitAnniversary,
    "fullstack-victim": certifyFullstackVictim,
}
```

### 身份与鉴权模型

**MVP 采用单 token 模式：** 服务管理员配置一个 GitHub Personal Access Token（PAT），所有 badge URL 请求均使用该 token 查询公开的 GitHub 数据。

- 通过 badge URL `?user=xxx` 指定要查询的 GitHub 用户名
- 认证式徽章仅基于**公开数据**进行验证（公开 commits、公开 issues、公开 repos）
- 如需查询私有数据，用户可在部署时配置具有 `repo` 权限的 token
- 无需 OAuth 登录流程或用户体系（自托管场景）
- CLI 的 `womm github-token` 命令用于设置服务端使用的 PAT
- CLI 的 `womm claim` / `womm generate` 等命令直接操作本地 SQLite，不涉及远程用户身份

### GitHub API 集成

- 库：`google/go-github`
- 认证：服务端配置的 PAT，存在 SQLite
- 速率限制：缓存 API 响应（默认 1 小时 TTL），避免重复请求
- 默认所需权限：`read:user`（公开数据足够）

### GitHubClient 接口

```go
type GitHubClient interface {
    ListCommits(ctx context.Context, user string) ([]Commit, error)
    ListPRs(ctx context.Context, user string) ([]PR, error)
    GetCIRuns(ctx context.Context, user, repo string) ([]CIRun, error)
    GetRepos(ctx context.Context, user string) ([]Repo, error)
    GetUser(ctx context.Context, user string) (*User, error)
    GetEvents(ctx context.Context, user string) ([]Event, error)
}
```

## 数据存储

### SQLite Schema

```sql
CREATE TABLE badge_state (
    github_user  TEXT NOT NULL,     -- GitHub 用户名（如 "torvalds"）
    badge_id     TEXT NOT NULL,
    unlocked     BOOLEAN NOT NULL DEFAULT FALSE,
    unlocked_at  TIMESTAMP,
    source       TEXT,              -- "claimed" | "certified"
    PRIMARY KEY (github_user, badge_id)
);

CREATE TABLE cert_cache (
    github_user TEXT NOT NULL,
    badge_id   TEXT NOT NULL,
    result     BOOLEAN NOT NULL,
    raw_data   TEXT,                -- JSON, API 原始响应摘要
    expires_at TIMESTAMP NOT NULL,
    PRIMARY KEY (github_user, badge_id)
);
```

无独立 `users` 表。以 `github_user` 字符串为自然主键，与"无用户体系"的服务身份模型保持一致。

## 错误处理

| 场景 | 处理方式 |
|------|---------|
| 未配置 GitHub Token | 返回 SVG："请先配置 GitHub Token" |
| GitHub API 限流 | 返回缓存结果（即使过期）+ `X-Cache: STALE` 头 |
| 不存在的 badge ID | 404 SVG："徽章不存在，但你的人生 bug 是真的" |
| 认证式徽章未解锁 | 返回"锁定"状态 SVG，显示解锁条件 |
| SVG 渲染失败 | 500 内部错误，不返回破损 SVG |

## 测试策略

```
单元测试
├── render/         — 每个主题+模板组合的 SVG 输出验证
├── certify/        — 每个认证函数的 mock GitHub 数据测试
└── rule/           — 阈值计算逻辑

集成测试
└── API 路由        — httptest + SQLite 内存数据库

E2E（可选）
└── 启动服务 → 请求 badge URL → 验证 SVG 结构
```

用 mock interface 隔离 GitHub API 调用，确保测试不依赖真实 API。

## 部署

### Docker

```dockerfile
FROM alpine:3.19
RUN mkdir -p /data
COPY womm /usr/local/bin/womm
EXPOSE 8080
VOLUME ["/data"]
ENTRYPOINT ["womm", "serve"]
```

选用 Alpine（而非 scratch）以确保 CA 证书存在（GitHub API 调用需要 TLS），同时 `/data` 目录作为挂载卷持久化 SQLite 数据库。

### 配置文件 womm.toml

```toml
[server]
port = 8080
host = "0.0.0.0"

[storage]
path = "/data/womm.db"

[github]
default_token = ""
rate_limit_ttl = "1h"

[cache]
ttl = "1h"

[themes]
default = "pixel"
```

## 依赖库

| 库 | 用途 |
|----|------|
| `go-chi/chi` | HTTP 路由（轻量、符合 net/http 接口） |
| `spf13/cobra` | CLI 框架 |
| `google/go-github` | GitHub API 客户端 |
| `BurntSushi/toml` | 配置文件解析 |
| `modernc.org/sqlite` | 纯 Go SQLite 驱动（无 CGO 依赖） |
