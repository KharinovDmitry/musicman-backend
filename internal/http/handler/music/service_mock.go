package music

import (
	"context"

	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

// MockService - мок сервиса для тестирования
type MockService struct {
	mock.Mock
}

func (m *MockService) GetSamples(ctx context.Context) ([]entity.Sample, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Sample), args.Error(1)
}

func (m *MockService) GetSampleDownloadURL(ctx context.Context, minioKey string) (string, error) {
	args := m.Called(ctx, minioKey)
	return args.String(0), args.Error(1)
}

func (m *MockService) GetSample(ctx context.Context, sampleID uuid.UUID) (entity.Sample, error) {
	args := m.Called(ctx, sampleID)
	return args.Get(0).(entity.Sample), args.Error(1)
}

func (m *MockService) CreateSample(ctx context.Context, sample entity.Sample, audioFilePath string) (entity.Sample, error) {
	args := m.Called(ctx, sample, audioFilePath)
	return args.Get(0).(entity.Sample), args.Error(1)
}

func (m *MockService) UpdateSample(ctx context.Context, id uuid.UUID, packID *uuid.UUID, title, author, description, genre *string) (entity.Sample, error) {
	args := m.Called(ctx, id, packID, title, author, description, genre)
	return args.Get(0).(entity.Sample), args.Error(1)
}

func (m *MockService) DeleteSample(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockService) GetAllPacks(ctx context.Context) ([]entity.Pack, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Pack), args.Error(1)
}

func (m *MockService) GetPack(ctx context.Context, id uuid.UUID) (entity.Pack, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Pack), args.Error(1)
}

func (m *MockService) CreatePack(ctx context.Context, pack entity.Pack) error {
	args := m.Called(ctx, pack)
	return args.Error(0)
}

func (m *MockService) UpdatePack(ctx context.Context, id uuid.UUID, name, description, genre *string) error {
	args := m.Called(ctx, id, name, description, genre)
	return args.Error(0)
}

func (m *MockService) DeletePack(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
