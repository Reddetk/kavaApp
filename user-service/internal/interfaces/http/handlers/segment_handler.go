// internal/interfaces/http/handlers/segment_handler.go
package handlers

import (
	"user-service/internal/application"
	"user-service/internal/domain/entities"
	"user-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SegmentHandler struct {
	segmentationService *application.SegmentationService
	logg                *logger.Logger
}

func NewSegmentHandler(ss *application.SegmentationService) *SegmentHandler {
	if ss == nil {
		panic("segmentation service cannot be nil")
	}
	logg := logger.NewLogger("APIsegment")
	if logg == nil {
		panic("failed to initialize logger")
	}
	logg.Info("Logger APIsegment initialized")
	return &SegmentHandler{segmentationService: ss, logg: logg}
}

func (h *SegmentHandler) CreateSegment(c *gin.Context) {
	if c == nil {
		h.logg.Error("Nil context received")
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logg.Error("Invalid create segment request")
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if input.Name == "" || input.Type == "" {
		h.logg.Error("Empty name or type received")
		c.JSON(400, gin.H{"error": "Name and type cannot be empty"})
		return
	}

	segment := &entities.Segment{
		ID:   uuid.New(),
		Name: input.Name,
		Type: input.Type,
	}

	if err := h.segmentationService.CreateSegment(c.Request.Context(), segment); err != nil {
		h.logg.Errorf("Failed to create segment: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create segment"})
		return
	}

	c.JSON(201, gin.H{"message": "Segment created", "id": segment.ID})
}

func (h *SegmentHandler) UpdateSegment(c *gin.Context) {
	if c == nil {
		h.logg.Error("Nil context received")
		return
	}

	var input struct {
		ID   string `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logg.Error("Invalid update request")
		c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if input.Name == "" || input.Type == "" {
		h.logg.Error("Empty name or type received")
		c.JSON(400, gin.H{"error": "Name and type cannot be empty"})
		return
	}

	segmentID, err := uuid.Parse(input.ID)
	if err != nil {
		h.logg.Errorf("Invalid UUID: %v", err)
		c.JSON(400, gin.H{"error": "Invalid UUID"})
		return
	}

	segment := &entities.Segment{
		ID:   segmentID,
		Name: input.Name,
		Type: input.Type,
	}

	if err := h.segmentationService.UpdateSegment(c.Request.Context(), segment); err != nil {
		h.logg.Errorf("Failed to update segment: %v", err)
		c.JSON(500, gin.H{"error": "Failed to update segment"})
		return
	}

	c.JSON(200, gin.H{"message": "Segment updated"})
}

func (h *SegmentHandler) PerformRFMSegmentation(c *gin.Context) {
	if c == nil {
		h.logg.Error("Nil context received")
		return
	}

	h.logg.Info("Starting RFM segmentation")

	if err := h.segmentationService.PerformRFMSegmentation(c.Request.Context()); err != nil {
		h.logg.Errorf("RFM segmentation failed: %v", err)
		c.JSON(500, gin.H{"error": "RFM segmentation failed", "details": err.Error()})
		return
	}

	h.logg.Info("RFM segmentation completed successfully")
	c.JSON(200, gin.H{"message": "RFM segmentation completed"})
}

func (h *SegmentHandler) PerformBehaviorSegmentation(c *gin.Context) {
	if c == nil {
		h.logg.Error("Nil context received")
		return
	}

	h.logg.Info("Starting behavior segmentation")

	if err := h.segmentationService.PerformBehaviorSegmentation(c.Request.Context()); err != nil {
		h.logg.Errorf("Behavior segmentation failed: %v", err)
		c.JSON(500, gin.H{"error": "Behavior segmentation failed", "details": err.Error()})
		return
	}

	h.logg.Info("Behavior segmentation completed successfully")
	c.JSON(200, gin.H{"message": "Behavior segmentation completed"})
}

func (h *SegmentHandler) AssignUserToSegment(c *gin.Context) {
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
		h.logg.Errorf("Invalid UUID: %v", err)
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := h.segmentationService.AssignUserToSegment(c.Request.Context(), uid); err != nil {
		h.logg.Errorf("Failed to assign user to segment: %v", err)
		c.JSON(500, gin.H{"error": "Assignment failed", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User assigned to segment"})
}

func (h *SegmentHandler) GetUserSegment(c *gin.Context) {
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
		h.logg.Errorf("Invalid UUID: %v", err)
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	segment, err := h.segmentationService.GetUserSegment(c.Request.Context(), uid)
	if err != nil {
		h.logg.Errorf("Segment not found or error: %v", err)
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	if segment == nil {
		h.logg.Error("No segment found for user")
		c.JSON(404, gin.H{"error": "No segment found for user"})
		return
	}

	c.JSON(200, segment)
}

func (h *SegmentHandler) GetAllSegmentsByType(c *gin.Context) {
	if c == nil {
		h.logg.Error("Nil context received")
		return
	}

	segType := c.Query("type")
	if segType == "" {
		h.logg.Error("Empty segment type received")
		c.JSON(400, gin.H{"error": "Segment type is required"})
		return
	}

	segments, err := h.segmentationService.GetAllSegmentsByType(c.Request.Context(), segType)
	if err != nil {
		h.logg.Errorf("Failed to get segments: %v", err)
		c.JSON(500, gin.H{"error": "Failed to get segments"})
		return
	}

	if segments == nil {
		segments = make([]*entities.Segment, 0)
	}

	c.JSON(200, segments)
}

func (h *SegmentHandler) GetSegment(c *gin.Context) {
	if c == nil {
		h.logg.Error("Nil context received")
		return
	}

	id := c.Param("id")
	if id == "" {
		h.logg.Error("Empty segment ID received")
		c.JSON(400, gin.H{"error": "Segment ID is required"})
		return
	}

	segmentID, err := uuid.Parse(id)
	if err != nil {
		h.logg.Errorf("Invalid segment ID: %v", err)
		c.JSON(400, gin.H{"error": "Invalid segment ID"})
		return
	}

	segment, err := h.segmentationService.GetSegment(c.Request.Context(), segmentID)
	if err != nil {
		h.logg.Errorf("Segment not found: %v", err)
		c.JSON(404, gin.H{"error": "Segment not found"})
		return
	}

	if segment == nil {
		h.logg.Error("No segment found")
		c.JSON(404, gin.H{"error": "Segment not found"})
		return
	}

	c.JSON(200, segment)
}
