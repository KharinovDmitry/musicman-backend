package payments

import (
	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain/entity"
	"time"
)

type Payment struct {
	ID            string    `db:"id"`
	UserUUID      uuid.UUID `db:"user_uuid"`
	PaymentStatus string    `db:"payment_status"`
	Description   string    `db:"description"`
	Amount        int       `db:"amount"`
	CapturedAt    time.Time `db:"captured_at"`
	CreatedAt     time.Time `db:"created_at"`
}

func (p Payment) ToEntity() entity.Payment {
	return entity.Payment{
		ID:            p.ID,
		UserUUID:      p.UserUUID,
		PaymentStatus: entity.PaymentStatus(p.PaymentStatus),
		Description:   p.Description,
		Amount:        p.Amount,
		CapturedAt:    p.CapturedAt,
		CreatedAt:     p.CreatedAt,
	}
}
