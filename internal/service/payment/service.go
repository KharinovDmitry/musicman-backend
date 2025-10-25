package payment

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain/constant"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/musicman-backend/pkg/client/yookassa"
	"log/slog"
	"time"
)

type BalanceController interface {
	UpdateUserBalance(ctx context.Context, userUUID uuid.UUID, amount int) error
}

type Repository interface {
	CreatePayment(ctx context.Context, payment entity.Payment) error
	UpdatePayment(ctx context.Context, payment entity.Payment) error
}

type YooKassa interface {
	CreatePayment(ctx context.Context, request yookassa.CreatePaymentRequest) (yookassa.CreatePaymentResponse, error)
	GetPayment(ctx context.Context, id string) (yookassa.PaymentByIDResponse, error)
}

type Service struct {
	yookassa YooKassa
	repo     Repository
	balance  BalanceController

	// tokenForRub - сколько токенов за рубль, да я делаю так и че ты мне сделаешь?
	tokenForRub int
}

func NewService(yookassa YooKassa, repo Repository, balance BalanceController) *Service {
	return &Service{
		repo:        repo,
		yookassa:    yookassa,
		balance:     balance,
		tokenForRub: 10,
	}
}

func (s *Service) CreatePayment(ctx context.Context, returnURI string, userUUID uuid.UUID, amount int) (string, error) {
	rub := amount / 100
	kop := amount % 100

	resp, err := s.yookassa.CreatePayment(ctx, yookassa.CreatePaymentRequest{
		Amount: yookassa.Amount{
			Value:    fmt.Sprintf("%d.%d", rub, kop),
			Currency: "RUB",
		},
		Description: "Покупка токенов",
		Test:        true,
		Confirmation: yookassa.ConfirmationCreate{
			Type:      "redirect",
			ReturnURL: returnURI,
		},
		MerchantCustomerID: userUUID.String(),
		Capture:            true,
	})
	if err != nil {
		return "", fmt.Errorf("create payment: %w", err)
	}

	err = s.repo.CreatePayment(ctx, entity.Payment{
		ID:            resp.ID,
		UserUUID:      userUUID,
		PaymentStatus: constant.PaymentStatusPending,
		Description:   "Покупка токенов",
		Amount:        amount,
		CapturedAt:    time.Time{},
		CreatedAt:     time.Now(),
	})
	if err != nil {
		return "", fmt.Errorf("save payment: %w", err)
	}

	return resp.Confirmation.ConfirmationURL, nil
}

func (s *Service) UpdatePaymentStatus(ctx context.Context, payment entity.Payment) error {
	resp, err := s.yookassa.GetPayment(ctx, payment.ID)
	if err != nil {
		return fmt.Errorf("get payment: %w", err)
	}

	if resp.Status == constant.PaymentStatusPending {
		return nil
	}

	if resp.Status == constant.PaymentStatusSucceeded {
		err = s.updateUserBalance(ctx, payment.UserUUID, payment.Amount)
		if err != nil {
			return fmt.Errorf("update user balance: %w", err)
		}
		slog.Info("payment updated",
			slog.String("payment_id", payment.ID),
			slog.String("user_uuid", payment.UserUUID.String()),
		)
	}

	payment.PaymentStatus = entity.PaymentStatus(resp.Status)

	return s.repo.UpdatePayment(ctx, payment)
}

func (s *Service) updateUserBalance(ctx context.Context, userUUID uuid.UUID, amount int) error {
	tokens := amount / 100 * s.tokenForRub
	return s.balance.UpdateUserBalance(ctx, userUUID, tokens)
}
