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
