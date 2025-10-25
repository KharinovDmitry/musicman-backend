package scheduler

import (
	"context"
	"github.com/musicman-backend/internal/domain/constant"
	"github.com/musicman-backend/internal/domain/entity"
	"log/slog"
	"sync"
	"time"
)

type PaymentExecutor interface {
	UpdatePaymentStatus(ctx context.Context, payment entity.Payment) error
}

type PaymentGetter interface {
	GetPaymentsByStatus(ctx context.Context, status entity.PaymentStatus) ([]entity.Payment, error)
}

type PaymentScheduler struct {
	interval time.Duration

	getter   PaymentGetter
	executor PaymentExecutor
}

func NewPaymentScheduler(interval time.Duration, getter PaymentGetter, executor PaymentExecutor) *PaymentScheduler {
	return &PaymentScheduler{
		interval: interval,
		getter:   getter,
		executor: executor,
	}
}

func (s *PaymentScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	for {
		select {
		case <-ticker.C:
			s.schedule(context.Background())
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (s *PaymentScheduler) schedule(ctx context.Context) {
	payments, err := s.getter.GetPaymentsByStatus(ctx, constant.PaymentStatusPending)
	if err != nil {
		slog.Error("failed to get payments by status", slog.String("err", err.Error()))
		return
	}

	wg := sync.WaitGroup{}
	for _, payment := range payments {
		go func() {
			wg.Add(1)
			defer wg.Done()

			err := s.executor.UpdatePaymentStatus(ctx, payment)
			if err != nil {
				slog.Error("failed to execute payment", slog.String("err", err.Error()))
				return
			}
		}()
	}

	wg.Wait()
}
