package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain/entity"
)

type CreateSampleRequest struct {
	Title       string       `json:"title" binding:"required"`
	Author      string       `json:"author" binding:"required"`
	Description string       `json:"description"`
	Genre       entity.Genre `json:"genre" binding:"required"`
	Duration    float64      `json:"duration" binding:"required"`
	PackID      *uuid.UUID   `json:"pack_id"`
}

type UpdateSampleRequest struct {
	Title       *string       `json:"title"`
	Author      *string       `json:"author"`
	Description *string       `json:"description"`
	Genre       *entity.Genre `json:"genre"`
	BPM         *int          `json:"bpm"`
	Key         *string       `json:"key"`
	PackID      *string       `json:"pack_id"`
}

type CreatePackRequest struct {
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description"`
	Genre       entity.Genre `json:"genre" binding:"required"`
	Author      string       `json:"author" binding:"required"`
}

type UpdatePackRequest struct {
	Name        *string       `json:"name"`
	Description *string       `json:"description"`
	Genre       *entity.Genre `json:"genre"`
	Author      *string       `json:"author"`
}

type SampleResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
	BPM         int       `json:"bpm"`
	Key         string    `json:"key"`
	Duration    float64   `json:"duration"`
	Size        int64     `json:"size"`
	PackID      *string   `json:"pack_id,omitempty"`
	DownloadURL string    `json:"download_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PackResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
	Author      string    `json:"author"`
	SampleCount int       `json:"sample_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PackWithSamplesResponse struct {
	PackResponse
	Samples []SampleResponse `json:"samples"`
}

type GenreResponse struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func ToSampleResponse(sample *entity.Sample, downloadURL string) SampleResponse {
	return SampleResponse{
		ID:          sample.ID,
		Title:       sample.Title,
		Author:      sample.Author,
		Description: sample.Description,
		Genre:       string(sample.Genre),
		Key:         sample.Key,
		Duration:    sample.Duration,
		Size:        sample.Size,
		PackID:      sample.PackID,
		DownloadURL: downloadURL,
		CreatedAt:   sample.CreatedAt,
		UpdatedAt:   sample.UpdatedAt,
	}
}

func ToPackResponse(pack *entity.Pack, sampleCount int) PackResponse {
	return PackResponse{
		ID:          pack.ID,
		Name:        pack.Name,
		Description: pack.Description,
		Genre:       string(pack.Genre),
		Author:      pack.Author,
		SampleCount: sampleCount,
		CreatedAt:   pack.CreatedAt,
		UpdatedAt:   pack.UpdatedAt,
	}
}
