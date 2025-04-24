// internal/interfaces/http/router.go
package http

import (
	"user-service/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler, segmentHandler *handlers.SegmentHandler) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("/:id", userHandler.GetUser)
			// Другие эндпоинты
		}

		segments := v1.Group("/segments")
		{
			segments.POST("/rfm", segmentHandler.PerformSegmentation)
			// Другие эндпоинты
		}
	}

	return router
}
