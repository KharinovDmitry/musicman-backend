package purchase

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/constant"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/musicman-backend/internal/http/dto"
)

type PurchaseService interface {
	PurchaseSample(ctx context.Context, userUUID, sampleID uuid.UUID) (entity.Purchase, error)
	GetUserPurchases(ctx context.Context, userUUID uuid.UUID) ([]entity.Purchase, error)
	IsPurchased(ctx context.Context, userUUID, sampleID uuid.UUID) (bool, error)
}

type SampleService interface {
	GetSample(ctx context.Context, sampleID uuid.UUID) (entity.Sample, error)
	GetSampleDownloadURL(ctx context.Context, minioKey string) (string, error)
}

type Handler struct {
	purchaseService PurchaseService
	sampleService   SampleService
}

func New(purchaseService PurchaseService, sampleService SampleService) *Handler {
	return &Handler{
		purchaseService: purchaseService,
		sampleService:   sampleService,
	}
}

// PurchaseSample godoc
// @Summary Покупка семпла
// @Description Покупает семпл за токены пользователя. После покупки семпл можно скачивать неограниченное количество раз. Если семпл уже куплен, возвращает ошибку 400.
// @Tags purchases
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sample ID"
// @Success 201 {object} dto.PurchaseDTO
// @Failure 400 {object} dto.ApiError
// @Failure 401 {object} dto.ApiError
// @Failure 404 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Router /samples/{id}/purchase [post]
func (h *Handler) PurchaseSample(c *gin.Context) {
	// Получить userUUID из контекста
	userUUIDStr := c.GetString(constant.CtxUserUUID)
	if userUUIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.NewApiError("user not authenticated"))
		return
	}

	userUUID, err := uuid.Parse(userUUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError("invalid user uuid"))
		return
	}

	// Получить sampleID из path параметра
	sampleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError("invalid sample id"))
		return
	}

	// Вызвать purchaseService.PurchaseSample
	purchase, err := h.purchaseService.PurchaseSample(c.Request.Context(), userUUID, sampleID)
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, dto.NewApiError(err.Error()))
		return
	}
	if errors.Is(err, domain.ErrAlreadyPurchased) {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}
	if errors.Is(err, domain.ErrInsufficientTokens) {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	// Вернуть PurchaseDTO со статусом 201
	c.JSON(http.StatusCreated, dto.ToPurchaseDTO(purchase))
}

// GetUserPurchases godoc
// @Summary Получить список всех покупок пользователя
// @Description Возвращает список всех купленных семплов текущего пользователя, отсортированный по дате покупки (от новых к старым)
// @Tags purchases
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.PurchaseDTO
// @Failure 401 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Router /purchases [get]
func (h *Handler) GetUserPurchases(c *gin.Context) {
	// Получить userUUID из контекста
	userUUIDStr := c.GetString(constant.CtxUserUUID)
	if userUUIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.NewApiError("user not authenticated"))
		return
	}

	userUUID, err := uuid.Parse(userUUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError("invalid user uuid"))
		return
	}

	// Вызвать purchaseService.GetUserPurchases
	purchases, err := h.purchaseService.GetUserPurchases(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	// Для каждой покупки получить семпл и download URL
	result := make([]dto.PurchaseDTO, len(purchases))
	for i, purchase := range purchases {
		sample, err := h.sampleService.GetSample(c.Request.Context(), purchase.SampleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
			return
		}

		downloadURL, err := h.sampleService.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
			return
		}

		sampleDTO := dto.ToSampleDTO(sample, downloadURL, downloadURL)
		purchaseDTO := dto.ToPurchaseDTO(purchase)
		purchaseDTO.Sample = &sampleDTO
		result[i] = purchaseDTO
	}

	c.JSON(http.StatusOK, result)
}
