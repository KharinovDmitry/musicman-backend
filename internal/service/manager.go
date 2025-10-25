package service

import (
	"github.com/musicman-backend/internal/repository"
	"github.com/musicman-backend/internal/service/auth"
	"github.com/musicman-backend/internal/service/payment"
	"github.com/musicman-backend/internal/service/token"
	"github.com/musicman-backend/pkg/client/yookassa"
)

type Manager struct {
	Token   *token.Service
	Auth    *auth.Service
	Payment *payment.Service
}

func NewManager(repository *repository.Manager, yookassa *yookassa.Client) *Manager {
	tokenService := token.New("secret")
	authService := auth.NewService(repository.UserRepository, tokenService)
	paymentService := payment.NewService(yookassa, repository.PaymentRepository, repository.UserRepository)

	return &Manager{
		Token:   tokenService,
		Auth:    authService,
		Payment: paymentService,
	}
}
