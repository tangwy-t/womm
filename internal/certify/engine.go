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
	Source string `json:"source"`
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
