package token_extracter

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrMissingAuthHeader = errors.New("missing auth header")
	ErrInvalidAuthHeader = errors.New("invalid auth header")
)

const (
	headerAuthorization = "Authorization"
)

func ExtractToken(req *http.Request) (string, error) {
	authHeader := req.Header.Get(headerAuthorization)
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}
