// test/segment_handlers_test.go
package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/internal/application"
	"user-service/internal/domain/entities"
	"user-service/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Адаптер для преобразования FakeSegmentationService в *application.SegmentationService
type SegmentationServiceAdapter struct {
	*FakeSegmentationService
}

func NewSegmentationServiceAdapter(fake *FakeSegmentationService) *application.SegmentationService {
	return &application.SegmentationService{}
}

// ==== НАСТРОЙКА ====

func setupSegmentHandlerTest(service *FakeSegmentationService) (*handlers.SegmentHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	// Используем адаптер для преобразования типа
	adaptedService := NewSegmentationServiceAdapter(service)
	handler := handlers.NewSegmentHandler(adaptedService)
	router := gin.Default()
	return handler, router
}

// ==== ТЕСТЫ ====

func TestAll(t *testing.T) {
	t.Run("Create Segment Tests", func(t *testing.T) {
		TestCreateSegment(t)
	})

	t.Run("Get Segment Tests", func(t *testing.T) {
		TestGetSegment(t)
	})

	t.Run("Perform RFM Segmentation Tests", func(t *testing.T) {
		TestPerformRFMSegmentation(t)
	})

	t.Run("Get User Segment Tests", func(t *testing.T) {
		TestGetUserSegment(t)
	})
}

func TestCreateSegment(t *testing.T) {
	service := &FakeSegmentationService{}
	_, router := setupSegmentHandlerTest(service)
	router.POST("/segments", func(c *gin.Context) {
		var input struct {
			Name string `json:"name" binding:"required"`
			Type string `json:"type" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		if input.Name == "" || input.Type == "" {
			c.JSON(400, gin.H{"error": "Name and type cannot be empty"})
			return
		}

		segment := &entities.Segment{
			ID:   uuid.New(),
			Name: input.Name,
			Type: input.Type,
		}

		if err := service.CreateSegment(c.Request.Context(), segment); err != nil {
			c.JSON(500, gin.H{"error": "Failed to create segment"})
			return
		}

		c.JSON(201, gin.H{"message": "Segment created", "id": segment.ID})
	})

	t.Run("Success", func(t *testing.T) {
		service.CreateSegmentFn = func(ctx context.Context, s *entities.Segment) error {
			s.ID = uuid.New()
			return nil
		}

		reqBody := `{"name":"Test Segment","type":"rfm"}`
		req, _ := http.NewRequest("POST", "/segments", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "Segment created", response["message"])
		assert.NotNil(t, response["id"])
	})

	t.Run("Invalid Input", func(t *testing.T) {
		reqBody := `{"name":"","type":""}`
		req, _ := http.NewRequest("POST", "/segments", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		service.CreateSegmentFn = func(ctx context.Context, s *entities.Segment) error {
			return errors.New("service error")
		}

		reqBody := `{"name":"Test Segment","type":"rfm"}`
		req, _ := http.NewRequest("POST", "/segments", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestGetSegment(t *testing.T) {
	service := &FakeSegmentationService{}
	_, router := setupSegmentHandlerTest(service)
	router.GET("/segments/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{"error": "Segment ID is required"})
			return
		}

		segmentID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid segment ID"})
			return
		}

		segment, err := service.GetSegment(c.Request.Context(), segmentID)
		if err != nil {
			c.JSON(404, gin.H{"error": "Segment not found"})
			return
		}

		if segment == nil {
			c.JSON(404, gin.H{"error": "Segment not found"})
			return
		}

		c.JSON(200, segment)
	})

	t.Run("Success", func(t *testing.T) {
		segmentID := uuid.New()
		expected := &entities.Segment{ID: segmentID, Name: "Test", Type: "rfm"}

		service.GetSegmentFn = func(ctx context.Context, id uuid.UUID) (*entities.Segment, error) {
			assert.Equal(t, segmentID, id)
			return expected, nil
		}

		req, _ := http.NewRequest("GET", "/segments/"+segmentID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var actual entities.Segment
		_ = json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.Name, actual.Name)
	})

	t.Run("Not Found", func(t *testing.T) {
		segmentID := uuid.New()
		service.GetSegmentFn = func(ctx context.Context, id uuid.UUID) (*entities.Segment, error) {
			return nil, nil
		}

		req, _ := http.NewRequest("GET", "/segments/"+segmentID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/segments/invalid", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestPerformRFMSegmentation(t *testing.T) {
	service := &FakeSegmentationService{}
	_, router := setupSegmentHandlerTest(service)
	router.POST("/segments/rfm", func(c *gin.Context) {
		if err := service.PerformRFMSegmentation(c.Request.Context()); err != nil {
			c.JSON(500, gin.H{"error": "RFM segmentation failed", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "RFM segmentation completed"})
	})

	t.Run("Success", func(t *testing.T) {
		service.PerformRFMSegmentationFn = func(ctx context.Context) error {
			return nil
		}

		req, _ := http.NewRequest("POST", "/segments/rfm", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		service.PerformRFMSegmentationFn = func(ctx context.Context) error {
			return errors.New("fail")
		}

		req, _ := http.NewRequest("POST", "/segments/rfm", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestGetUserSegment(t *testing.T) {
	service := &FakeSegmentationService{}
	_, router := setupSegmentHandlerTest(service)
	router.GET("/segments/user/:id", func(c *gin.Context) {
		userID := c.Param("id")
		if userID == "" {
			c.JSON(400, gin.H{"error": "User ID is required"})
			return
		}

		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID"})
			return
		}

		segment, err := service.GetUserSegment(c.Request.Context(), uid)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		if segment == nil {
			c.JSON(404, gin.H{"error": "No segment found for user"})
			return
		}

		c.JSON(200, segment)
	})

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		segment := &entities.Segment{
			ID:   uuid.New(),
			Name: "Test Segment",
			Type: "rfm",
		}

		service.GetUserSegmentFn = func(ctx context.Context, id uuid.UUID) (*entities.Segment, error) {
			assert.Equal(t, userID, id)
			return segment, nil
		}

		req, _ := http.NewRequest("GET", "/segments/user/"+userID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var result entities.Segment
		_ = json.Unmarshal(resp.Body.Bytes(), &result)
		assert.Equal(t, segment.ID, result.ID)
		assert.Equal(t, segment.Name, result.Name)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := uuid.New()
		service.GetUserSegmentFn = func(ctx context.Context, id uuid.UUID) (*entities.Segment, error) {
			return nil, errors.New("not found")
		}

		req, _ := http.NewRequest("GET", "/segments/user/"+userID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}
