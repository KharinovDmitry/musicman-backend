package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain/entity"
)

type UpdateSampleRequest struct {
	Title       *string    `json:"title"`
	Author      *string    `json:"author"`
	Description *string    `json:"description"`
	Genre       *string    `json:"genre"`
	PackID      *uuid.UUID `json:"pack_id"`
	Price       *int       `json:"price"`
}

type CreatePackRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Genre       string `json:"genre" binding:"required"`
	Author      string `json:"author" binding:"required"`
}

type UUIDResponse struct {
	UUID uuid.UUID `json:"uuid"`
}

type DownloadURLResponse struct {
	DownloadURL string `json:"download_url"`
}

type UpdatePackRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Genre       *string `json:"genre"`
	Author      *string `json:"author"`
}

type SampleDTO struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Description string     `json:"description"`
	Genre       string     `json:"genre"`
	Duration    float64    `json:"duration"`
	Size        int64      `json:"size"`
	PackID      *uuid.UUID `json:"pack_id,omitempty"`
	Price       int        `json:"price"`
	ListenURL   string     `json:"listen_url"`
	DownloadURL string     `json:"download_url,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateSampleRequest struct {
	Title       string     `json:"title" binding:"required"`
	Author      string     `json:"author" binding:"required"`
	Description string     `json:"description" binding:"required"`
	Genre       string     `json:"genre" binding:"required"`
	PackID      *uuid.UUID `json:"pack_id,omitempty"`
	Price       int        `json:"price" binding:"required"`
}

type PackDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PackWithSamplesResponse struct {
	PackDTO `json:"pack"`
	Samples []SampleDTO `json:"samples"`
}

func ToSampleDTO(sample entity.Sample, listenURL string, downloadURL string) SampleDTO {
	return SampleDTO{
		ID:          sample.ID,
		Title:       sample.Title,
		Author:      sample.Author,
		Description: sample.Description,
		Genre:       sample.Genre,
		Duration:    sample.Duration,
		Size:        sample.Size,
		PackID:      sample.PackID,
		Price:       sample.Price,
		ListenURL:   listenURL,
		DownloadURL: downloadURL,
		CreatedAt:   sample.CreatedAt,
		UpdatedAt:   sample.UpdatedAt,
	}
}
func (s *SampleDTO) ToEntity() entity.Sample {
	return entity.Sample{
		ID:          s.ID,
		Title:       s.Title,
		Author:      s.Author,
		Description: s.Description,
		Genre:       s.Genre,
		Duration:    s.Duration,
		Size:        s.Size,
		PackID:      s.PackID,
		Price:       s.Price,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func ToPackDTO(pack entity.Pack) PackDTO {
	return PackDTO{
		ID:          pack.ID.String(),
		Name:        pack.Name,
		Description: pack.Description,
		Genre:       pack.Genre,
		Author:      pack.Author,
		CreatedAt:   pack.CreatedAt,
		UpdatedAt:   pack.UpdatedAt,
	}
}
