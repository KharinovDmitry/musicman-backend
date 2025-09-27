package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/musicman-backend/config"
	"github.com/musicman-backend/internal/app"
)

const (
	defaultPath = "config/config.yaml"
)

func main() {
	path := defaultPath
	if len(os.Args) > 2 {
		path = os.Args[1]
	}

	cfg, err := config.ParseConfig(path)
	if err != nil {
		log.Fatalf("parsing config: %s", err.Error())
	}

	ctx := context.Background()

	application, err := app.BuildApp(ctx, cfg)
	if err != nil {
		log.Fatalf("creating application: %s", err.Error())
	}

	go func(ctx context.Context) {
		err = application.Run(ctx)
		if err != nil {
			log.Fatalf("running application: %s", err.Error())
		}
	}(ctx)

	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, os.Interrupt, syscall.SIGTERM)

	<-closeCh

	application.Shutdown(ctx)
}
