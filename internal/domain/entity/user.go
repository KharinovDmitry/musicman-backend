package entity

import "github.com/google/uuid"

type User struct {
	UUID      uuid.UUID
	Login     string
	PassHash  string
	Subscribe bool
}
