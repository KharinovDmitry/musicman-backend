package purchase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
)

type PurchaseRepository interface {
	Create(ctx context.Context, purchase entity.Purchase) (uuid.UUID, error)
	GetByUserAndSample(ctx context.Context, userUUID, sampleID uuid.UUID) (entity.Purchase, error)
	GetByUser(ctx context.Context, userUUID uuid.UUID) ([]entity.Purchase, error)
}

type SampleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (entity.Sample, error)
}

type UrlGetter interface {
	GetSampleDownloadURL(ctx context.Context, minioKey string) (string, error)
}

type UserRepository interface {
	GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (entity.User, error)
	UpdateUserBalance(ctx context.Context, userUUID uuid.UUID, amount int) error
}

type Service struct {
	purchaseRepo PurchaseRepository
	sampleRepo   SampleRepository
	userRepo     UserRepository
	url          UrlGetter
}

func New(purchaseRepo PurchaseRepository, sampleRepo SampleRepository, userRepo UserRepository, url UrlGetter) *Service {
	return &Service{
		purchaseRepo: purchaseRepo,
		sampleRepo:   sampleRepo,
		userRepo:     userRepo,
		url:          url,
	}
}

func (s *Service) PurchaseSample(ctx context.Context, userUUID, sampleID uuid.UUID) (entity.Purchase, error) {
	// 1. Получить семпл по ID
	sample, err := s.sampleRepo.GetByID(ctx, sampleID)
	if errors.Is(err, domain.ErrNotFound) {
		return entity.Purchase{}, err
	}
	if err != nil {
		return entity.Purchase{}, fmt.Errorf("failed to get sample: %w", err)
	}

	// 2. Проверить, не куплен ли уже семпл
	_, err = s.purchaseRepo.GetByUserAndSample(ctx, userUUID, sampleID)
	if err == nil {
		return entity.Purchase{}, domain.ErrAlreadyPurchased
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return entity.Purchase{}, fmt.Errorf("failed to check purchase: %w", err)
	}

	// 3. Получить пользователя
	user, err := s.userRepo.GetUserByUUID(ctx, userUUID)
	if errors.Is(err, domain.ErrNotFound) {
		return entity.Purchase{}, err
	}
	if err != nil {
		return entity.Purchase{}, fmt.Errorf("failed to get user: %w", err)
	}

	if sample.Price == 0 {
		return entity.Purchase{}, domain.ErrSampleIsFree
	}

	// 4. Проверить баланс токенов
	if user.Tokens < sample.Price {
		return entity.Purchase{}, domain.ErrInsufficientTokens
	}

	sampleURL, err := s.url.GetSampleDownloadURL(ctx, sample.MinioKey)
	if err != nil {
		return entity.Purchase{}, fmt.Errorf("failed to get sample download url: %w", err)
	}

	// 5. Создать покупку
	purchase := entity.Purchase{
		ID:          uuid.New(),
		UserUUID:    userUUID,
		SampleID:    sampleID,
		Price:       sample.Price,
		CreatedAt:   time.Now(),
		Sample:      &sample,
		ListenURL:   sampleURL,
		DownloadURL: sampleURL,
	}

	purchaseID, err := s.purchaseRepo.Create(ctx, purchase)
	if err != nil {
		return entity.Purchase{}, fmt.Errorf("failed to create purchase: %w", err)
	}
	purchase.ID = purchaseID

	// 6. Списывать токены
	err = s.userRepo.UpdateUserBalance(ctx, userUUID, -sample.Price)
	if err != nil {
		return entity.Purchase{}, fmt.Errorf("failed to update user balance: %w", err)
	}

	return purchase, nil
}

func (s *Service) GetUserPurchases(ctx context.Context, userUUID uuid.UUID) ([]entity.Purchase, error) {
	purchases, err := s.purchaseRepo.GetByUser(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user purchases: %w", err)
	}

	return purchases, nil
}

func (s *Service) IsPurchased(ctx context.Context, userUUID, sampleID uuid.UUID) (bool, error) {
	_, err := s.purchaseRepo.GetByUserAndSample(ctx, userUUID, sampleID)
	if errors.Is(err, domain.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check purchase: %w", err)
	}

	return true, nil
}
