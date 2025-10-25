package entity

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	Login    string    `json:"login"`
	jwt.RegisteredClaims
}
