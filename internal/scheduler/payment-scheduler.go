package scheduler

import (
	"context"
	"time"
)

type PaymentController interface {
	GetPaymentsByStatus(ctx context.Context)
}

type PaymentScheduler struct {
	interval time.Duration
}
