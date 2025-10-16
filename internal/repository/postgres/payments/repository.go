package payments

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/musicman-backend/internal/domain/entity"
)

type Repository struct {
	db *pgxpool.Pool
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) GetPayments(ctx context.Context) ([]entity.Payment, error) {
	return nil, nil
}

func (r *Repository) GetPaymentsByStatus(ctx context.Context, status entity.PaymentStatus) ([]entity.Payment, error) {
	return nil, nil
}

func (r *Repository) CreatePayment(ctx context.Context, payment entity.Payment) error {
	return nil
}

func (r *Repository) UpdatePayment(ctx context.Context, payment entity.Payment) error {
	return nil
}
