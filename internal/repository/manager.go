package repository

import (
	"context"
	"fmt"
	"github.com/musicman-backend/cmd/migrator"
	"github.com/musicman-backend/internal/repository/postgres/payments"
	"github.com/musicman-backend/internal/repository/postgres/users"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/musicman-backend/config"
	"github.com/musicman-backend/internal/repository/postgres"
)

type Manager struct {
	UserRepository    *users.Repository
	PaymentRepository *payments.Repository

	pg *pgxpool.Pool
}

func Init(ctx context.Context, cfg *config.Config) (*Manager, error) {
	var manager Manager
	var err error

	manager.pg, err = postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	err = migrator.Migrate(cfg.Postgres.ToDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	manager.UserRepository = users.NewRepository(manager.pg)
	manager.PaymentRepository = payments.New(manager.pg)

	return &manager, nil
}

func (m *Manager) Close() {
	m.pg.Close()
}
