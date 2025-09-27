package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/kforge/kforge-sdk/sdkconfig"
)

func Connect(ctx context.Context, cfg sdkconfig.Postgres) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.ToDSN())
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New error: %w", err)
	}

	return pool, nil
}
