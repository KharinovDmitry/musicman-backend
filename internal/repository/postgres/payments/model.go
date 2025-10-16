package payments

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID            uuid.UUID  `db:"id"`
	UserUUID      uuid.UUID  `db:"user_uuid"`
	PaymentStatus string     `db:"payment_status"`
	Description   string     `db:"description"`
	Amount        int        `db:"amount"`
	CapturedAt    *time.Time `db:"captured_at"`
	CreatedAt     time.Time  `db:"created_at"`
}
