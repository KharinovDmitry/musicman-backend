package music

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
)

const BucketName = "samples"

type SampleRepository interface {
	Create(ctx context.Context, sample entity.Sample) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Sample, error)
	GetAll(ctx context.Context) ([]entity.Sample, error)
	GetByPack(ctx context.Context, packID uuid.UUID) ([]entity.Sample, error)
	Update(ctx context.Context, sample entity.Sample) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PackRepository interface {
	Create(ctx context.Context, pack entity.Pack) (uuid.UUID, error)
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
	CreateBucketIfNotExists(ctx context.Context, bucketName string) error
}

type UserRepository interface {
	GetUserByLogin(ctx context.Context, login string) (entity.User, error)
}

type Service struct {
	sampleRepo SampleRepository
	packRepo   PackRepository
	fileRepo   FileRepository
	userRepo   UserRepository
}

func New(
	sampleRepo SampleRepository,
	packRepo PackRepository,
	fileRepo FileRepository,
	userRepo UserRepository,
) *Service {
	return &Service{
		sampleRepo: sampleRepo,
		packRepo:   packRepo,
		fileRepo:   fileRepo,
		userRepo:   userRepo,
	}
}

func (s *Service) CreateSample(ctx context.Context, author, title, description, genre string, packID *uuid.UUID, price int) (uuid.UUID, error) {
	var sampleID uuid.UUID
	if packID != nil {
		_, err := s.packRepo.GetByID(ctx, *packID)
		if errors.Is(err, domain.ErrNotFound) {
			return sampleID, err
		}
		if err != nil {
			return sampleID, fmt.Errorf("error getting pack while creating sample: %w", err)
		}
	}
	if _, err := s.userRepo.GetUserByLogin(ctx, author); err != nil {
		return sampleID, fmt.Errorf("error getting user by login: %w", err)
	}

	if author == "" {
		return uuid.Nil, fmt.Errorf("author is empty")
	}
	if title == "" {
		return uuid.Nil, fmt.Errorf("description is empty")
	}
	if description == "" {
		return uuid.Nil, fmt.Errorf("genre is empty")
	}
	if genre == "" {
		return uuid.Nil, fmt.Errorf("genre is empty")
	}

	sample := entity.Sample{
		Title:       title,
		Author:      author,
		Description: description,
		Genre:       genre,
		MinioKey:    fmt.Sprintf("sample_%d_%s.wav", time.Now().Unix(), title),
		PackID:      packID,
		Price:       price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	sampleID, err := s.sampleRepo.Create(ctx, sample)
	if err != nil {
		return sampleID, fmt.Errorf("failed to create sample: %w", err)
	}

	return sampleID, nil
}

func (s *Service) UploadAudio(ctx context.Context, audioFilePath string, sampleID uuid.UUID) error {
	sample, err := s.sampleRepo.GetByID(ctx, sampleID)
	if errors.Is(err, domain.ErrNotFound) {
		return err
	}
	if err != nil {
		return fmt.Errorf("failed to get sample: %w", err)
	}

	if err := s.fileRepo.CreateBucketIfNotExists(ctx, BucketName); err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	if err := s.fileRepo.UploadFile(ctx, BucketName, sample.MinioKey, audioFilePath); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (s *Service) GetSample(ctx context.Context, id uuid.UUID) (entity.Sample, error) {
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

func (s *Service) UpdateSample(ctx context.Context, id uuid.UUID, packID *uuid.UUID, title, author, description, genre *string, price *int, size *int64, duration *float64) (entity.Sample, error) {
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
	if price != nil {
		existing.Price = *price
	}

	if size != nil {
		existing.Size = *size
	}
	if duration != nil {
		existing.Duration = *duration
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

	return nil
}

func (s *Service) GetSampleDownloadURL(ctx context.Context, minioKey string) (string, error) {
	url, err := s.fileRepo.GetFileURL(ctx, BucketName, minioKey)
	if err != nil {
		return "", fmt.Errorf("failed to get sample download url: %w", err)
	}

	return url, nil
}

func (s *Service) CreatePack(ctx context.Context, name, description, genre, author string) (uuid.UUID, error) {
	if _, err := s.userRepo.GetUserByLogin(ctx, author); err != nil {
		return uuid.Nil, fmt.Errorf("error getting user by login: %w", err)
	}

	if name == "" {
		return uuid.Nil, fmt.Errorf("name is empty")
	}
	if description == "" {
		return uuid.Nil, fmt.Errorf("description is empty")
	}
	if genre == "" {
		return uuid.Nil, fmt.Errorf("genre is empty")
	}
	if author == "" {
		return uuid.Nil, fmt.Errorf("author is empty")
	}

	return s.packRepo.Create(ctx, entity.Pack{
		Name:        name,
		Description: description,
		Genre:       genre,
		Author:      author,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
}

func (s *Service) GetPack(ctx context.Context, id uuid.UUID) (entity.Pack, error) {
	return s.packRepo.GetByID(ctx, id)
}

func (s *Service) GetAllPacks(ctx context.Context) ([]entity.Pack, error) {
	return s.packRepo.GetAll(ctx)
}

func (s *Service) UpdatePack(ctx context.Context, id uuid.UUID, name, description, genre *string) error {
	pack, err := s.packRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed get pack by id to update it: %w", err)
	}

	if name != nil {
		pack.Name = *name
	}
	if description != nil {
		pack.Description = *description
	}
	if genre != nil {
		pack.Genre = *genre
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
