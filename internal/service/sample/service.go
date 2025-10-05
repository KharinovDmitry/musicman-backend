package sample

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/musicman-backend/internal/http/dto"
)

const BucketName = "samples"

type SampleRepository interface {
	Create(ctx context.Context, sample entity.Sample) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.Sample, error)
	GetAll(ctx context.Context) ([]entity.Sample, error)
	GetByPack(ctx context.Context, packID uuid.UUID) ([]entity.Sample, error)
	Update(ctx context.Context, sample entity.Sample) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PackRepository interface {
	Create(ctx context.Context, pack entity.Pack) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.Pack, error)
	GetAll(ctx context.Context) ([]entity.Pack, error)
	Update(ctx context.Context, pack entity.Pack) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type FileRepository interface {
	UploadFile(ctx context.Context, bucketName, objectName, filePath string) error
	DownloadFile(ctx context.Context, bucketName, objectName, filePath string) error
	GetFileURL(ctx context.Context, bucketName, objectName string) (string, error)
	DeleteFile(ctx context.Context, bucketName, objectName string) error
}

type SampleService struct {
	sampleRepo SampleRepository
	packRepo   PackRepository
	fileRepo   FileRepository
}

func NewSampleService(
	sampleRepo SampleRepository,
	packRepo PackRepository,
	fileRepo FileRepository,
) *SampleService {
	return &SampleService{
		sampleRepo: sampleRepo,
		packRepo:   packRepo,
		fileRepo:   fileRepo,
	}
}

func (s *SampleService) CreateSample(ctx context.Context, req dto.CreateSampleRequest, filePath string, fileSize int64) (entity.Sample, error) {
	if req.PackID != nil {
		_, err := s.packRepo.GetByID(ctx, *req.PackID)
		if err != nil {
			return nil, fmt.Errorf("pack not found while creating sample: %w", err)
		}
	}

	minioKey := fmt.Sprintf("sample_%d_%s.wav", time.Now().Unix(), req.Title)

	err := s.fileRepo.UploadFile(ctx, BucketName, minioKey, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	sample := entity.Sample{
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
		Genre:       req.Genre,
		BPM:         req.BPM,
		Key:         req.Key,
		Duration:    req.Duration,
		Size:        fileSize,
		MinioKey:    minioKey,
		PackID:      req.PackID,
	}

	err = s.sampleRepo.Create(ctx, sample)
	if err != nil {
		s.storage.DeleteFile(ctx, s.bucketName, minioKey)
		return nil, fmt.Errorf("failed to create sample: %w", err)
	}

	return sample, nil
}

func (s *SampleService) GetSample(ctx context.Context, id string) (*domain.Sample, error) {
	return s.sampleRepo.GetByID(ctx, id)
}

func (s *SampleService) GetAllSamples(ctx context.Context) ([]domain.Sample, error) {
	return s.sampleRepo.GetAll(ctx)
}

func (s *SampleService) UpdateSample(ctx context.Context, id string, req dto.UpdateSampleRequest) (*domain.Sample, error) {
	existing, err := s.sampleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Author != nil {
		existing.Author = *req.Author
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Genre != nil {
		existing.Genre = *req.Genre
	}
	if req.BPM != nil {
		existing.BPM = *req.BPM
	}
	if req.Key != nil {
		existing.Key = *req.Key
	}
	if req.PackID != nil {
		if *req.PackID != "" {
			_, err := s.packRepo.GetByID(ctx, *req.PackID)
			if err != nil {
				return nil, fmt.Errorf("pack not found: %w", err)
			}
		}
		existing.PackID = req.PackID
	}

	err = s.sampleRepo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *SampleService) DeleteSample(ctx context.Context, id string) error {
	sample, err := s.sampleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.storage.DeleteFile(ctx, s.bucketName, sample.MinioKey)
	if err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	return s.sampleRepo.Delete(ctx, id)
}

func (s *SampleService) GetSampleDownloadURL(ctx context.Context, id string) (string, error) {
	sample, err := s.sampleRepo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	return s.storage.GetFileURL(ctx, s.bucketName, sample.MinioKey)
}

// Pack methods
func (s *SampleService) CreatePack(ctx context.Context, req dto.CreatePackRequest) (*domain.Pack, error) {
	pack := &domain.Pack{
		Name:        req.Name,
		Description: req.Description,
		Genre:       req.Genre,
		Author:      req.Author,
	}

	err := s.packRepo.Create(ctx, pack)
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (s *SampleService) GetPack(ctx context.Context, id string) (*domain.Pack, error) {
	return s.packRepo.GetByID(ctx, id)
}

func (s *SampleService) GetAllPacks(ctx context.Context) ([]domain.Pack, error) {
	return s.packRepo.GetAll(ctx)
}

func (s *SampleService) UpdatePack(ctx context.Context, id string, req dto.UpdatePackRequest) (*domain.Pack, error) {
	existing, err := s.packRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Genre != nil {
		existing.Genre = *req.Genre
	}
	if req.Author != nil {
		existing.Author = *req.Author
	}

	err = s.packRepo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *SampleService) DeletePack(ctx context.Context, id string) error {
	return s.packRepo.Delete(ctx, id)
}

func (s *SampleService) GetPackWithSamples(ctx context.Context, id string) (*domain.PackWithSamples, error) {
	pack, err := s.packRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	samples, err := s.sampleRepo.GetByPack(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.PackWithSamples{
		Pack:    *pack,
		Samples: samples,
	}, nil
}

func (s *SampleService) GetSampleCountByPack(ctx context.Context, packID string) (int, error) {
	return s.sampleRepo.CountByPack(ctx, packID)
}

func (s *SampleService) GetGenres() []dto.GenreResponse {
	return []dto.GenreResponse{
		{Value: string(domain.GenreHipHop), Label: "Hip Hop"},
		{Value: string(domain.GenreRock), Label: "Rock"},
		{Value: string(domain.GenreElectronic), Label: "Electronic"},
		{Value: string(domain.GenreJazz), Label: "Jazz"},
		{Value: string(domain.GenreClassical), Label: "Classical"},
		{Value: string(domain.GenrePop), Label: "Pop"},
		{Value: string(domain.GenreRB), Label: "R&B"},
	}
}
