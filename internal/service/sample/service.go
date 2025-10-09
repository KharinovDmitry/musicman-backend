package sample

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/musicman-backend/internal/http/dto"
	"github.com/musicman-backend/internal/repository"
	"github.com/musicman-backend/internal/service"
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

type Sample struct {
	sampleRepo SampleRepository
	packRepo   PackRepository
	fileRepo   FileRepository
}

func NewSampleService(
	sampleRepo SampleRepository,
	packRepo PackRepository,
	fileRepo FileRepository,
) *Sample {
	return &Sample{
		sampleRepo: sampleRepo,
		packRepo:   packRepo,
		fileRepo:   fileRepo,
	}
}

func (s *Sample) CreateSample(ctx context.Context, req dto.CreateSampleRequest, byte []string, fileSize int64) (entity.Sample, error) {
	var sample entity.Sample

	if req.PackID != nil {
		if _, err := s.packRepo.GetByID(ctx, *req.PackID); err != nil {
			return sample, fmt.Errorf("pack not found while creating music: %w", err)
		}
	}

	minioKey := fmt.Sprintf("sample_%d_%s.wav", time.Now().Unix(), req.Title)

	err := s.fileRepo.UploadFile(ctx, BucketName, minioKey, filePath)
	if err != nil {
		return sample, fmt.Errorf("failed to upload sample file: %w", err)
	}

	sample = entity.Sample{
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
		Genre:       req.Genre,
		Duration:    req.Duration,
		Size:        fileSize,
		MinioKey:    minioKey,
		PackID:      req.PackID,
	}

	err = s.sampleRepo.Create(ctx, sample)
	if err != nil {
		s.fileRepo.DeleteFile(ctx, BucketName, minioKey)
		return sample, fmt.Errorf("failed to create sample: %w", err)
	}

	return sample, nil
}

func (s *Sample) GetSample(ctx context.Context, id uuid.UUID) (entity.Sample, error) {
	sample, err := s.sampleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.SampleNotFound) {
			return sample, service.SampleNotFoundError
		}

		return sample, fmt.Errorf("failed to get sample: %w", err)
	}

	return sample, nil
}

func (s *Sample) GetAllSamples(ctx context.Context) ([]entity.Sample, error) {
	samples, err := s.sampleRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get samples: %w", err)
	}
}

func (s *Sample) UpdateSample(ctx context.Context, id uuid.UUID, packID *uuid.UUID, title, author, description *string, genre *entity.Genre) (entity.Sample, error) {
	existing, err := s.sampleRepo.GetByID(ctx, id)
	if err != nil {
		return existing, fmt.Errorf("failed to get sample by id: %w", err)
	}

	if title != nil {
		existing.Title = *title
	}
	if author != nil {
		existing.Author = *author
	}
	if description != nil {
		existing.Description = *description
	}
	if genre != nil {
		existing.Genre = *genre
	}
	if packID != nil {
		if _, err := s.packRepo.GetByID(ctx, *packID); err != nil {
			return existing, fmt.Errorf("pack not found: %w", err)
		}
		existing.PackID = packID
	}

	if err = s.sampleRepo.Update(ctx, existing); err != nil {
		return existing, fmt.Errorf("failed to update sample: %w", err)
	}

	return existing, nil
}

func (s *Sample) DeleteSample(ctx context.Context, id uuid.UUID) error {
	sample, err := s.sampleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err = s.fileRepo.DeleteFile(ctx, BucketName, sample.MinioKey); err != nil {
		return fmt.Errorf("failed to delete sample file : %w", err)
	}

	if err = s.sampleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete sample: %w", err)
	}
}

func (s *Sample) GetSampleDownloadURL(ctx context.Context, id uuid.UUID) (string, error) {
	sample, err := s.sampleRepo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get sample download url: %w", err)
	}

	return s.fileRepo.GetFileURL(ctx, BucketName, sample.MinioKey)
}

func (s *Sample) CreatePack(ctx context.Context, name, description, author string, genre entity.Genre) (entity.Pack, error) {
	pack := entity.Pack{
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

func (s *Sample) GetPack(ctx context.Context, id uuid.UUID) (entity.Pack, error) {
	pack, err := s.packRepo.GetByID(ctx, id)
	if err != nil {
		return pack, fmt.Errorf("failed to get all packs: %w", err)
	}

	return pack, nil
}

func (s *Sample) GetAllPacks(ctx context.Context) ([]entity.Pack, error) {
	packs, err := s.packRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all packs: %w", err)
	}

	return packs, nil
}

func (s *Sample) UpdatePack(ctx context.Context, id uuid.UUID, name, description string, genre entity.Genre) error {
	pack, err := s.packRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed get pack by id to update it: %w", err)
	}

	if name == "" {
		pack.Name = name
	}
	if description == "" {
		pack.Description = description
	}
	if genre == "" {
		pack.Genre = genre
	}

	err = s.packRepo.Update(ctx, pack)
	if err != nil {
		return fmt.Errorf("failed udpate pack: %w", err)
	}

	return nil
}

func (s *Sample) DeletePack(ctx context.Context, id uuid.UUID) error {
	return s.packRepo.Delete(ctx, id)
}

func (s *Sample) GetPackWithSamples(ctx context.Context, id uuid.UUID) (entity.Pack, []entity.Sample, error) {
	var pack entity.Pack
	var samples []entity.Sample

	pack, err := s.packRepo.GetByID(ctx, id)
	if err != nil {
		return pack, samples, fmt.Errorf("failed to get pack: %w", err)
	}

	samples, err = s.sampleRepo.GetByPack(ctx, id)
	if err != nil {
		return pack, samples, fmt.Errorf("failed to get samples by pack: %w", err)
	}

	return pack, samples, nil
}
