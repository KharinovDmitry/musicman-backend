package entity

import "time"

type Genre string

// Sample - доменная модель сэмпла
type Sample struct {
	ID          string
	Title       string
	Author      string
	Description string
	Genre       Genre
	Duration    float64
	Size        int64
	MinioKey    string
	PackID      *string //   ну типа нуллабл :)
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Pack - доменная модель пака
type Pack struct {
	ID          string
	Name        string
	Description string
	Genre       Genre
	Author      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
