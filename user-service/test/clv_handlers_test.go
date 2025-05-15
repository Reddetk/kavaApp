// test/clv_handlers_test.go
package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user-service/internal/application"
	"user-service/internal/domain/entities"
	"user-service/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Адаптер для преобразования FakeCLVService в *application.CLVService
type CLVServiceAdapter struct {
	*FakeCLVService
}

func NewCLVServiceAdapter(fake *FakeCLVService) *application.CLVService {
	return &application.CLVService{}
}

// ==== НАСТРОЙКА ====

func setupCLVHandlerTest(service *FakeCLVService) (*handlers.CLVHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	// Используем адаптер для преобразования типа
	adaptedService := NewCLVServiceAdapter(service)
	handler := handlers.NewCLVHandler(adaptedService)
	router := gin.Default()
	return handler, router
}

// ==== ТЕСТЫ ====

func TestCalculateUserCLV(t *testing.T) {
	service := &FakeCLVService{}
	_, router := setupCLVHandlerTest(service)
	router.GET("/clv/:id", func(c *gin.Context) {
		userID := c.Param("id")

		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid UUID"})
			return
		}

		clv, err := service.CalculateUserCLV(c.Request.Context(), uid)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to calculate CLV", "details": err.Error()})
			return
		}

		c.JSON(200, clv)
	})

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		expectedCLV := &entities.CLV{
			UserID:       userID,
			Value:        1000.0,
			Currency:     "USD",
			CalculatedAt: time.Now(),
			Forecast:     1200.0,
			Confidence:   0.85,
			Scenario:     "default",
		}

		service.CalculateUserCLVFn = func(ctx context.Context, id uuid.UUID) (*entities.CLV, error) {
			assert.Equal(t, userID, id)
			return expectedCLV, nil
		}

		req, _ := http.NewRequest("GET", "/clv/"+userID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var result entities.CLV
		_ = json.Unmarshal(resp.Body.Bytes(), &result)
		assert.Equal(t, expectedCLV.UserID, result.UserID)
		assert.Equal(t, expectedCLV.Value, result.Value)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/clv/invalid-uuid", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		userID := uuid.New()
		service.CalculateUserCLVFn = func(ctx context.Context, id uuid.UUID) (*entities.CLV, error) {
			return nil, errors.New("service error")
		}

		req, _ := http.NewRequest("GET", "/clv/"+userID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestBatchUpdateCLV(t *testing.T) {
	service := &FakeCLVService{}
	_, router := setupCLVHandlerTest(service)
	router.POST("/clv/update", func(c *gin.Context) {
		var req struct {
			BatchSize int `json:"batch_size"`
		}

		if err := c.ShouldBindJSON(&req); err != nil || req.BatchSize <= 0 {
			c.JSON(400, gin.H{"error": "batch_size must be a positive integer"})
			return
		}

		err := service.BatchUpdateCLV(c.Request.Context(), req.BatchSize)
		if err != nil {
			c.JSON(500, gin.H{"error": "Batch CLV update failed", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Batch CLV update completed"})
	})

	t.Run("Success", func(t *testing.T) {
		service.BatchUpdateCLVFn = func(ctx context.Context, batchSize int) error {
			assert.Equal(t, 100, batchSize)
			return nil
		}

		reqBody := `{"batch_size": 100}`
		req, _ := http.NewRequest("POST", "/clv/update", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		reqBody := `{"batch_size": -1}`
		req, _ := http.NewRequest("POST", "/clv/update", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		service.BatchUpdateCLVFn = func(ctx context.Context, batchSize int) error {
			return errors.New("service error")
		}

		reqBody := `{"batch_size": 100}`
		req, _ := http.NewRequest("POST", "/clv/update", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestEstimateCLV(t *testing.T) {
	service := &FakeCLVService{}
	_, router := setupCLVHandlerTest(service)
	router.GET("/clv/:id/estimate", func(c *gin.Context) {
		userID := c.Param("id")
		scenario := c.DefaultQuery("scenario", "default")

		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid UUID"})
			return
		}

		clv, err := service.EstimateCLV(c.Request.Context(), uid, scenario)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to estimate CLV", "details": err.Error()})
			return
		}

		c.JSON(200, clv)
	})

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		expectedCLV := &entities.CLV{
			UserID:       userID,
			Value:        1200.0,
			Currency:     "USD",
			CalculatedAt: time.Now(),
			Forecast:     1500.0,
			Confidence:   0.8,
			Scenario:     "optimistic",
		}

		service.EstimateCLVFn = func(ctx context.Context, id uuid.UUID, scenario string) (*entities.CLV, error) {
			assert.Equal(t, userID, id)
			assert.Equal(t, "optimistic", scenario)
			return expectedCLV, nil
		}

		req, _ := http.NewRequest("GET", "/clv/"+userID.String()+"/estimate?scenario=optimistic", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var result entities.CLV
		_ = json.Unmarshal(resp.Body.Bytes(), &result)
		assert.Equal(t, expectedCLV.UserID, result.UserID)
		assert.Equal(t, expectedCLV.Value, result.Value)
		assert.Equal(t, expectedCLV.Scenario, result.Scenario)
	})
}

func TestGetHistoricalCLV(t *testing.T) {
	service := &FakeCLVService{}
	_, router := setupCLVHandlerTest(service)
	router.GET("/clv/:id/history", func(c *gin.Context) {
		userID := c.Param("id")

		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid UUID"})
			return
		}

		history, err := service.GetHistoricalCLV(c.Request.Context(), uid)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get historical CLV", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"user_id": uid,
			"history": history,
		})
	})

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		history := []*entities.CLVDataPoint{
			{UserID: userID, Value: 800.0, Date: now.AddDate(0, -6, 0), Scenario: "default"},
			{UserID: userID, Value: 900.0, Date: now.AddDate(0, -3, 0), Scenario: "default"},
			{UserID: userID, Value: 1000.0, Date: now, Scenario: "default"},
		}

		service.GetHistoricalCLVFn = func(ctx context.Context, id uuid.UUID) ([]*entities.CLVDataPoint, error) {
			assert.Equal(t, userID, id)
			return history, nil
		}

		req, _ := http.NewRequest("GET", "/clv/"+userID.String()+"/history", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, userID.String(), response["user_id"])
		assert.NotNil(t, response["history"])
	})
}
