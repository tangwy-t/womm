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