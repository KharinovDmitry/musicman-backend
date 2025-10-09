package music

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
	"golang.org/x/net/context"
)

type Sample struct {
	db *pgxpool.Pool
}

func NewSample(db *pgxpool.Pool) *Sample {
	return &Sample{db: db}
}

func (r *Sample) Create(ctx context.Context, sample entity.Sample) error {
	query := `
	INSERT INTO samples (title, author, description, genre, duration, size, minio_key, pack_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $8, $9, $10, $11, $12, $13)`

	_, err := r.db.Exec(ctx, query,
		sample.Title, sample.Author, sample.Description, string(sample.Genre), sample.Duration, sample.Size, sample.MinioKey,
		sample.PackID, sample.CreatedAt, sample.UpdatedAt)

	return fmt.Errorf("failed to create sample in DB: %w", err)
}

func (r *Sample) GetByID(ctx context.Context, id uuid.UUID) (entity.Sample, error) {
	query := `
	SELECT id, title, author, description, genre, duration, size, minio_key, pack_id, created_at, updated_at
	FROM samples WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)
	sample, err := r.scanSample(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return sample, domain.ErrNotFound
	}
	if err != nil {
		return sample, fmt.Errorf("failed get pack from db: %w", err)
	}

	return sample, nil
}

func (r *Sample) GetAll(ctx context.Context) ([]entity.Sample, error) {
	query := `
	SELECT id, title, author, description, genre, duration, size, minio_key, pack_id, created_at, updated_at
	FROM samples ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error get all samples from DB: %w", err)
	}

	var samples []entity.Sample
	for rows.Next() {
		var sample entity.Sample
		var genre string
		var packID sql.Null[uuid.UUID]

		err = rows.Scan(
			&sample.ID, &sample.Title, &sample.Author, &sample.Description, &genre,
			&sample.Duration, &sample.Size, &sample.MinioKey,
			&packID, &sample.CreatedAt, &sample.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error get all samples from DB: %w", err)
		}

		sample.Genre = entity.Genre(genre)
		if packID.Valid {
			sample.PackID = &packID.V
		}

		samples = append(samples, sample)
	}

	return samples, nil
}

func (r *Sample) GetByPack(ctx context.Context, packID uuid.UUID) ([]entity.Sample, error) {
	query := `
	SELECT id, title, author, description, genre, key, duration, size, minio_key, pack_id, created_at, updated_at
	FROM samples WHERE pack_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, packID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	var samples []entity.Sample
	for rows.Next() {
		var sample entity.Sample
		var genre string
		var packID sql.Null[uuid.UUID]

		err = rows.Scan(
			&sample.ID, &sample.Title, &sample.Author, &sample.Description, &genre,
			&sample.Duration, &sample.Size, &sample.MinioKey,
			&packID, &sample.CreatedAt, &sample.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error get all samples from DB: %w", err)
		}

		sample.Genre = entity.Genre(genre)
		if packID.Valid {
			sample.PackID = &packID.V
		}

		samples = append(samples, sample)
	}

	return samples, nil
}

func (r *Sample) Update(ctx context.Context, sample entity.Sample) error {
	query := `
	UPDATE samples SET title=$1, author=$2, description=$3, genre=$4, 
	                   duration=$7, size=$8, minio_key=$9, pack_id=$10, updated_at=$11
	WHERE id=$12`

	_, err := r.db.Exec(ctx, query,
		sample.Title, sample.Author, sample.Description, sample.Genre,
		sample.Duration, sample.Size, sample.MinioKey, sample.PackID,
		sample.UpdatedAt, sample.ID)

	return err
}

func (r *Sample) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM samples WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)

	return fmt.Errorf("failed delete sample from DB: %w", err)
}

func (r *Sample) scanSample(row pgx.Row) (entity.Sample, error) {
	var sample entity.Sample
	var genre string
	var packID sql.Null[uuid.UUID]

	err := row.Scan(
		&sample.ID, &sample.Title, &sample.Author, &sample.Description, &genre,
		&sample.Duration, &sample.Size, &sample.MinioKey,
		&packID, &sample.CreatedAt, &sample.UpdatedAt,
	)

	sample.Genre = entity.Genre(genre)
	if packID.Valid {
		sample.PackID = &packID.V
	}

	return sample, err
}
