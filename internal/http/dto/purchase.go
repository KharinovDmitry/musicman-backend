package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain/entity"
)

// PurchaseDTO - DTO для покупки
type PurchaseDTO struct {
	ID          uuid.UUID  `json:"id"`
	SampleID    uuid.UUID  `json:"sampleId"`
	Sample      *SampleDTO `json:"sample,omitempty"` // опционально, для списка покупок
	Price       int        `json:"price"`
	PurchasedAt time.Time  `json:"purchasedAt"`
}

// ToPurchaseDTO - конвертирует entity.Purchase в PurchaseDTO
func ToPurchaseDTO(purchase entity.Purchase) PurchaseDTO {
	sample := ToSampleDTO(*purchase.Sample, purchase.ListenURL, purchase.DownloadURL)
	return PurchaseDTO{
		ID:          purchase.ID,
		SampleID:    purchase.SampleID,
		Price:       purchase.Price,
		PurchasedAt: purchase.CreatedAt,
		Sample:      &sample,
	}
}

// PurchasesToDTO - конвертирует слайс покупок в слайс DTO
func PurchasesToDTO(purchases []entity.Purchase) []PurchaseDTO {
	result := make([]PurchaseDTO, len(purchases))
	for i, p := range purchases {
		result[i] = ToPurchaseDTO(p)
	}
	return result
}
