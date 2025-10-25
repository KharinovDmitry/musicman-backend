package dto

import "github.com/google/uuid"

type UserProfile struct {
	UUID   uuid.UUID `json:"uuid"`
	Login  string    `json:"login"`
	Tokens int       `json:"tokens"`
}
