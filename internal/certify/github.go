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