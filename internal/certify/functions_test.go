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
	var commits []Commit
	base := time.Date(2025, 6, 2, 10, 0, 0, 0, time.UTC)
	for i := 0; i < 100; i++ {
		day := base.AddDate(0, 0, i)
		if i < 30 {
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
