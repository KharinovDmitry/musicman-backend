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
	UUID     uuid.UUID `db:"uuid"`
	Login    string    `db:"login"`
	PassHash string    `db:"password"`
	Tokens   int       `db:"tokens"`
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
		RETURNING uuid, login, password, tokens 
	`

	var user User
	err := r.db.QueryRow(ctx, query, login, passHash).Scan(
		&user.UUID,
		&user.Login,
		&user.PassHash,
		&user.Tokens,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return toEntity(user), nil
}

func (r *Repository) GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (entity.User, error) {
	const query = `
		SELECT uuid, login, password, tokens 
		FROM users 
		WHERE uuid = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, userUUID).Scan(
		&user.UUID,
		&user.Login,
		&user.PassHash,
		&user.Tokens,
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
		SELECT uuid, login, password, tokens 
		FROM users 
		WHERE login = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, login).Scan(
		&user.UUID,
		&user.Login,
		&user.PassHash,
		&user.Tokens,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, domain.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return toEntity(user), nil
}

func (r *Repository) UpdateUserBalance(ctx context.Context, userUUID uuid.UUID, amount int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	var curBalance int
	err = tx.QueryRow(ctx, `select tokens from users where uuid = $1`, userUUID).Scan(&curBalance)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to get user balance: %w", err)
	}

	newBalance := curBalance + amount

	_, err = tx.Exec(ctx, `update users set tokens = $1`, newBalance)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	return tx.Commit(ctx)
}

func toEntity(user User) entity.User {
	return entity.User{
		UUID:     user.UUID,
		Login:    user.Login,
		PassHash: user.PassHash,
		Tokens:   user.Tokens,
	}
}
