package payments

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/musicman-backend/internal/domain/entity"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetPaymentsByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Payment, error) {
	const query = `
		SELECT id, user_uuid, payment_status, description, amount, captured_at, created_at 
		FROM payments 
		WHERE user_uuid = $1 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by user: %w", err)
	}
	defer rows.Close()

	var payments []entity.Payment
	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.ID,
			&payment.UserUUID,
			&payment.PaymentStatus,
			&payment.Description,
			&payment.Amount,
			&payment.CapturedAt,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment.ToEntity())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payments: %w", err)
	}

	return payments, nil
}

func (r *Repository) GetPaymentsByStatus(ctx context.Context, status entity.PaymentStatus) ([]entity.Payment, error) {
	const query = `
		SELECT id, user_uuid, payment_status, description, amount, captured_at, created_at 
		FROM payments 
		WHERE payment_status = $1 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by status: %w", err)
	}
	defer rows.Close()

	var payments []entity.Payment
	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.ID,
			&payment.UserUUID,
			&payment.PaymentStatus,
			&payment.Description,
			&payment.Amount,
			&payment.CapturedAt,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment.ToEntity())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payments: %w", err)
	}

	return payments, nil
}

func (r *Repository) CreatePayment(ctx context.Context, payment entity.Payment) error {
	const query = `
		INSERT INTO payments (id, user_uuid, payment_status, description, amount, captured_at, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		payment.ID,
		payment.UserUUID,
		string(payment.PaymentStatus),
		payment.Description,
		payment.Amount,
		payment.CapturedAt,
		payment.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

func (r *Repository) UpdatePayment(ctx context.Context, payment entity.Payment) error {
	const query = `
		UPDATE payments 
		SET payment_status = $1, description = $2, amount = $3, captured_at = $4 
		WHERE id = $5
	`

	result, err := r.db.Exec(ctx, query,
		string(payment.PaymentStatus),
		payment.Description,
		payment.Amount,
		payment.CapturedAt,
		payment.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("payment with id %s not found", payment.ID)
	}

	return nil
}
