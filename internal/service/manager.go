package service

import (
	"github.com/musicman-backend/internal/repository"
	"github.com/musicman-backend/internal/service/auth"
	"github.com/musicman-backend/internal/service/payment"
	"github.com/musicman-backend/internal/service/token"
)

type Manager struct {
	Token   *token.Service
	Auth    *auth.Service
	Payment *payment.Service
}

func NewManager(repository *repository.Manager) *Manager {
	tokenService := token.New("secret")
	authService := auth.NewService(repository.UserRepository, tokenService)
	paymentService := payment.NewService()

	return &Manager{
		Token:   tokenService,
		Auth:    authService,
		Payment: paymentService,
	}
}
