package app

import (
	"github.com/womm/womm/internal/badge"
	"github.com/womm/womm/internal/certify"
	"github.com/womm/womm/internal/config"
	"github.com/womm/womm/internal/render"
	"github.com/womm/womm/internal/server"
	"github.com/womm/womm/internal/store"
)

type App struct {
	Config   *config.Config
	Registry *badge.Registry
	Renderer *render.Renderer
	CertEng  *certify.Engine
	Store    *store.Store
	Server   *server.Server
}

func New(cfg *config.Config) (*App, error) {
	reg := badge.NewRegistry()
	badge.RegisterAll(reg)

	renderer := render.NewRenderer()

	s, err := store.Open(cfg.Storage.Path)
	if err != nil {
		return nil, err
	}

	var certEng *certify.Engine
	if cfg.GitHub.DefaultToken != "" {
		ghClient := certify.NewRealGitHubClient(cfg.GitHub.DefaultToken)
		certEng = certify.NewEngine(ghClient, s)
	}

	srv := server.NewServer(reg, renderer, certEng, s)

	return &App{
		Config:   cfg,
		Registry: reg,
		Renderer: renderer,
		CertEng:  certEng,
		Store:    s,
		Server:   srv,
	}, nil
}
