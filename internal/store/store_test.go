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
