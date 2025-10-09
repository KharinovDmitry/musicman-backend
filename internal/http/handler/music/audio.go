package music

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type AudioHandler struct {
	service Service
}

func NewAudioHandler(service Service) *AudioHandler {
	return &AudioHandler{
		service: service,
	}
}

// UploadAudio godoc
// @Summary Upload WAV audio file
// @Description Upload a WAV audio file with metadata
// @Tags audio
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Audio title"
// @Param description formData string false "Audio description"
// @Param author formData string true "Audio author"
// @Param genre formData string true "Audio genre" Enums(hip_hop, rock, electronic, jazz, classical, pop, r_b)
// @Param bpm formData int false "BPM"
// @Param key formData string false "Musical key"
// @Param tags formData string false "Comma-separated tags"
// @Param is_public formData bool false "Is public"
// @Param file formData file true "WAV audio file"
// @Success 201 {object} dto.AudioResponse
// @Router /audio/upload [post]
func (h *AudioHandler) UploadAudio(c *gin.Context) {
	var req dto.UploadAudioRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audio, err := h.service.UploadAudio(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	downloadURL, _ := h.service.GetDownloadURL(c.Request.Context(), audio.ID)
	response := dto.ConvertToAudioResponse(audio, downloadURL, "")

	c.JSON(http.StatusCreated, response)
}

// GetAudio godoc
// @Summary Get audio file info
// @Description Get audio file metadata by ID
// @Tags audio
// @Produce json
// @Security BearerAuth
// @Param id path string true "Audio ID"
// @Success 200 {object} dto.AudioResponse
// @Router /audio/{id} [get]
func (h *AudioHandler) GetAudio(c *gin.Context) {
	id := c.Param("id")

	audio, err := h.service.GetAudio(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audio not found"})
		return
	}

	downloadURL, _ := h.service.GetDownloadURL(c.Request.Context(), id)
	response := dto.ConvertToAudioResponse(audio, downloadURL, "")

	c.JSON(http.StatusOK, response)
}

// ListAudio godoc
// @Summary List audio files
// @Description Get paginated list of audio files with filtering
// @Tags audio
// @Produce json
// @Security BearerAuth
// @Param genre query string false "Genre filter"
// @Param author query string false "Author filter"
// @Param bpm_from query int false "Minimum BPM"
// @Param bpm_to query int false "Maximum BPM"
// @Param key query string false "Key filter"
// @Param is_public query bool false "Public filter"
// @Param tags query string false "Comma-separated tags"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /audio [get]
func (h *AudioHandler) ListAudio(c *gin.Context) {
	var filter dto.AudioFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Парсим теги
	if c.Query("tags") != "" {
		filter.Tags = strings.Split(c.Query("tags"), ",")
	}

	// Устанавливаем значения по умолчанию
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 20
	}

	audioFiles, total, err := h.service.GetAllAudio(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]dto.AudioResponse, len(audioFiles))
	for i, audio := range audioFiles {
		downloadURL, _ := h.service.GetDownloadURL(c.Request.Context(), audio.ID)
		responses[i] = dto.ConvertToAudioResponse(&audio, downloadURL, "")
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        filter.Page,
			"page_size":   filter.PageSize,
			"total":       total,
			"total_pages": (total + filter.PageSize - 1) / filter.PageSize,
		},
	})
}

// UpdateAudio godoc
// @Summary Update audio metadata
// @Description Update audio file metadata
// @Tags audio
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Audio ID"
// @Param request body dto.UpdateAudioRequest true "Update data"
// @Success 200 {object} dto.AudioResponse
// @Router /audio/{id} [put]
func (h *AudioHandler) UpdateAudio(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateAudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audio, err := h.service.UpdateAudio(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	downloadURL, _ := h.service.GetDownloadURL(c.Request.Context(), id)
	response := dto.ConvertToAudioResponse(audio, downloadURL, "")

	c.JSON(http.StatusOK, response)
}

// DeleteAudio godoc
// @Summary Delete audio file
// @Description Delete audio file and its metadata
// @Tags audio
// @Security BearerAuth
// @Param id path string true "Audio ID"
// @Success 204
// @Router /audio/{id} [delete]
func (h *AudioHandler) DeleteAudio(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteAudio(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DownloadAudio godoc
// @Summary Download audio file
// @Description Download the WAV audio file
// @Tags audio
// @Produce audio/wav
// @Security BearerAuth
// @Param id path string true "Audio ID"
// @Success 200 {file} binary
// @Router /audio/{id}/download [get]
func (h *AudioHandler) DownloadAudio(c *gin.Context) {
	id := c.Param("id")

	audio, err := h.service.GetAudio(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audio not found"})
		return
	}

	// Создаем временный файл
	tempDir := "/tmp/audio_downloads"
	os.MkdirAll(tempDir, 0755)
	tempFile := filepath.Join(tempDir, fmt.Sprintf("%s.wav", audio.ID))

	// Скачиваем из MinIO
	downloadURL, err := h.service.GetDownloadURL(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Увеличиваем счетчик прослушиваний
	h.service.IncrementPlayCount(c.Request.Context(), id)

	// Перенаправляем на presigned URL
	c.Redirect(http.StatusTemporaryRedirect, downloadURL)
}

// StreamAudio godoc
// @Summary Stream audio file
// @Description Stream audio file directly from storage
// @Tags audio
// @Produce audio/wav
// @Security BearerAuth
// @Param id path string true "Audio ID"
// @Success 200 {file} binary
// @Router /audio/{id}/stream [get]
func (h *AudioHandler) StreamAudio(c *gin.Context) {
	id := c.Param("id")

	downloadURL, err := h.service.GetDownloadURL(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Увеличиваем счетчик прослушиваний
	h.service.IncrementPlayCount(c.Request.Context(), id)

	c.Redirect(http.StatusTemporaryRedirect, downloadURL)
}

// GetAudioStats godoc
// @Summary Get audio statistics
// @Description Get statistics about audio files
// @Tags audio
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /audio/stats [get]
func (h *AudioHandler) GetAudioStats(c *gin.Context) {
	stats, err := h.service.GetAudioStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
