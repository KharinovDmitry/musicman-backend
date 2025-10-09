package music

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
)

type Pack struct {
	db *pgxpool.Pool
}

func NewPack(db *pgxpool.Pool) *Pack {
	return &Pack{db: db}
}

func (r *Pack) Create(ctx context.Context, pack entity.Pack) error {
	query := `
	INSERT INTO packs (id, name, description, genre, author, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(ctx, query,
		pack.ID, pack.Name, pack.Description, string(pack.Genre), pack.Author,
		pack.CreatedAt, pack.UpdatedAt)

	return err
}

func (r *Pack) GetByID(ctx context.Context, id uuid.UUID) (entity.Pack, error) {
	query := `
	SELECT id, name, description, genre, author, created_at, updated_at
	FROM packs WHERE id = $1`

	var pack entity.Pack
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&pack.ID, &pack.Name, &pack.Description, &pack.Genre, &pack.Author,
		&pack.CreatedAt, &pack.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return pack, domain.ErrNotFound
	}
	if err != nil {
		return pack, fmt.Errorf("failed get pack from db: %w", err)
	}

	return pack, nil
}

func (r *Pack) GetAll(ctx context.Context) ([]entity.Pack, error) {
	query := `
	SELECT id, name, description, genre, author, created_at, updated_at
	FROM packs ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("failed get all packs from db: %w", err)
	}

	var packs []entity.Pack
	for rows.Next() {
		var pack entity.Pack
		err = rows.Scan(
			&pack.ID, &pack.Name, &pack.Description, &pack.Genre, &pack.Author,
			&pack.CreatedAt, &pack.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed get map packs from db: %w", err)
		}
		packs = append(packs, pack)
	}

	return packs, nil
}

func (r *Pack) Update(ctx context.Context, pack entity.Pack) error {
	query := `
	UPDATE packs SET name=$1, description=$2, genre=$3, author=$4, updated_at=$5
	WHERE id=$6`

	_, err := r.db.Exec(ctx, query,
		pack.Name, pack.Description, pack.Genre, pack.Author,
		pack.UpdatedAt, pack.ID)

	return fmt.Errorf("failed update pack from db: %w", err)
}

func (r *Pack) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM packs WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return fmt.Errorf("failed delete pack from db: %w", err)
}
