package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/musicman-backend/config"
	"github.com/musicman-backend/internal/repository/postgres"
)

type Manager struct {
	pg *pgxpool.Pool
}

func Init(ctx context.Context, cfg *config.Config) (*Manager, error) {
	var manager Manager
	var err error

	manager.pg, err = postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &manager, nil
}

func (m *Manager) Close() {
	m.pg.Close()
}
