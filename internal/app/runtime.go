package app

import (
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

type Runtime struct {
	Config config.Config
	Client *api.Client
}

func Load() (*Runtime, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return New(cfg), nil
}

func New(cfg config.Config) *Runtime {
	return &Runtime{
		Config: cfg,
		Client: api.NewClient(cfg),
	}
}
