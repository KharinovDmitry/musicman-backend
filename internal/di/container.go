package di

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/musicman-backend/config"
	"github.com/musicman-backend/internal/repository"
	"github.com/musicman-backend/internal/service"
	"github.com/musicman-backend/pkg/client/yookassa"
	"net/url"
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

	yookassaHost, err := url.Parse(cfg.YooKassa.Host)
	if err != nil {
		return nil, fmt.Errorf("parse yookassa host: %w", err)
	}

	yookassaClient := yookassa.New(resty.New(), yookassa.Config{
		Host:      yookassaHost,
		SecretKey: cfg.YooKassa.SecretKey,
		AccountID: cfg.YooKassa.AccountID,
	})

	container.Service = service.NewManager(container.Repository, yookassaClient)

	return &container, nil
}
