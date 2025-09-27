package service

import (
	"github.com/musicman-backend/internal/repository"
)

type Manager struct {
}

func NewManager(repository *repository.Manager) *Manager {
	return &Manager{}
}
