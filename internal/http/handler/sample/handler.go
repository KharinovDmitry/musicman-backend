package sample

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/musicman-backend/internal/http/dto"
	"github.com/musicman-backend/internal/service"
)

type Music struct {
	sampleService *service.SampleService
}

func NewSampleHandler(sampleService *service.SampleService) *Music {
	return &Music{
		sampleService: sampleService,
	}
}

// GetSamples godoc
// @Summary Get all samples
// @Description Get all audio samples
// @Tags samples
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.SampleResponse
// @Router /samples [get]
func (h *Music) GetSamples(c *gin.Context) {
	samples, err := h.sampleService.GetAllSamples(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.SampleResponse, len(samples))
	for i, sample := range samples {
		downloadURL, _ := h.sampleService.GetSampleDownloadURL(c.Request.Context(), sample.ID)
		response[i] = dto.ToSampleResponse(&sample, downloadURL)
	}

	c.JSON(http.StatusOK, response)
}

// GetSample godoc
// @Summary Get music by ID
// @Description Get audio music by ID
// @Tags samples
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sample ID"
// @Success 200 {object} dto.SampleResponse
// @Router /samples/{id} [get]
func (h *Music) GetSample(c *gin.Context) {
	id := c.Param("id")

	sample, err := h.sampleService.GetSample(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sample not found"})
		return
	}

	downloadURL, _ := h.sampleService.GetSampleDownloadURL(c.Request.Context(), id)
	response := dto.ToSampleResponse(sample, downloadURL)

	c.JSON(http.StatusOK, response)
}

// CreateSample godoc
// @Summary Create a new music
// @Description Create a new audio music
// @Tags samples
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Sample title"
// @Param author formData string true "Sample author"
// @Param description formData string false "Sample description"
// @Param genre formData string true "Sample genre" Enums(hip_hop, rock, electronic, jazz, classical, pop, r_b)
// @Param bpm formData int false "BPM"
// @Param key formData string false "Musical key"
// @Param duration formData number true "Duration in seconds"
// @Param pack_id formData string false "Pack ID"
// @Param file formData file true "Audio file"
// @Success 201 {object} dto.SampleResponse
// @Router /samples [post]
func (h *Music) CreateSample(c *gin.Context) {
	title := c.PostForm("title")
	author := c.PostForm("author")
	description := c.PostForm("description")
	genre := c.PostForm("genre")
	bpmStr := c.PostForm("bpm")
	key := c.PostForm("key")
	durationStr := c.PostForm("duration")
	packID := c.PostForm("pack_id")

	var bpm int
	if bpmStr != "" {
		bpm, _ = strconv.Atoi(bpmStr)
	}

	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	if file.Header.Get("Content-Type") != "audio/wav" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only WAV files are supported"})
		return
	}

	filePath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	req := dto.CreateSampleRequest{
		Title:       title,
		Author:      author,
		Description: description,
		Genre:       domain.Genre(genre),
		BPM:         bpm,
		Key:         key,
		Duration:    duration,
	}

	if packID != "" {
		req.PackID = &packID
	}

	sample, err := h.sampleService.CreateSample(c.Request.Context(), req, filePath, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	downloadURL, _ := h.sampleService.GetSampleDownloadURL(c.Request.Context(), sample.ID)
	response := dto.ToSampleResponse(sample, downloadURL)

	c.JSON(http.StatusCreated, response)
}

// UpdateSample godoc
// @Summary Update music
// @Description Update audio music
// @Tags samples
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sample ID"
// @Param request body dto.UpdateSampleRequest true "Update data"
// @Success 200 {object} dto.SampleResponse
// @Router /samples/{id} [put]
func (h *Music) UpdateSample(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sample, err := h.sampleService.UpdateSample(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	downloadURL, _ := h.sampleService.GetSampleDownloadURL(c.Request.Context(), id)
	response := dto.ToSampleResponse(sample, downloadURL)

	c.JSON(http.StatusOK, response)
}

// DeleteSample godoc
// @Summary Delete music
// @Description Delete audio music
// @Tags samples
// @Security BearerAuth
// @Param id path string true "Sample ID"
// @Success 204
// @Router /samples/{id} [delete]
func (h *Music) DeleteSample(c *gin.Context) {
	id := c.Param("id")

	err := h.sampleService.DeleteSample(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Pack handlers
// GetPacks godoc
// @Summary Get all packs
// @Description Get all music packs
// @Tags packs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.PackResponse
// @Router /packs [get]
func (h *Music) GetPacks(c *gin.Context) {
	packs, err := h.sampleService.GetAllPacks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.PackResponse, len(packs))
	for i, pack := range packs {
		count, _ := h.sampleService.GetSampleCountByPack(c.Request.Context(), pack.ID)
		response[i] = dto.ToPackResponse(&pack, count)
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
func (h *Music) GetPack(c *gin.Context) {
	id := c.Param("id")

	packWithSamples, err := h.sampleService.GetPackWithSamples(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pack not found"})
		return
	}

	response := dto.PackWithSamplesResponse{
		PackResponse: dto.ToPackResponse(&packWithSamples.Pack, len(packWithSamples.Samples)),
		Samples:      make([]dto.SampleResponse, len(packWithSamples.Samples)),
	}

	for i, sample := range packWithSamples.Samples {
		downloadURL, _ := h.sampleService.GetSampleDownloadURL(c.Request.Context(), sample.ID)
		response.Samples[i] = dto.ToSampleResponse(&sample, downloadURL)
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
// @Success 201 {object} dto.PackResponse
// @Router /packs [post]
func (h *Music) CreatePack(c *gin.Context) {
	var req dto.CreatePackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pack, err := h.sampleService.CreatePack(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToPackResponse(pack, 0)
	c.JSON(http.StatusCreated, response)
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
// @Success 200 {object} dto.PackResponse
// @Router /packs/{id} [put]
func (h *Music) UpdatePack(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdatePackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pack, err := h.sampleService.UpdatePack(c.Request.Context(), id, req