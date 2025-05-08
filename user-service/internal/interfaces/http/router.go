package http

import (
	"user-service/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *handlers.UserHandler,
	segmentHandler *handlers.SegmentHandler,
	retentionHandler *handlers.RetentionHandler,
	clvHandler *handlers.CLVHandler,
) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		// --- Пользователи ---
		users := v1.Group("/users")
		{
			users.GET("/:id", userHandler.GetUser)
			users.POST("/", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
		}

		// --- Сегменты ---
		segments := v1.Group("/segments")
		{
			segments.POST("/rfm", segmentHandler.PerformRFMSegmentation)
			segments.POST("/behavior", segmentHandler.PerformBehaviorSegmentation)
			segments.POST("/", segmentHandler.CreateSegment)
			segments.PUT("/", segmentHandler.UpdateSegment)
			segments.GET("/", segmentHandler.GetAllSegmentsByType)
			segments.GET("/:id", segmentHandler.GetSegment)
			segments.PUT("/assign/:id", segmentHandler.AssignUserToSegment)
			segments.GET("/user/:id", segmentHandler.GetUserSegment)
		}

		// --- Удержание ---
		retention := v1.Group("/retention")
		{
			retention.GET("/:id/churn", retentionHandler.PredictChurnProbability)
			retention.GET("/:id/time", retentionHandler.PredictTimeToEvent)
		}

		// --- CLV ---
		clv := v1.Group("/clv")
		{
			clv.GET("/:id", clvHandler.CalculateUserCLV)
			clv.POST("/update", clvHandler.BatchUpdateCLV)
		}
	}

	return router
}
