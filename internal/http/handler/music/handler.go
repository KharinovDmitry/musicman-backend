package music

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/musicman-backend/internal/http/dto"
)

type Service interface {
	GetSamples(ctx context.Context) ([]entity.Sample, error)
	GetSampleDownloadURL(ctx context.Context, minioKey string) (string, error)
	GetSample(ctx context.Context, sampleID uuid.UUID) (entity.Sample, error)
	CreateSample(ctx context.Context, sample entity.Sample, audioFilePath string) (entity.Sample, error)
	UpdateSample(ctx context.Context, id uuid.UUID, packID *uuid.UUID, title, author, description, genre *string) (entity.Sample, error)
	DeleteSample(ctx context.Context, id uuid.UUID) error

	GetAllPacks(ctx context.Context) ([]entity.Pack, error)
	GetPack(ctx context.Context, id uuid.UUID) (entity.Pack, error)
	CreatePack(ctx context.Context, pack entity.Pack) error
	UpdatePack(ctx context.Context, id uuid.UUID, name, description, genre *string) error
	DeletePack(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetSamples godoc
// @Summary Получение всех семплов (без файлов)
// @Description Получение всех семплов с информацией о них
// @Tags samples
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.SampleDTO
// @Router /api/v1/samples [get]
func (h *Handler) GetSamples(c *gin.Context) {
	samples, err := h.service.GetSamples(c.Request.Context())
	if errors.Is(err, domain.ErrNotFound) {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := make([]dto.SampleDTO, len(samples))
	for i, sample := range samples {
		downloadURL, err := h.service.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("error get download url: %s for sample id %s", err.Error(), sample.ID.String())})
			return
		}

		response[i] = dto.ToSampleDTO(sample, downloadURL)
	}

	c.JSON(http.StatusOK, response)
}

// GetSample godoc
// @Summary Получение семпла по ID
// @Description Получение семпла по ID
// @Tags samples
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sample ID"
// @Success 200 {object} dto.SampleDTO
// @Router /api/v1/samples/{id} [get]
func (h *Handler) GetSample(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sample, err := h.service.GetSample(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Sample not found"})
		return
	}

	downloadURL, err := h.service.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := dto.ToSampleDTO(sample, downloadURL)

	c.JSON(http.StatusOK, response)
}

// CreateSample godoc
// @Summary Create a new sample
// @Description Create a new sample
// @Tags samples
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Body input body  dto.SampleDTO
// @Param file formData file true "Audio file (sample)"
// @Success 201 {object} dto.SampleDTO
// @Router /samples [post]
func (h *Handler) CreateSample(c *gin.Context) {
	sampleDto := dto.SampleDTO{}
	if err := c.BindJSON(&sampleDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File is required"})
		return
	}

	if file.Header.Get("Content-Type") != "audio/wav" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only WAV files are supported"})
		return
	}

	filePath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file"})
		return
	}

	sample, err := h.service.CreateSample(c.Request.Context(), sampleDto.ToEntity(), filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	downloadURL, err := h.service.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := dto.ToSampleDTO(sample, downloadURL)

	c.JSON(http.StatusCreated, response)
}

// UpdateSample godoc
// @Summary Update music
// @Description Update audio music
// @Tags samples
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateSampleRequest true "Update data"
// @Success 200 {object} dto.SampleDTO
// @Router /samples [put]
func (h *Handler) UpdateSample(c *gin.Context) {
	var req dto.UpdateSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sample, err := h.service.UpdateSample(c.Request.Context(), req.ID, req.PackID, req.Title, req.Author, req.Description, req.Genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	downloadURL, err := h.service.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ToSampleDTO(sample, downloadURL))
}

// DeleteSample godoc
// @Summary Delete music
// @Description Delete audio music
// @Tags samples
// @Security BearerAuth
// @Param id path string true "Sample ID"
// @Success 204
// @Router /samples/{id} [delete]
func (h *Handler) DeleteSample(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.service.DeleteSample(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPacks godoc
// @Summary Get all packs
// @Description Get all music packs
// @Tags packs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.PackDTO
// @Router /packs [get]
func (h *Handler) GetPacks(c *gin.Context) {
	packs, err := h.service.GetAllPacks(c.Request.Context())
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := make([]dto.PackDTO, len(packs))
	for i, pack := range packs {
		response[i] = dto.ToPackDTO(pack)
	}

	c.JSON(http.StatusOK, response)
}

// GetPack godoc
// @Summary Get pack by ID
// @Description Get music pack by ID with samples
// @Tags packs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Pack ID"
// @Success 200 {object} dto.PackWithSamplesResponse
// @Router /packs/{id} [get]
func (h *Handler) GetPack(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	pack, err := h.service.GetPack(c.Request.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	samples, err := h.service.GetSamples(c.Request.Context())
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	packSamples := make([]dto.SampleDTO, 0)
	for _, sample := range samples {
		if sample.PackID == nil {
			continue
		}
		url, err := h.service.GetSampleDownloadURL(c.Request.Context(), sample.MinioKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if *sample.PackID == pack.ID {
			packSamples = append(packSamples, dto.ToSampleDTO(sample, url))
		}
	}

	response := dto.PackWithSamplesResponse{
		PackDTO: dto.ToPackDTO(pack),
		Samples: packSamples,
	}

	c.JSON(http.StatusOK, response)
}

// CreatePack godoc
// @Summary Create a new pack
// @Description Create a new music pack
// @Tags packs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePackRequest true "Pack data"
// @Success 201 {object} dto.PackDTO
// @Router /packs [post]
func (h *Handler) CreatePack(c *gin.Context) {
	var req dto.CreatePackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := h.service.CreatePack(c.Request.Context(), entity.Pack{
		Name:        req.Name,
		Description: req.Description,
		Genre:       req.Genre,
		Author:      req.Author,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// UpdatePack godoc
// @Summary Update pack
// @Description Update music pack
// @Tags packs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Pack ID"
// @Param request body dto.UpdatePackRequest true "Update data"
// @Success 200 {object} dto.PackDTO
// @Router /packs/{id} [put]
func (h *Handler) UpdatePack(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var req dto.UpdatePackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.service.UpdatePack(c.Request.Context(), id, req.Name, req.Description, req.Genre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeletePack godoc
// @Summary Delete pack
// @Description Delete audio music
// @Tags packs
// @Security BearerAuth
// @Param id path string true "Pack ID"
// @Success 204
// @Router /packs/{id} [delete]
func (h *Handler) DeletePack(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.service.DeletePack(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
