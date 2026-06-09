package server

import (
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
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
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

func TestBadge_Certified_NoToken(t *testing.T) {
	srv := setup()
	req := httptest.NewRequest("GET", "/api/badge/midnight-coder?user=octocat", nil)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "svg") {
		t.Error("expected SVG response")
	}
	if !strings.Contains(body, "github") && !strings.Contains(body, "Token") && !strings.Contains(body, "配置") {
		t.Errorf("expected token-needed message, got: %s", body)
	}
}
