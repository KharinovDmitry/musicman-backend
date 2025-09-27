package users

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

type User struct {
	UUID      uuid.UUID `db:"uuid"`
	Login     string    `db:"login"`
	PassHash  string    `db:"password"`
	Subscribe bool      `db:"subscribe"`
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, login, passHash string) (entity.User, error) {
	const query = `
		INSERT INTO users (login, password) 
		VALUES ($1, $2) 
		RETURNING uuid, login, password, subscribe 
	`

	var user User
	err := r.db.QueryRow(ctx, query, login, passHash).Scan(
		&user.UUID,
		&user.Login,
		&user.PassHash,
		&user.Subscribe,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return toEntity(user), nil
}

func (r *Repository) GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (entity.User, error) {
	const query = `
		SELECT uuid, login, password, subscribe 
		FROM users 
		WHERE uuid = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, userUUID).Scan(
		&user.UUID,
		&user.Login,
		&user.PassHash,
		&user.Subscribe,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, domain.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return toEntity(user), nil
}

func (r *Repository) GetUserByLogin(ctx context.Context, login string) (entity.User, error) {
	const query = `
		SELECT uuid, login, password, subscribe 
		FROM users 
		WHERE login = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, login).Scan(
		&user.UUID,
		&user.Login,
		&user.PassHash,
		&user.Subscribe,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, domain.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return toEntity(user), nil
}

func (r *Repository) SetSubscribe(ctx context.Context, userUUID uuid.UUID, isSubscribe bool) error {
	const query = `
		UPDATE users 
		SET subscribe = $1 
		WHERE uuid = $2
	`

	result, err := r.db.Exec(ctx, query, isSubscribe, userUUID)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with UUID %s not found", userUUID)
	}

	return nil
}

func toEntity(user User) entity.User {
	return entity.User{
		UUID:      user.UUID,
		Login:     user.Login,
		PassHash:  user.PassHash,
		Subscribe: user.Subscribe,
	}
}
