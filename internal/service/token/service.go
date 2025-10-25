package token

import (
	"context"
	"fmt"
	"github.com/musicman-backend/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/musicman-backend/internal/domain/entity"
)

const (
	Day = 24 * time.Hour

	ExpireTime = 365 * Day

	issuer = "musicman-backend"
)

type Service struct {
	secret []byte
}

func New(secret string) *Service {
	return &Service{secret: []byte(secret)}
}

func (s *Service) CreateToken(ctx context.Context, user entity.User) (string, error) {
	claims := entity.JWTClaims{
		UserUUID: user.UUID,
		Login:    user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   user.UUID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign token failed: %w", err)
	}

	return tokenString, nil
}

func (s *Service) VerifyToken(ctx context.Context, tokenString string) (entity.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.JWTClaims{}, s.keyFunc)
	if err != nil {
		return entity.JWTClaims{}, fmt.Errorf("parse token failed: %w", err)
	}

	if claims, ok := token.Claims.(*entity.JWTClaims); ok && token.Valid {
		return *claims, nil
	}

	return entity.JWTClaims{}, domain.ErrInvalidToken
}

func (s *Service) keyFunc(token *jwt.Token) (interface{}, error) {
	return s.secret, nil
}
