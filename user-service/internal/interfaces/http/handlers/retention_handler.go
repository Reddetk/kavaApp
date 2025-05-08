// internal/domain/services/retention-handler.go
package handlers

import (
	"user-service/internal/application"
	"user-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RetentionHandler struct {
	retntionService *application.RetentionService
	logg            *logger.Logger
}

func NewRetentionHandler(rs *application.RetentionService) *RetentionHandler {
	logg := logger.NewLogger("APIretention")
	logg.Info("Logger APIretention initialized")
	return &RetentionHandler{
		retntionService: rs,
		logg:            logg,
	}
}

func (h *RetentionHandler) PredictChurnProbability(c *gin.Context) {
	if c == nil {
		h.logg.Error("Nil context received")
		return
	}

	userID := c.Param("id")
	if userID == "" {
		h.logg.Error("Empty user ID received")
		c.JSON(400, gin.H{"error": "User ID is required"})
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		h.logg.Errorf("Invalid user ID: %v", err)
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}

	if c.Request == nil || c.Request.Context() == nil {
		h.logg.Error("Invalid request context")
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	prob, err := h.retntionService.PredictChurnProbability(c.Request.Context(), uid)
	if err != nil {
		h.logg.Errorf("Failed to predict churn: %v", err)
		c.JSON(500, gin.H{"error": "Prediction failed", "details": err.Error()})
		return
	}

	if prob < 0 || prob > 1 {
		h.logg.Errorf("Invalid probability value: %v", prob)
		c.JSON(500, gin.H{"error": "Invalid probability calculation"})
		return
	}

	c.JSON(200, gin.H{
		"user_id":                 uid,
		"churn_probability":       prob,
		"churn_probability_score": formatProbability(prob),
	})
}

func (h *RetentionHandler) PredictTimeToEvent(c *gin.Context) {
	if c == nil {
		h.logg.Error("Nil context received")
		return
	}

	userID := c.Param("id")
	if userID == "" {
		h.logg.Error("Empty user ID received")
		c.JSON(400, gin.H{"error": "User ID is required"})
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		h.logg.Errorf("Invalid user ID: %v", err)
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}

	if c.Request == nil || c.Request.Context() == nil {
		h.logg.Error("Invalid request context")
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	timeEstimate, err := h.retntionService.PredictTimeToEvent(c.Request.Context(), uid)
	if err != nil {
		h.logg.Errorf("Failed to predict time to event: %v", err)
		c.JSON(500, gin.H{"error": "Time-to-event prediction failed", "details": err.Error()})
		return
	}

	if timeEstimate < 0 {
		h.logg.Errorf("Invalid time estimate value: %v", timeEstimate)
		c.JSON(500, gin.H{"error": "Invalid time estimate calculation"})
		return
	}

	c.JSON(200, gin.H{
		"user_id":        uid,
		"time_to_event":  timeEstimate,
		"estimated_days": formatDays(timeEstimate),
	})
}

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

func formatDays(days float64) string {
	if days < 0 {
		return "invalid days value"
	}

	if days < 1 {
		return "less than a day"
	} else if days < 7 {
		return "within a week"
	} else if days < 30 {
		return "within a month"
	}
	return "more than a month"
}
