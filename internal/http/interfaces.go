package http

import (
	"context"

	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain/entity"
)

type MusicService interface {
	GetSamples(ctx context.Context) ([]entity.Sample, error)
	GetSampleDownloadURL(ctx context.Context, sample entity.Sample) (string, error)
	GetSample(ctx context.Context, sampleID uuid.UUID) (entity.Sample, error)
	CreateSample(ctx context.Context, sample entity.Sample, audioFilePath string) (entity.Sample, error)
	UpdateSample(ctx context.Context, id uuid.UUID, packID *uuid.UUID, title, author, description, genre *string) (entity.Sample, error)
	DeleteSample(ctx context.Context, id uuid.UUID) error

	GetAllPacks(ctx context.Context) ([]entity.Pack, error)
	GetPack(ctx context.Context, id uuid.UUID) (entity.Pack, error)
	CreatePack(ctx context.Context, pack entity.Pack) error
	UpdatePack(ctx context.Context, id uuid.UUID, name, description, genre *string) error
	DeletePack(ctx context.Context, id uuid.UUID) error
}
