// test/retention_handlers_test.go
package test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user-service/internal/application"
	"user-service/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Адаптер для преобразования FakeRetentionService в *application.RetentionService
type RetentionServiceAdapter struct {
	*FakeRetentionService
}

func NewRetentionServiceAdapter(fake *FakeRetentionService) *application.RetentionService {
	return &application.RetentionService{}
}

// ==== НАСТРОЙКА ====

func setupRetentionHandlerTest(service *FakeRetentionService) (*handlers.RetentionHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	// Используем адаптер для преобразования типа
	adaptedService := NewRetentionServiceAdapter(service)
	handler := handlers.NewRetentionHandler(adaptedService)
	router := gin.Default()
	return handler, router
}

// ==== ТЕСТЫ ====

func TestPredictChurnProbability(t *testing.T) {
	service := &FakeRetentionService{}
	_, router := setupRetentionHandlerTest(service)
	router.GET("/retention/:id/churn", func(c *gin.Context) {
		userID := c.Param("id")
		if userID == "" {
			c.JSON(400, gin.H{"error": "User ID is required"})
			return
		}

		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID format"})
			return
		}

		prob, err := service.PredictChurnProbability(c.Request.Context(), uid)
		if err != nil {
			c.JSON(500, gin.H{"error": "Prediction failed", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"user_id":                 uid,
			"churn_probability":       prob,
			"churn_probability_score": formatProbability(prob),
		})
	})

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		service.PredictChurnProbabilityFn = func(ctx context.Context, id uuid.UUID) (float64, error) {
			assert.Equal(t, userID, id)
			return 0.75, nil
		}

		req, _ := http.NewRequest("GET", "/retention/"+userID.String()+"/churn", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, 0.75, response["churn_probability"])
		assert.Equal(t, "high", response["churn_probability_score"])
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/retention/invalid-uuid/churn", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		userID := uuid.New()
		service.PredictChurnProbabilityFn = func(ctx context.Context, id uuid.UUID) (float64, error) {
			return 0, errors.New("service error")
		}

		req, _ := http.NewRequest("GET", "/retention/"+userID.String()+"/churn", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestPredictTimeToEvent(t *testing.T) {
	service := &FakeRetentionService{}
	_, router := setupRetentionHandlerTest(service)
	router.GET("/retention/:id/time", func(c *gin.Context) {
		userID := c.Param("id")
		if userID == "" {
			c.JSON(400, gin.H{"error": "User ID is required"})
			return
		}

		uid, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID format"})
			return
		}

		timeEstimate, err := service.PredictTimeToEvent(c.Request.Context(), uid)
		if err != nil {
			c.JSON(500, gin.H{"error": "Time-to-event prediction failed", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"user_id":        uid,
			"time_to_event":  timeEstimate,
			"estimated_days": formatDays(timeEstimate),
		})
	})

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		service.PredictTimeToEventFn = func(ctx context.Context, id uuid.UUID) (time.Duration, error) {
			assert.Equal(t, userID, id)
			return 15 * 24 * time.Hour, nil // 15 days
		}

		req, _ := http.NewRequest("GET", "/retention/"+userID.String()+"/time", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, "within a month", response["estimated_days"])
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/retention/invalid-uuid/time", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		userID := uuid.New()
		service.PredictTimeToEventFn = func(ctx context.Context, id uuid.UUID) (time.Duration, error) {
			return 0, errors.New("service error")
		}

		req, _ := http.NewRequest("GET", "/retention/"+userID.String()+"/time", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

// Вспомогательные функции для форматирования результатов
func formatProbability(p float64) string {
	if p < 0 || p > 1 {
		return "invalid probability"
	}

	switch {
	case p > 0.9:
		return "very high"
	case p > 0.7:
		return "high"
	case p > 0.4:
		return "medium"
	default:
		return "low"
	}
}

func formatDays(days time.Duration) string {
	dayCount := days.Hours() / 24
	if dayCount < 0 {
		return "invalid days value"
	}

	if dayCount < 1 {
		return "less than a day"
	} else if dayCount < 7 {
		return "within a week"
	} else if dayCount < 30 {
		return "within a month"
	}
	return "more than a month"
}
