// internal/interfaces/http/handlers/segment_handler.go
package handlers

import (
	"user-service/internal/application"

	"github.com/gin-gonic/gin"
)

type SegmentHandler struct {
	segmentationService *application.SegmentationService
}

func NewSegmentHandler(ss *application.SegmentationService) *SegmentHandler {
	return &SegmentHandler{segmentationService: ss}
}

func (h *SegmentHandler) PerformSegmentation(c *gin.Context) {
	// Implementation
}
