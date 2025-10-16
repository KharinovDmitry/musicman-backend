package entity

import (
	"github.com/google/uuid"
	"time"
)

type PaymentStatus string

type Payment struct {
	ID            string
	UserUUID      uuid.UUID
	PaymentStatus PaymentStatus
	Description   string
	Amount        int

	CapturedAt time.Time
	CreatedAt  time.Time
}
