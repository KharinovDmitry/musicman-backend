package music

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/musicman-backend/internal/domain"
	"github.com/musicman-backend/internal/domain/entity"
	"github.com/musicman-backend/internal/http/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetSamples_E2E - E2E тест для получения всех семплов
func TestGetSamples_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	sampleID1 := uuid.New()
	sampleID2 := uuid.New()
	packID := uuid.New()

	expectedSamples := []entity.Sample{
		{
			ID:          sampleID1,
			Title:       "Sample 1",
			Author:      "Author 1",
			Description: "Description 1",
			Genre:       "rock",
			Duration:    120.5,
			Size:        1024000,
			PackID:      &packID,
			MinioKey:    "key1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          sampleID2,
			Title:       "Sample 2",
			Author:      "Author 2",
			Description: "Description 2",
			Genre:       "jazz",
			Duration:    180.2,
			Size:        2048000,
			PackID:      nil,
			MinioKey:    "key2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Mock expectations
	mockService.On("GetSamples", mock.Anything).Return(expectedSamples, nil)
	mockService.On("GetSampleDownloadURL", mock.Anything, "key1").Return("http://localhost:9000/samples/key1", nil)
	mockService.On("GetSampleDownloadURL", mock.Anything, "key2").Return("http://localhost:9000/samples/key2", nil)

	// Execute
	router := gin.Default()
	router.GET("/api/v1/samples", handler.GetSamples)

	req, _ := http.NewRequest("GET", "/api/v1/samples", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.SampleDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, sampleID1, response[0].ID)
	assert.Equal(t, "Sample 1", response[0].Title)
	assert.Equal(t, "http://localhost:9000/samples/key1", response[0].DownloadURL)

	mockService.AssertExpectations(t)
}

// TestGetSamples_NotFound_E2E - E2E тест для случая когда семплы не найдены
func TestGetSamples_NotFound_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Mock expectations
	mockService.On("GetSamples", mock.Anything).Return([]entity.Sample{}, domain.ErrNotFound)

	// Execute
	router := gin.Default()
	router.GET("/api/v1/samples", handler.GetSamples)

	req, _ := http.NewRequest("GET", "/api/v1/samples", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetSamples_InternalError_E2E - E2E тест для внутренней ошибки сервера
func TestGetSamples_InternalError_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Mock expectations
	mockService.On("GetSamples", mock.Anything).Return([]entity.Sample{}, fmt.Errorf("internal error"))

	// Execute
	router := gin.Default()
	router.GET("/api/v1/samples", handler.GetSamples)

	req, _ := http.NewRequest("GET", "/api/v1/samples", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetSample_E2E - E2E тест для получения семпла по ID
func TestGetSample_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	sampleID := uuid.New()
	packID := uuid.New()
	expectedSample := entity.Sample{
		ID:          sampleID,
		Title:       "Test Sample",
		Author:      "Test Author",
		Description: "Test Description",
		Genre:       "electronic",
		Duration:    150.75,
		Size:        1536000,
		PackID:      &packID,
		MinioKey:    "test-key",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock expectations
	mockService.On("GetSample", mock.Anything, sampleID).Return(expectedSample, nil)
	mockService.On("GetSampleDownloadURL", mock.Anything, "test-key").Return("http://localhost:9000/samples/test-key", nil)

	// Execute
	router := gin.Default()
	router.GET("/api/v1/samples/:id", handler.GetSample)

	req, _ := http.NewRequest("GET", "/api/v1/samples/"+sampleID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.SampleDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, sampleID, response.ID)
	assert.Equal(t, "Test Sample", response.Title)
	assert.Equal(t, "http://localhost:9000/samples/test-key", response.DownloadURL)

	mockService.AssertExpectations(t)
}

// TestGetSample_InvalidID_E2E - E2E тест для невалидного ID
func TestGetSample_InvalidID_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Execute
	router := gin.Default()
	router.GET("/api/v1/samples/:id", handler.GetSample)

	req, _ := http.NewRequest("GET", "/api/v1/samples/invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetSample_NotFound_E2E - E2E тест для случая когда семпл не найден
func TestGetSample_NotFound_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	sampleID := uuid.New()

	// Mock expectations
	mockService.On("GetSample", mock.Anything, sampleID).Return(entity.Sample{}, fmt.Errorf("not found"))

	// Execute
	router := gin.Default()
	router.GET("/api/v1/samples/:id", handler.GetSample)

	req, _ := http.NewRequest("GET", "/api/v1/samples/"+sampleID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// TestCreateSample_E2E - E2E тест для создания семпла
func TestCreateSample_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	sampleID := uuid.New()
	packID := uuid.New()
	createRequest := dto.SampleDTO{
		Title:       "New Sample",
		Author:      "New Author",
		Description: "New Description",
		Genre:       "hiphop",
		Duration:    200.0,
		PackID:      &packID,
	}

	expectedSample := entity.Sample{
		ID:          sampleID,
		Title:       "New Sample",
		Author:      "New Author",
		Description: "New Description",
		Genre:       "hiphop",
		Duration:    200.0,
		Size:        2048000,
		PackID:      &packID,
		MinioKey:    "new-key",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock expectations
	mockService.On("CreateSample", mock.Anything, mock.AnythingOfType("entity.Sample"), mock.AnythingOfType("string")).
		Return(expectedSample, nil)
	mockService.On("GetSampleDownloadURL", mock.Anything, "new-key").Return("http://localhost:9000/samples/new-key", nil)

	// Execute
	router := gin.Default()
	router.POST("/samples", handler.CreateSample)

	// Create multipart form with file and JSON
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add JSON part
	jsonPart, _ := writer.CreateFormField("sample")
	jsonData, _ := json.Marshal(createRequest)
	if _, err := jsonPart.Write(jsonData); err != nil {
		t.Fatal(err)
	}

	// Add file part
	filePart, _ := writer.CreateFormFile("file", "test.wav")
	if _, err := filePart.Write([]byte("fake wav content")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/samples", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.SampleDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, sampleID, response.ID)
	assert.Equal(t, "New Sample", response.Title)

	mockService.AssertExpectations(t)
}

// TestCreateSample_InvalidInput_E2E - E2E тест для невалидных входных данных
func TestCreateSample_InvalidInput_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Execute
	router := gin.Default()
	router.POST("/samples", handler.CreateSample)

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/samples", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateSample_E2E - E2E тест для обновления семпла
func TestUpdateSample_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	sampleID := uuid.New()
	packID := uuid.New()
	title := "Updated Title"
	author := "Updated Author"
	description := "Updated Description"
	genre := "updated-genre"

	updateRequest := dto.UpdateSampleRequest{
		ID:          sampleID,
		Title:       &title,
		Author:      &author,
		Description: &description,
		Genre:       &genre,
		PackID:      &packID,
	}

	updatedSample := entity.Sample{
		ID:          sampleID,
		Title:       "Updated Title",
		Author:      "Updated Author",
		Description: "Updated Description",
		Genre:       "updated-genre",
		Duration:    150.75,
		Size:        1536000,
		PackID:      &packID,
		MinioKey:    "updated-key",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock expectations
	mockService.On("UpdateSample", mock.Anything, sampleID, &packID, &title, &author, &description, &genre).
		Return(updatedSample, nil)
	mockService.On("GetSampleDownloadURL", mock.Anything, "updated-key").Return("http://localhost:9000/samples/updated-key", nil)

	// Execute
	router := gin.Default()
	router.PUT("/samples", handler.UpdateSample)

	jsonData, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/samples", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.SampleDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", response.Title)
	assert.Equal(t, "Updated Author", response.Author)

	mockService.AssertExpectations(t)
}

// TestUpdateSample_InvalidInput_E2E - E2E тест для невалидных данных обновления
func TestUpdateSample_InvalidInput_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Execute
	router := gin.Default()
	router.PUT("/samples", handler.UpdateSample)

	// Invalid JSON
	req, _ := http.NewRequest("PUT", "/samples", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestDeleteSample_E2E - E2E тест для удаления семпла
func TestDeleteSample_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	sampleID := uuid.New()

	// Mock expectations
	mockService.On("DeleteSample", mock.Anything, sampleID).Return(nil)

	// Execute
	router := gin.Default()
	router.DELETE("/samples/:id", handler.DeleteSample)

	req, _ := http.NewRequest("DELETE", "/samples/"+sampleID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

// TestDeleteSample_InvalidID_E2E - E2E тест для удаления с невалидным ID
func TestDeleteSample_InvalidID_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Execute
	router := gin.Default()
	router.DELETE("/samples/:id", handler.DeleteSample)

	req, _ := http.NewRequest("DELETE", "/samples/invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetPacks_E2E - E2E тест для получения всех паков
func TestGetPacks_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	pack1ID := uuid.New()
	pack2ID := uuid.New()

	expectedPacks := []entity.Pack{
		{
			ID:          pack1ID,
			Name:        "Pack 1",
			Description: "Description 1",
			Genre:       "rock",
			Author:      "Author 1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          pack2ID,
			Name:        "Pack 2",
			Description: "Description 2",
			Genre:       "electronic",
			Author:      "Author 2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Mock expectations
	mockService.On("GetAllPacks", mock.Anything).Return(expectedPacks, nil)

	// Execute
	router := gin.Default()
	router.GET("/packs", handler.GetPacks)

	req, _ := http.NewRequest("GET", "/packs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.PackDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Pack 1", response[0].Name)
	assert.Equal(t, "Pack 2", response[1].Name)

	mockService.AssertExpectations(t)
}

// TestGetPacks_NotFound_E2E - E2E тест для случая когда паки не найдены
func TestGetPacks_NotFound_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Mock expectations
	mockService.On("GetAllPacks", mock.Anything).Return([]entity.Pack{}, domain.ErrNotFound)

	// Execute
	router := gin.Default()
	router.GET("/packs", handler.GetPacks)

	req, _ := http.NewRequest("GET", "/packs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetPack_E2E - E2E тест для получения пака по ID
func TestGetPack_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	packID := uuid.New()
	sampleID1 := uuid.New()
	sampleID2 := uuid.New()

	expectedPack := entity.Pack{
		ID:          packID,
		Name:        "Test Pack",
		Description: "Test Description",
		Genre:       "rock",
		Author:      "Test Author",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	expectedSamples := []entity.Sample{
		{
			ID:          sampleID1,
			Title:       "Sample 1",
			Author:      "Author 1",
			Description: "Description 1",
			Genre:       "rock",
			Duration:    120.5,
			Size:        1024000,
			PackID:      &packID,
			MinioKey:    "key1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          sampleID2,
			Title:       "Sample 2",
			Author:      "Author 2",
			Description: "Description 2",
			Genre:       "rock",
			Duration:    180.2,
			Size:        2048000,
			PackID:      &packID,
			MinioKey:    "key2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Mock expectations
	mockService.On("GetPack", mock.Anything, packID).Return(expectedPack, nil)
	mockService.On("GetSamples", mock.Anything).Return(expectedSamples, nil)
	mockService.On("GetSampleDownloadURL", mock.Anything, "key1").Return("http://localhost:9000/samples/key1", nil)
	mockService.On("GetSampleDownloadURL", mock.Anything, "key2").Return("http://localhost:9000/samples/key2", nil)

	// Execute
	router := gin.Default()
	router.GET("/packs/:id", handler.GetPack)

	req, _ := http.NewRequest("GET", "/packs/"+packID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.PackWithSamplesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Pack", response.Name)
	assert.Len(t, response.Samples, 2)
	assert.Equal(t, "Sample 1", response.Samples[0].Title)

	mockService.AssertExpectations(t)
}

// TestGetPack_InvalidID_E2E - E2E тест для невалидного ID пака
func TestGetPack_InvalidID_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Execute
	router := gin.Default()
	router.GET("/packs/:id", handler.GetPack)

	req, _ := http.NewRequest("GET", "/packs/invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetPack_NotFound_E2E - E2E тест для случая когда пак не найден
func TestGetPack_NotFound_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	packID := uuid.New()

	// Mock expectations
	mockService.On("GetPack", mock.Anything, packID).Return(entity.Pack{}, domain.ErrNotFound)

	// Execute
	router := gin.Default()
	router.GET("/packs/:id", handler.GetPack)

	req, _ := http.NewRequest("GET", "/packs/"+packID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// TestCreatePack_E2E - E2E тест для создания пака
func TestCreatePack_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	createRequest := dto.CreatePackRequest{
		Name:        "New Pack",
		Description: "New Description",
		Genre:       "electronic",
		Author:      "New Author",
	}

	expectedPack := entity.Pack{
		Name:        "New Pack",
		Description: "New Description",
		Genre:       "electronic",
		Author:      "New Author",
	}

	// Mock expectations
	mockService.On("CreatePack", mock.Anything, expectedPack).Return(nil)

	// Execute
	router := gin.Default()
	router.POST("/packs", handler.CreatePack)

	jsonData, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/packs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

// TestCreatePack_InvalidInput_E2E - E2E тест для невалидных данных создания пака
func TestCreatePack_InvalidInput_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Execute
	router := gin.Default()
	router.POST("/packs", handler.CreatePack)

	// Missing required fields
	invalidRequest := map[string]interface{}{
		"name": "Only name provided", // missing genre and author
	}
	jsonData, _ := json.Marshal(invalidRequest)

	req, _ := http.NewRequest("POST", "/packs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdatePack_E2E - E2E тест для обновления пака
func TestUpdatePack_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	packID := uuid.New()
	name := "Updated Pack"
	description := "Updated Description"
	genre := "updated-genre"

	updateRequest := dto.UpdatePackRequest{
		Name:        &name,
		Description: &description,
		Genre:       &genre,
	}

	// Mock expectations
	mockService.On("UpdatePack", mock.Anything, packID, &name, &description, &genre).Return(nil)

	// Execute
	router := gin.Default()
	router.PUT("/packs/:id", handler.UpdatePack)

	jsonData, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/packs/"+packID.String(), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

// TestUpdatePack_InvalidID_E2E - E2E тест для обновления с невалидным ID
func TestUpdatePack_InvalidID_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Execute
	router := gin.Default()
	router.PUT("/packs/:id", handler.UpdatePack)

	updateRequest := dto.UpdatePackRequest{
		Name: stringPtr("Updated Name"),
	}
	jsonData, _ := json.Marshal(updateRequest)

	req, _ := http.NewRequest("PUT", "/packs/invalid-uuid", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestDeletePack_E2E - E2E тест для удаления пака
func TestDeletePack_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Test data
	packID := uuid.New()

	// Mock expectations
	mockService.On("DeletePack", mock.Anything, packID).Return(nil)

	// Execute
	router := gin.Default()
	router.DELETE("/packs/:id", handler.DeletePack)

	req, _ := http.NewRequest("DELETE", "/packs/"+packID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

// TestDeletePack_InvalidID_E2E - E2E тест для удаления с невалидным ID
func TestDeletePack_InvalidID_E2E(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := New(mockService)

	// Execute
	router := gin.Default()
	router.DELETE("/packs/:id", handler.DeletePack)

	req, _ := http.NewRequest("DELETE", "/packs/invalid-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
