package service

import (
	"github.com/musicman-backend/internal/repository"
	"github.com/musicman-backend/internal/service/auth"
	"github.com/musicman-backend/internal/service/token"
)

type Manager struct {
	Token *token.Service
	Auth  *auth.Service
}

func NewManager(repository *repository.Manager) *Manager {
	tokenService := token.New("да я храню секрет тут, и че ты мне сделаешь")
	authService := auth.NewService(repository.UserRepository, tokenService)

	return &Manager{
		Token: tokenService,
		Auth:  authService,
	}
}
