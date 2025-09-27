package auth

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
)

type UserController interface {
	GetUserByLogin(ctx context.Context, login string) (entity.User, error)
	CreateUser(ctx context.Context, login, passHash string) (entity.User, error)
}

type Tokenizer interface {
	CreateToken(ctx context.Context, user entity.User) (string, error)
}

type Service struct {
	user  UserController
	token Tokenizer
}

func NewService(user UserController, token Tokenizer) *Service {
	return &Service{
		user:  user,
		token: token,
	}
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.user.GetUserByLogin(ctx, login)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return "", fmt.Errorf("get user failed: %w", err)
	}

	if errors.Is(err, domain.ErrNotFound) {
		return "", domain.ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)) != nil {
		return "", domain.ErrInvalidCredentials
	}

	return s.token.CreateToken(ctx, user)
}

func (s *Service) Register(ctx context.Context, login, password string) (string, error) {
	_, err := s.user.GetUserByLogin(ctx, login)
	if err == nil {
		return "", domain.ErrUserAlreadyExists
	}

	if !errors.Is(err, domain.ErrNotFound) {
		return "", fmt.Errorf("get user failed: %w", err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate password failed: %w", err)
	}

	user, err := s.user.CreateUser(ctx, login, string(passHash))
	if err != nil {
		return "", fmt.Errorf("create user failed: %w", err)
	}

	return s.token.CreateToken(ctx, user)
}
