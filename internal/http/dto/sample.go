package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	entity "github.com/musicman-backend/internal/domain/entity"
)

type CreateSampleRequest struct {
	Title       string       `json:"title" binding:"required"`
	Author      string       `json:"author" binding:"required"`
	Description string       `json:"description"`
	Genre       entity.Genre `json:"genre" binding:"required"`
	Duration    float64      `json:"duration" binding:"required"`
	PackID      *uuid.UUID   `json:"pack_id"`
}

type CreateSampleFileRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UpdateSampleRequest struct {
	ID          uuid.UUID  `json:"id" binding:"required"`
	Title       *string    `json:"title"`
	Author      *string    `json:"author"`
	Description *string    `json:"description"`
	Genre       *string    `json:"genre"`
	PackID      *uuid.UUID `json:"pack_id"`
}

type CreatePackRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Genre       string `json:"genre" binding:"required"`
	Author      string `json:"author" binding:"required"`
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
	DownloadURL string     `json:"download_url,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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

func ToSampleDTO(sample entity.Sample, downloadURL string) SampleDTO {
	return SampleDTO{
		ID:          sample.ID,
		Title:       sample.Title,
		Author:      sample.Author,
		Description: sample.Description,
		Genre:       sample.Genre,
		Duration:    sample.Duration,
		Size:        sample.Size,
		PackID:      sample.PackID,
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
