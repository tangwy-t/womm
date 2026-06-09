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

func NewRealGitHubClient(token string) GitHubClient {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return &realGitHubClient{client: github.NewClient(tc)}
}

func extractRepoFromURL(url string) string {
	url = strings.TrimRight(url, "/")
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return url
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
		opt := &github.CommitsListOptions{ListOptions: github.ListOptions{PerPage: 100}}
		maxPages := 3
		for page := 0; page < maxPages; page++ {
			commits, resp, err := r.client.Repositories.ListCommits(ctx, username, repo.Name, opt)
			if err != nil {
				break
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
			if resp.NextPage == 0 {
				break
			}
			opt.ListOptions.Page = resp.NextPage
		}
	}
	return all, nil
}

func (r *realGitHubClient) ListPRs(ctx context.Context, username string) ([]PR, error) {
	var all []PR
	opt := &github.SearchOptions{ListOptions: github.ListOptions{PerPage: 100}}
	query := "author:" + username + " type:pr"
	maxPages := 5
	for page := 0; page < maxPages; page++ {
		result, resp, err := r.client.Search.Issues(ctx, query, opt)
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
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
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
	maxPages := 5
	for page := 0; page < maxPages; page++ {
		repos, resp, err := r.client.Repositories.ListByUser(ctx, username, opt)
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
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return all, nil
}

func (r *realGitHubClient) ListEvents(ctx context.Context, username string) ([]Event, error) {
	var all []Event
	opt := &github.ListOptions{PerPage: 100}
	maxPages := 3
	for page := 0; page < maxPages; page++ {
		events, resp, err := r.client.Activity.ListEventsPerformedByUser(ctx, username, false, opt)
		if err != nil {
			return nil, err
		}
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
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return all, nil
}