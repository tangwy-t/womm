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
			} else if run.Conclusion == "success" {
				consecutive = 0
			}
		}
	}
	return false, nil
}

func certifyY2KHunter(ctx context.Context, gh GitHubClient, user string) (bool, error) {
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
	return earliest.Before(time.Now().AddDate(-5, 0, 0)), nil
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
