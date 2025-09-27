package di

import (
	"context"
	"fmt"
	"github.com/musicman-backend/config"
	"github.com/musicman-backend/internal/repository"
	"github.com/musicman-backend/internal/service"
)

type Container struct {
	Repository *repository.Manager
	Service    *service.Manager
}

func CreateContainer(ctx context.Context, cfg *config.Config) (*Container, error) {
	var container Container
	var err error

	container.Repository, err = repository.Init(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init repository: %w", err)
	}

	container.Service = service.NewManager(container.Repository)

	return &container, nil
}
