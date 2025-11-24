package purchases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, purchase entity.Purchase) (uuid.UUID, error) {
	const query = `
		INSERT INTO purchases (id, user_uuid, sample_id, price, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query,
		purchase.ID,
		purchase.UserUUID,
		purchase.SampleID,
		purchase.Price,
		purchase.CreatedAt,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create purchase: %w", err)
	}

	return id, nil
}

func (r *Repository) GetByUserAndSample(ctx context.Context, userUUID, sampleID uuid.UUID) (entity.Purchase, error) {
	const query = `
		SELECT id, user_uuid, sample_id, price, created_at
		FROM purchases
		WHERE user_uuid = $1 AND sample_id = $2
	`

	var purchase entity.Purchase
	err := r.db.QueryRow(ctx, query, userUUID, sampleID).Scan(
		&purchase.ID,
		&purchase.UserUUID,
		&purchase.SampleID,
		&purchase.Price,
		&purchase.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Purchase{}, domain.ErrNotFound
		}
		return entity.Purchase{}, fmt.Errorf("failed to get purchase: %w", err)
	}

	return purchase, nil
}

func (r *Repository) GetByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Purchase, error) {
	const query = `
		SELECT id, user_uuid, sample_id, price, created_at
		FROM purchases
		WHERE user_uuid = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchases by user: %w", err)
	}
	defer rows.Close()

	var purchases []entity.Purchase
	for rows.Next() {
		var purchase entity.Purchase
		err := rows.Scan(
			&purchase.ID,
			&purchase.UserUUID,
			&purchase.SampleID,
			&purchase.Price,
			&purchase.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase: %w", err)
		}
		purchases = append(purchases, purchase)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating purchases: %w", err)
	}

	return purchases, nil
}
