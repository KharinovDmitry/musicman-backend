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
}
