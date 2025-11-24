package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrAlreadyPurchased   = errors.New("sample already purchased")
	ErrInsufficientTokens = errors.New("insufficient tokens")
)
