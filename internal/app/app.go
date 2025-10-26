package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/musicman-backend/config"
	"github.com/musicman-backend/internal/di"
	"github.com/musicman-backend/internal/http"
	"github.com/musicman-backend/internal/scheduler"
	"time"
)

type App struct {
	container *di.Container
	http      *config.HttpConfig
	router    *gin.Engine
	scheduler *scheduler.PaymentScheduler
}

func BuildApp(ctx context.Context, cfg *config.Config) (*App, error) {
	var app App
	var err error

	app.container, err = di.CreateContainer(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create container: %w", err)
	}

	app.http = &cfg.Http

	app.router = http.SetupRouter(app.container)

	app.scheduler = scheduler.NewPaymentScheduler(time.Second, app.container.Repository.PaymentRepository, app.container.Service.Payment)

	return &app, nil
}

func (a *App) Run(ctx context.Context) error {
	errChan := make(chan error)

	go func(a *App) {
		err := a.router.Run(a.http.Addr)
		if err != nil {
			errChan <- err
		}
	}(a)

	go func(a *App) {
		a.scheduler.Start(ctx)
	}(a)

	err := <-errChan
	if err != nil {
		return fmt.Errorf("http server err: %w", err)
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) {
	a.container.Repository.Close()
}
