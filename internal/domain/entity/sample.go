package entity

import (
	"time"

	"github.com/google/uuid"
)

type Genre string

// Sample - доменная модель сэмпла
type Sample struct {
	ID          uuid.UUID
	Title       string
	Author      string
	Description string
	Genre       Genre
	Duration    float64
	Size        int64
	MinioKey    string
	PackID      *uuid.UUID //   ну типа нуллабл :)
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Pack - доменная модель пака
type Pack struct {
	ID          uuid.UUID
	Name        string
	Description string
	Genre       Genre
	Author      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
