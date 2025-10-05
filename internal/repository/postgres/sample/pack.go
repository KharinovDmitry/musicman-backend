package sample

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/musicman-backend/internal/domain/entity"
)

type PostgresPackRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPackRepository(db *pgxpool.Pool) *PostgresPackRepository {
	return &PostgresPackRepository{db: db}
}

func (r *PostgresPackRepository) Create(ctx context.Context, pack entity.Pack) error {
	query := `
	INSERT INTO packs (id, name, description, genre, author, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(ctx, query,
		pack.ID, pack.Name, pack.Description, string(pack.Genre), pack.Author,
		pack.CreatedAt, pack.UpdatedAt)

	return err
}

func (r *PostgresPackRepository) GetByID(ctx context.Context, id uuid.UUID) (entity.Pack, error) {
	query := `
	SELECT id, name, description, genre, author, created_at, updated_at
	FROM packs WHERE id = $1`

	var pack entity.Pack
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&pack.ID, &pack.Name, &pack.Description, &pack.Genre, &pack.Author,
		&pack.CreatedAt, &pack.UpdatedAt,
	)

	return pack, fmt.Errorf("failed get pack from db: %w", err)
}

func (r *PostgresPackRepository) GetAll(ctx context.Context) ([]entity.Pack, error) {
	query := `
	SELECT id, name, description, genre, author, created_at, updated_at
	FROM packs ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
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
			return nil, fmt.Errorf("failed get all packs from db: %w", err)
		}
		packs = append(packs, pack)
	}

	return packs, nil
}

func (r *PostgresPackRepository) Update(ctx context.Context, pack entity.Pack) error {
	query := `
	UPDATE packs SET name=$1, description=$2, genre=$3, author=$4, updated_at=$5
	WHERE id=$6`

	_, err := r.db.Exec(ctx, query,
		pack.Name, pack.Description, pack.Genre, pack.Author,
		pack.UpdatedAt, pack.ID)

	return fmt.Errorf("failed update pack from db: %w", err)
}

func (r *PostgresPackRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM packs WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return fmt.Errorf("failed delete pack from db: %w", err)
}
