package music

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/constant"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/musicman-backend/internal/http/dto"
)

type PurchaseChecker interface {
	IsPurchased(ctx context.Context, userUUID, sampleID uuid.UUID) (bool, error)
}

type Service interface {
	GetSamples(ctx context.Context) ([]entity.Sample, error)
	GetSampleDownloadURL(ctx context.Context, minioKey string) (string, error)
	GetSample(ctx context.Context, sampleID uuid.UUID) (entity.Sample, error)
	CreateSample(ctx context.Context, author, title, description, genre string, packID *uuid.UUID, price int) (uuid.UUID, error)
	UploadAudio(ctx context.Context, audioFilePath string, sampleID uuid.UUID) error
	UpdateSample(ctx context.Context, id uuid.UUID, packID *uuid.UUID, title, author, description, genre *string, price *int, size *int64, duration *float64) (entity.Sample, error)
	DeleteSample(ctx context.Context, id uuid.UUID) error

	GetAllPacks(ctx context.Context) ([]entity.Pack, error)
	GetPack(ctx context.Context, id uuid.UUID) (entity.Pack, error)
	CreatePack(ctx context.Context, name, description, genre, author string) (uuid.UUID, error)
	UpdatePack(ctx context.Context, id uuid.UUID, name, description, genre *string) error
	DeletePack(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	service         Service
	purchaseChecker PurchaseChecker // может быть nil, если нет авторизации
}

func New(service Service, purchaseChecker PurchaseChecker) *Handler {
	return &Handler{
		service:         service,
		purchaseChecker: purchaseChecker,
	}
}

// GetSamples godoc
// @Summary Получение всех семплов, но без файлов
// @Tags samples
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.SampleDTO
// @Success 404 {object} dto.ApiError
// @Success 500 {object} dto.ApiError
// @Router /samples [get]
func (h *Handler) GetSamples(c *gin.Context) {
	samples, err := h.service.GetSamples(c.Request.Context())
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, dto.NewApiError(err.Error()))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	// Получить userUUID из контекста (если есть)
	userUUIDStr := c.GetString(constant.CtxUserUUID)
	var userUUID *uuid.UUID
	if userUUIDStr != "" {
		parsed, err := uuid.Parse(userUUIDStr)
		if err == nil {
			userUUID = &parsed
		}
	}

	response := make([]dto.SampleDTO, len(samples))
	for i, sample := range samples {
		listenURL, err := h.service.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewApiError(fmt.Sprintf("error get download url: %s for sample id %s", err.Error(), sample.ID.String())))
			return
		}

		// Определить downloadURL: только если куплен или бесплатный
		var downloadURL string
		if sample.Price == 0 {
			// Бесплатный семпл - всегда доступен для скачивания
			downloadURL = listenURL
		} else if userUUID != nil && h.purchaseChecker != nil {
			// Проверить покупку
			isPurchased, err := h.purchaseChecker.IsPurchased(c.Request.Context(), *userUUID, sample.ID)
			if err == nil && isPurchased {
				downloadURL = listenURL
			}
		}

		response[i] = dto.ToSampleDTO(sample, listenURL, downloadURL)
	}

	c.JSON(http.StatusOK, response)
}

// GetSample godoc
// @Summary Получение семпла по ID
// @Tags samples
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sample ID"
// @Success 200 {object} dto.SampleDTO
// @Success 400 {object} dto.ApiError
// @Success 404 {object} dto.ApiError
// @Success 500 {object} dto.ApiError
// @Router /samples/{id} [get]
func (h *Handler) GetSample(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	sample, err := h.service.GetSample(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewApiError(err.Error()))
		return
	}

	listenURL, err := h.service.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	// Получить userUUID из контекста (если есть)
	userUUIDStr := c.GetString(constant.CtxUserUUID)
	var downloadURL string
	if sample.Price == 0 {
		// Бесплатный семпл - всегда доступен для скачивания
		downloadURL = listenURL
	} else if userUUIDStr != "" && h.purchaseChecker != nil {
		userUUID, err := uuid.Parse(userUUIDStr)
		if err == nil {
			// Проверить покупку
			isPurchased, err := h.purchaseChecker.IsPurchased(c.Request.Context(), userUUID, sample.ID)
			if err == nil && isPurchased {
				downloadURL = listenURL
			}
		}
	}

	c.JSON(http.StatusOK, dto.ToSampleDTO(sample, listenURL, downloadURL))
}

// UploadAudio godoc
// @Summary Загрузка .wav аудио файла для созданного семпла
// @Tags samples
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Аудио файл (sample)"
// @Param id path string true "Sample ID"
// @Success 201 {object} dto.DownloadURLResponse
// @Success 400 {object} dto.ApiError
// @Success 404 {object} dto.ApiError
// @Success 500 {object} dto.ApiError
// @Router /samples/{id} [post]
func (h *Handler) UploadAudio(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	filename := file.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".wav") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Invalid file extension",
			"expected": ".wav",
			"received": filename,
		})
		return
	}

	filePath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError("failed to save file"))
		return
	}

	sample, err := h.service.GetSample(c.Request.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, dto.NewApiError(err.Error()))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	if err := h.service.UploadAudio(c.Request.Context(), filePath, sample.ID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}
	// Открываем файл для получения duration
	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError("failed to open file"))
		return
	}
	defer fileReader.Close()

	duration, err := getWAVDuration(fileReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	sizeMB := file.Size

	if _, err := h.service.UpdateSample(c.Request.Context(), id, nil, nil, nil, nil, nil, nil, &sizeMB, &duration); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	downloadURL, err := h.service.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.DownloadURLResponse{DownloadURL: downloadURL})
}

// Функция для получения длительности WAV файла
func getWAVDuration(file multipart.File) (float64, error) {
	// Сохраняем начальную позицию
	_, err := file.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	// Читаем RIFF header
	riffHeader := make([]byte, 12)
	_, err = file.Read(riffHeader)
	if err != nil {
		return 0, err
	}

	if string(riffHeader[0:4]) != "RIFF" || string(riffHeader[8:12]) != "WAVE" {
		return 0, fmt.Errorf("invalid WAV file format")
	}

	// Ищем fmt chunk
	foundFmt := false
	var byteRate uint32

	for !foundFmt {
		chunkHeader := make([]byte, 8)
		_, err := file.Read(chunkHeader)
		if err != nil {
			return 0, err
		}

		chunkID := string(chunkHeader[0:4])
		chunkSize := binary.LittleEndian.Uint32(chunkHeader[4:8])

		if chunkID == "fmt " {
			fmtData := make([]byte, chunkSize)
			_, err := file.Read(fmtData)
			if err != nil {
				return 0, err
			}

			// Получаем byteRate (байты 8-11 в fmt chunk)
			byteRate = binary.LittleEndian.Uint32(fmtData[8:12])
			foundFmt = true
		} else {
			// Пропускаем другие chunks
			_, err = file.Seek(int64(chunkSize), 1)
			if err != nil {
				return 0, err
			}
		}
	}

	// Ищем data chunk
	var dataSize uint32
	foundData := false

	for !foundData {
		chunkHeader := make([]byte, 8)
		_, err := file.Read(chunkHeader)
		if err != nil {
			return 0, err
		}

		chunkID := string(chunkHeader[0:4])
		dataSize = binary.LittleEndian.Uint32(chunkHeader[4:8])

		if chunkID == "data" {
			foundData = true
			break
		} else {
			// Пропускаем другие chunks
			_, err = file.Seek(int64(dataSize), 1)
			if err != nil {
				return 0, err
			}
		}
	}

	// Вычисляем длительность
	if byteRate == 0 {
		return 0, fmt.Errorf("invalid byte rate")
	}

	duration := float64(dataSize) / float64(byteRate)
	return duration, nil
}

// CreateSample godoc
// @Summary Создает новый семпл (аудио загружается для созданного семпла через UploadAudio эндпоинт по ID семпла)
// @Tags samples
// @Accept application/json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateSampleRequest true "Pack data"
// @Success 201 {object} dto.UUIDResponse
// @Success 400 {object} dto.ApiError
// @Success 500 {object} dto.ApiError
// @Router /samples [post]
func (h *Handler) CreateSample(c *gin.Context) {
	var sampleDto dto.CreateSampleRequest
	if err := c.BindJSON(&sampleDto); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	id, err := h.service.CreateSample(c.Request.Context(), sampleDto.Author, sampleDto.Title, sampleDto.Description, sampleDto.Genre, sampleDto.PackID, sampleDto.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.UUIDResponse{UUID: id})
}

// UpdateSample godoc
// @Summary Обновляет семпл
// @Tags samples
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateSampleRequest true "Update data"
// @Success 200 {object} dto.SampleDTO
// @Success 400 {object} dto.ApiError
// @Success 500 {object} dto.ApiError
// @Router /samples/{id} [put]
func (h *Handler) UpdateSample(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	var req dto.UpdateSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	sample, err := h.service.UpdateSample(c.Request.Context(), id, req.PackID, req.Title, req.Author, req.Description, req.Genre, req.Price, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	listenURL, err := h.service.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	// Получить userUUID из контекста (если есть)
	userUUIDStr := c.GetString(constant.CtxUserUUID)
	var downloadURL string
	if sample.Price == 0 {
		// Бесплатный семпл - всегда доступен для скачивания
		downloadURL = listenURL
	} else if userUUIDStr != "" && h.purchaseChecker != nil {
		userUUID, err := uuid.Parse(userUUIDStr)
		if err == nil {
			// Проверить покупку
			isPurchased, err := h.purchaseChecker.IsPurchased(c.Request.Context(), userUUID, sample.ID)
			if err == nil && isPurchased {
				downloadURL = listenURL
			}
		}
	}

	c.JSON(http.StatusOK, dto.ToSampleDTO(sample, listenURL, downloadURL))
}

// DeleteSample godoc
// @Summary Удаляет семпл
// @Tags samples
// @Security BearerAuth
// @Param id path string true "Sample ID"
// @Success 204
// @Success 200 {object} dto.SampleDTO
// @Success 400 {object} dto.ApiError
// @Success 500 {object} dto.ApiError
// @Router /samples/{id} [delete]
func (h *Handler) DeleteSample(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	if err := h.service.DeleteSample(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPacks godoc
// @Summary Получение всех паков (без семплов)
// @Tags packs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.PackDTO
// @Success 404 {object} dto.ApiError
// @Success 500 {object} dto.ApiError
// @Router /packs [get]
func (h *Handler) GetPacks(c *gin.Context) {
	packs, err := h.service.GetAllPacks(c.Request.Context())
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, dto.NewApiError(err.Error()))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	response := make([]dto.PackDTO, len(packs))
	for i, pack := range packs {
		response[i] = dto.ToPackDTO(pack)
	}

	c.JSON(http.StatusOK, response)
}

// GetPack godoc
// @Summary Получает пак вместе с семплами (но без аудио содержимого)
// @Tags packs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Pack ID"
// @Success 200 {object} dto.PackWithSamplesResponse
// @Success 400 {object} dto.ApiError
// @Success 404 {object} dto.ApiError
// @Success 500 {object} dto.ApiError
// @Router /packs/{id} [get]
func (h *Handler) GetPack(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	pack, err := h.service.GetPack(c.Request.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, dto.NewApiError(err.Error()))
		return
	}
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewApiError(err.Error()))
		return
	}

	samples, err := h.service.GetSamples(c.Request.Context())
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, dto.NewApiError(err.Error()))
		return
	}
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NewApiError(err.Error()))
		return
	}

	// Получить userUUID из контекста (если есть)
	userUUIDStr := c.GetString(constant.CtxUserUUID)
	var userUUID *uuid.UUID
	if userUUIDStr != "" {
		parsed, err := uuid.Parse(userUUIDStr)
		if err == nil {
			userUUID = &parsed
		}
	}

	packSamples := make([]dto.SampleDTO, 0)
	for _, sample := range samples {
		if sample.PackID == nil {
			continue
		}
		if *sample.PackID != pack.ID {
			continue
		}

		listenURL, err := h.service.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
			return
		}

		// Определить downloadURL: только если куплен или бесплатный
		var downloadURL string
		if sample.Price == 0 {
			// Бесплатный семпл - всегда доступен для скачивания
			downloadURL = listenURL
		} else if userUUID != nil && h.purchaseChecker != nil {
			// Проверить покупку
			isPurchased, err := h.purchaseChecker.IsPurchased(c.Request.Context(), *userUUID, sample.ID)
			if err == nil && isPurchased {
				downloadURL = listenURL
			}
		}

		packSamples = append(packSamples, dto.ToSampleDTO(sample, listenURL, downloadURL))
	}

	response := dto.PackWithSamplesResponse{
		PackDTO: dto.ToPackDTO(pack),
		Samples: packSamples,
	}

	c.JSON(http.StatusOK, response)
}

// CreatePack godoc
// @Summary Создает пак семплов (без семплов)
// @Tags packs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePackRequest true "Pack data"
// @Success 201 {object} dto.UUIDResponse
// @Success 400 {object} dto.ApiError
// @Success 500 {object} dto.ApiError
// @Router /packs [post]
func (h *Handler) CreatePack(c *gin.Context) {
	var req dto.CreatePackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	id, err := h.service.CreatePack(c.Request.Context(),
		req.Name,
		req.Description,
		req.Genre,
		req.Author,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.UUIDResponse{UUID: id})
}

// UpdatePack godoc
// @Summary Обновляет пак
// @Tags packs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Pack ID"
// @Param request body dto.UpdatePackRequest true "Update data"
// @Success 204
// @Success 400 {object} dto.ApiError
// @Success 500 {object} dto.ApiError
// @Router /packs/{id} [put]
func (h *Handler) UpdatePack(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}
	var req dto.UpdatePackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	if err := h.service.UpdatePack(c.Request.Context(), id, req.Name, req.Description, req.Genre); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

// DeletePack godoc
// @Summary Удаляет пак
// @Tags packs
// @Security BearerAuth
// @Param id path string true "Pack ID"
// @Success 204
// @Success 500 {object} dto.ApiError
// @Router /packs/{id} [delete]
func (h *Handler) DeletePack(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewApiError(err.Error()))
		return
	}

	if err := h.service.DeletePack(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
