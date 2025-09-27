package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/musicman-backend/config"
)

func Connect(ctx context.Context, cfg config.Postgres) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.ToDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	connConfig.MaxConns = int32(cfg.MaxConns)

	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New error: %w", err)
	}

	return pool, nil
}
