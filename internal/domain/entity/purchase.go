package entity

import (
	"time"

	"github.com/google/uuid"
)

// Purchase - доменная модель покупки семпла
type Purchase struct {
	ID        uuid.UUID
	UserUUID  uuid.UUID
	SampleID  uuid.UUID
	Price     int // цена в токенах на момент покупки
	CreatedAt time.Time

	Sample      *Sample
	DownloadURL string
	ListenURL   string
}
