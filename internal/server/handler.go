package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/womm/womm/internal/badge"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleBadge(w http.ResponseWriter, r *http.Request) {
	badgeID := chi.URLParam(r, "id")
	theme := r.URL.Query().Get("theme")
	if theme == "" {
		theme = "pixel"
	}
	style := r.URL.Query().Get("style")
	if style == "" {
		style = "badge"
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "zh"
	}
	user := r.URL.Query().Get("user")

	b, ok := s.registry.Lookup(badgeID)
	if !ok {
		s.writeErrorSVG(w, "badge not found")
		return
	}

	if b.Type == badge.Certified && s.certEng != nil && user != "" {
		result, err := s.certEng.TryCertify(r.Context(), user, badgeID)
		if err != nil || !result.Passed {
			s.writeLockedSVG(w)
			return
		}
	}

	svg, err := s.renderer.Render(b, theme, style, lang)
	if err != nil {
		s.writeErrorSVG(w, "render error")
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write([]byte(svg))
}

func (s *Server) handleBadges(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.URL.Query().Get("user")
	if user != "" {
		states, err := s.store.GetUserBadges(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(states)
		return
	}

	badges := s.registry.ListAll()
	json.NewEncoder(w).Encode(badges)
}

func (s *Server) writeErrorSVG(w http.ResponseWriter, msg string) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="280" height="30"><rect width="280" height="30" fill="#1a1a2e" rx="4"/><text x="12" y="20" fill="#ff5555" font-family="monospace" font-size="10">error: ` + msg + `</text></svg>`
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(svg))
}

func (s *Server) writeLockedSVG(w http.ResponseWriter) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="180" height="30"><rect width="180" height="30" fill="#1a1a2e" rx="4"/><text x="12" y="20" fill="#666" font-family="monospace" font-size="10">locked - not unlocked yet</text></svg>`
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(svg))
}
