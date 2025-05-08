// internal/domain/services/clv_handler.go
package handlers

import (
	"user-service/internal/application"
	"user-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CLVHandler struct {
	clvService *application.CLVService
	logg       *logger.Logger
}

func NewCLVSHandler(clv *application.CLVService) *CLVHandler {
	logg := logger.NewLogger("APIclv")
	logg.Info("Logger APIclv initialized")
	return &CLVHandler{
		clvService: clv,
		logg:       logg,
	}
}

func (h *CLVHandler) CalculateUserCLV(c *gin.Context) {
	userID := c.Param("id")

	uid, err := uuid.Parse(userID)
	if err != nil {
		h.logg.Errorf("Invalid user ID format: %v", err)
		c.JSON(400, gin.H{"error": "Invalid UUID"})
		return
	}

	clv, err := h.clvService.CalculateUserCLV(c.Request.Context(), uid)
	if err != nil {
		h.logg.Errorf("CLV calculation failed: %v", err)
		c.JSON(500, gin.H{"error": "Failed to calculate CLV", "details": err.Error()})
		return
	}

	h.logg.Infof("Successfully calculated CLV for user %s: %f", uid, clv)
	c.JSON(200, gin.H{
		"user_id": uid,
		"clv":     clv,
	})
}

func (h *CLVHandler) BatchUpdateCLV(c *gin.Context) {
	var req struct {
		BatchSize int `json:"batch_size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.BatchSize <= 0 {
		h.logg.Info("Invalid batch size input")
		c.JSON(400, gin.H{"error": "batch_size must be a positive integer"})
		return
	}

	err := h.clvService.BatchUpdateCLV(c.Request.Context(), req.BatchSize)
	if err != nil {
		h.logg.Errorf("Batch CLV update failed: %v", err)
		c.JSON(500, gin.H{"error": "Batch CLV update failed", "details": err.Error()})
		return
	}

	h.logg.Infof("Batch CLV update completed with batch size %d", req.BatchSize)
	c.JSON(200, gin.H{"message": "Batch CLV update completed"})
}
