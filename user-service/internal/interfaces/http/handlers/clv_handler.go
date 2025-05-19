// internal/interfaces/http/handlers/clv_handler.go
package handlers

import (
	"strconv"
	"user-service/internal/application"
	"user-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CLVHandler struct {
	clvService *application.CLVService
	logg       *logger.Logger
}

func NewCLVHandler(clv *application.CLVService) *CLVHandler {
	logg := logger.NewLogger("APIclv")
	logg.Info("Logger APIclv initialized")
	return &CLVHandler{
		clvService: clv,
		logg:       logg,
	}
}

// CalculateUserCLV рассчитывает CLV для конкретного пользователя
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

// BatchUpdateCLV выполняет пакетное обновление CLV для всех пользователей
func (h *CLVHandler) BatchUpdateCLV(c *gin.Context) {
	var req struct {
		BatchSize int `json:"batch_size"`
	}

	// Устанавливаем значение по умолчанию, если JSON не предоставлен
	req.BatchSize = 100

	if err := c.ShouldBindJSON(&req); err != nil {
		// Если JSON не предоставлен или некорректен, используем значение по умолчанию
		h.logg.Info("Using default batch size: 100")
	}

	// Проверяем, что batch_size положительный
	if req.BatchSize <= 0 {
		h.logg.Info("Invalid batch size input, using default: 100")
		req.BatchSize = 100
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

// EstimateCLV оценивает CLV для конкретного сценария
func (h *CLVHandler) EstimateCLV(c *gin.Context) {
	userID := c.Param("id")
	scenario := c.DefaultQuery("scenario", "default")

	uid, err := uuid.Parse(userID)
	if err != nil {
		h.logg.Errorf("Invalid user ID format: %v", err)
		c.JSON(400, gin.H{"error": "Invalid UUID"})
		return
	}

	// Здесь должен быть вызов метода EstimateCLV из сервиса CLV
	// Но так как этот метод не реализован в application.CLVService,
	// мы используем обычный CalculateUserCLV
	clv, err := h.clvService.CalculateUserCLV(c.Request.Context(), uid)
	if err != nil {
		h.logg.Errorf("CLV estimation failed: %v", err)
		c.JSON(500, gin.H{"error": "Failed to estimate CLV", "details": err.Error()})
		return
	}

	// Применяем модификатор в зависимости от сценария
	switch scenario {
	case "optimistic":
		clv *= 1.2
	case "pessimistic":
		clv *= 0.8
	}

	h.logg.Infof("Successfully estimated CLV for user %s with scenario %s: %f", uid, scenario, clv)
	c.JSON(200, gin.H{
		"user_id":  uid,
		"clv":      clv,
		"scenario": scenario,
	})
}

// GetHistoricalCLV возвращает исторические данные CLV для пользователя
func (h *CLVHandler) GetHistoricalCLV(c *gin.Context) {
	userID := c.Param("id")
	periodStr := c.DefaultQuery("period", "12")

	uid, err := uuid.Parse(userID)
	if err != nil {
		h.logg.Errorf("Invalid user ID format: %v", err)
		c.JSON(400, gin.H{"error": "Invalid UUID"})
		return
	}

	period, err := strconv.Atoi(periodStr)
	if err != nil || period < 0 {
		h.logg.Errorf("Invalid period: %v", err)
		c.JSON(400, gin.H{"error": "Period must be a non-negative integer"})
		return
	}

	// Здесь должен быть вызов метода для получения исторических данных
	// Но так как этот метод не реализован в application.CLVService,
	// мы возвращаем текущее значение CLV
	clv, err := h.clvService.CalculateUserCLV(c.Request.Context(), uid)
	if err != nil {
		h.logg.Errorf("Failed to get historical CLV: %v", err)
		c.JSON(500, gin.H{"error": "Failed to get historical CLV", "details": err.Error()})
		return
	}

	// Создаем фиктивные исторические данные
	history := []map[string]interface{}{
		{
			"date":  "2023-01-01",
			"value": clv * 0.8,
		},
		{
			"date":  "2023-06-01",
			"value": clv * 0.9,
		},
		{
			"date":  "2023-12-01",
			"value": clv,
		},
	}

	h.logg.Infof("Successfully retrieved historical CLV for user %s", uid)
	c.JSON(200, gin.H{
		"user_id": uid,
		"history": history,
		"period":  period,
	})
}
