package music

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
	errors2 "github.com/musicman-backend/internal/errors"
	"github.com/musicman-backend/internal/http/dto"
	errors3 "github.com/musicman-backend/internal/service/errors"
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

type Service struct {
	sampleRepo SampleRepository
	packRepo   PackRepository
	fileRepo   FileRepository
}

func New(
	sampleRepo SampleRepository,
	packRepo PackRepository,
	fileRepo FileRepository,
) *Service {
	return &Service{
		sampleRepo: sampleRepo,
		packRepo:   packRepo,
		fileRepo:   fileRepo,
	}
}

func (s *Service) CreateSample(ctx context.Context, sample entity.Sample) (entity.Sample, error) {
	if sample.PackID != nil {
		_, err := s.packRepo.GetByID(ctx, *sample.PackID)
		if errors.Is(err, domain.ErrNotFound) {
			return sample, err
		}
		if err != nil {
			return sample, fmt.Errorf("error getting pack while creating sample: %w", err)
		}
	}

	minioKey := fmt.Sprintf("sample_%d_%s.wav", time.Now().Unix(), sample.Title)

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

func (s *Service) GetSampleBy(ctx context.Context, id uuid.UUID) (entity.Sample, error) {
	sample, err := s.sampleRepo.GetByID(ctx, id)
	if errors.Is(err, domain.ErrNotFound) {
		return sample, err
	}
	if err != nil {
		return sample, fmt.Errorf("failed to get sample: %w", err)
	}

	return sample, nil
}

func (s *Service) GetSamples(ctx context.Context) ([]entity.Sample, error) {
	samples, err := s.sampleRepo.GetAll(ctx)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get samples: %w", err)
	}

	return samples, nil
}

func (s *Service) UpdateSample(ctx context.Context, id uuid.UUID, packID *uuid.UUID, title, author, description, genre *string) (entity.Sample, error) {
	existing, err := s.sampleRepo.GetByID(ctx, id)
	if errors.Is(err, domain.ErrNotFound) {
		return existing, err
	}
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

func (s *Service) DeleteSample(ctx context.Context, id uuid.UUID) error {
	sample, err := s.sampleRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete sample by id: %w", err)
	}

	if err = s.fileRepo.DeleteFile(ctx, BucketName, sample.MinioKey); err != nil {
		return fmt.Errorf("failed to delete sample file : %w", err)
	}

	if err = s.sampleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete sample: %w", err)
	}
}

func (s *Service) GetSampleDownloadURL(ctx context.Context, minioKey string) (string, error) {
	url, err := s.fileRepo.GetFileURL(ctx, BucketName, minioKey)
	if err != nil {
		return "", fmt.Errorf("failed to get sample download url: %w", err)
	}

	return url, nil
}

func (s *Service) CreatePack(ctx context.Context, name, description, author string, genre entity.Genre) (entity.Pack, error) {
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

func (s *Service) GetPack(ctx context.Context, id uuid.UUID) (entity.Pack, error) {
	pack, err := s.packRepo.GetByID(ctx, id)
	if err != nil {
		return pack, fmt.Errorf("failed to get all packs: %w", err)
	}

	return pack, nil
}

func (s *Service) GetAllPacks(ctx context.Context) ([]entity.Pack, error) {
	packs, err := s.packRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all packs: %w", err)
	}

	return packs, nil
}

func (s *Service) UpdatePack(ctx context.Context, id uuid.UUID, name, description string, genre entity.Genre) error {
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

func (s *Service) DeletePack(ctx context.Context, id uuid.UUID) error {
	return s.packRepo.Delete(ctx, id)
}

func (s *Service) GetPackWithSamples(ctx context.Context, id uuid.UUID) (entity.Pack, []entity.Sample, error) {
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
