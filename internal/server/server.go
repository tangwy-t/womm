package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/womm/womm/internal/badge"
	"github.com/womm/womm/internal/certify"
	"github.com/womm/womm/internal/config"
	"github.com/womm/womm/internal/render"
	"github.com/womm/womm/internal/store"
)

type Server struct {
	router   chi.Router
	registry *badge.Registry
	renderer *render.Renderer
	certEng  *certify.Engine
	store    *store.Store
}

func NewServer(reg *badge.Registry, renderer *render.Renderer, certEng *certify.Engine, s *store.Store) *Server {
	srv := &Server{
		registry: reg,
		renderer: renderer,
		certEng:  certEng,
		store:    s,
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/api/health", srv.handleHealth)
	r.Get("/api/badge/{id}", srv.handleBadge)
	r.Get("/api/badges", srv.handleBadges)
	srv.router = r
	return srv
}

func (s *Server) ListenAndServe(cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	return http.ListenAndServe(addr, s.router)
}
