package migrator

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	dialect = "postgres"
	path    = "migrations"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Migrate(connStr string) error {
	db, err := sql.Open(dialect, connStr)
	if err != nil {
		return fmt.Errorf("sql.Open: %v", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err := goose.Up(db, path); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}
