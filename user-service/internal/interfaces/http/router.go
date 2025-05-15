// internal/interfaces/http/router.go
package http

import (
	"user-service/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает все маршруты API для сервиса пользователей
// Принимает обработчики для различных доменных областей и возвращает настроенный роутер
func SetupRouter(
	userHandler *handlers.UserHandler,
	segmentHandler *handlers.SegmentHandler,
	clvHandler *handlers.CLVHandler,
	retentionHandler *handlers.RetentionHandler,
) *gin.Engine {
	router := gin.Default()

	// Группа API v1
	v1 := router.Group("/api/v1")
	{
		// --- Пользователи ---
		users := v1.Group("/users")
		{
			// GET /api/v1/users/:id - Получение информации о пользователе по ID
			users.GET("/:id", userHandler.GetUser)
			
			// POST /api/v1/users - Создание нового пользователя
			users.POST("/", userHandler.CreateUser)
			
			// PUT /api/v1/users/:id - Обновление информации о пользователе
			users.PUT("/:id", userHandler.UpdateUser)
		}

		// --- Сегменты ---
		segments := v1.Group("/segments")
		{
			// POST /api/v1/segments/rfm - Запуск RFM сегментации для всех пользователей
			segments.POST("/rfm", segmentHandler.PerformRFMSegmentation)
			
			// POST /api/v1/segments/behavior - Запуск поведенческой сегментации
			segments.POST("/behavior", segmentHandler.PerformBehaviorSegmentation)
			
			// POST /api/v1/segments - Создание нового сегмента
			segments.POST("/", segmentHandler.CreateSegment)
			
			// PUT /api/v1/segments - Обновление существующего сегмента
			segments.PUT("/", segmentHandler.UpdateSegment)
			
			// GET /api/v1/segments?type=X - Получение всех сегментов определенного типа
			segments.GET("/", segmentHandler.GetAllSegmentsByType)
			
			// GET /api/v1/segments/:id - Получение информации о сегменте по ID
			segments.GET("/:id", segmentHandler.GetSegment)
			
			// PUT /api/v1/segments/assign/:id - Назначение пользователя в сегмент
			segments.PUT("/assign/:id", segmentHandler.AssignUserToSegment)
			
			// GET /api/v1/segments/user/:id - Получение сегмента пользователя
			segments.GET("/user/:id", segmentHandler.GetUserSegment)
		}

		// --- Удержание ---
		retention := v1.Group("/retention")
		{
			// GET /api/v1/retention/:id/churn - Прогноз вероятности оттока для пользователя
			retention.GET("/:id/churn", retentionHandler.PredictChurnProbability)
			
			// GET /api/v1/retention/:id/time - Прогноз времени до события (оттока)
			retention.GET("/:id/time", retentionHandler.PredictTimeToEvent)
		}

		// --- CLV (Customer Lifetime Value) ---
		clv := v1.Group("/clv")
		{
			// GET /api/v1/clv/:id - Расчет CLV для пользователя
			clv.GET("/:id", clvHandler.CalculateUserCLV)
			
			// POST /api/v1/clv/update - Пакетное обновление CLV для всех пользователей
			clv.POST("/update", clvHandler.BatchUpdateCLV)
			
			// GET /api/v1/clv/:id/estimate?scenario=X - Оценка CLV по сценарию
			clv.GET("/:id/estimate", clvHandler.EstimateCLV)
			
			// GET /api/v1/clv/:id/history - Получение исторических данных CLV
			clv.GET("/:id/history", clvHandler.GetHistoricalCLV)
		}
	}

	return router
}