// internal/interfaces/http/handlers/user_handler.go
package handlers

import (
	"user-service/internal/application"
	"user-service/internal/domain/entities"
	"user-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *application.UserService
	logg        *logger.Logger
}

func NewUserHandler(us *application.UserService) *UserHandler {
	// Инициализация логгера
	logg := logger.NewLogger("APIuser")
	if logg == nil {
		panic("failed to initialize logger")
	}
	logg.Info("Logger API successfully initialized")

	if us == nil {
		panic("user service cannot be nil")
	}

	return &UserHandler{
		userService: us,
		logg:        logg,
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	if c == nil {
		h.logg.Error("nil context received")
		return
	}

	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		h.logg.Error("user id is required")
		c.JSON(400, gin.H{"error": "user id is required"})
		return
	}

	// Validate UUID format
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		h.logg.Errorf("invalid user id format: %v", err)
		c.JSON(400, gin.H{"error": "invalid user id format - must be UUID"})
		return
	}

	h.logg.Debugf("getting user with id: %s", parsedID)
	user, err := h.userService.GetUser(c.Request.Context(), parsedID)
	if err != nil {
		h.logg.Errorf("error getting user: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return user data
	if user == nil {
		h.logg.Info("user not found")
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	h.logg.Infof("successfully retrieved user with id: %s", parsedID)
	c.JSON(200, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	if c == nil {
		h.logg.Error("nil context received")
		return
	}

	var userRequest struct {
		Email  string `json:"email" binding:"required,email"`
		Phone  string `json:"phone" binding:"omitempty"`
		Age    int    `json:"age" binding:"required,min=0"`
		Gender string `json:"gender" binding:"required"`
		City   string `json:"city" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		h.logg.Errorf("invalid request data: %v", err)
		c.JSON(400, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Validate required fields
	if userRequest.Email == "" || userRequest.Gender == "" || userRequest.City == "" {
		h.logg.Error("missing required fields")
		c.JSON(400, gin.H{"error": "Missing required fields"})
		return
	}

	// Create new user object from request
	user := &entities.User{
		Email:  userRequest.Email,
		Phone:  userRequest.Phone,
		Age:    userRequest.Age,
		Gender: userRequest.Gender,
		City:   userRequest.City,
	}

	h.logg.Debugf("creating new user with email: %s", user.Email)

	// Call service to create user
	err := h.userService.CreateUser(c.Request.Context(), user)
	if err != nil {
		switch err.Error() {
		case "email already exists":
			h.logg.Errorf("user with email %s already exists", user.Email)
			c.JSON(409, gin.H{"error": "User with this email already exists"})
		case "invalid email format":
			h.logg.Error("invalid email format")
			c.JSON(400, gin.H{"error": "Invalid email format"})
		case "invalid phone format":
			h.logg.Error("invalid phone format")
			c.JSON(400, gin.H{"error": "Invalid phone number format"})
		case "invalid age":
			h.logg.Error("invalid age value")
			c.JSON(400, gin.H{"error": "Invalid age value"})
		case "invalid gender":
			h.logg.Error("invalid gender value")
			c.JSON(400, gin.H{"error": "Invalid gender value"})
		case "invalid city":
			h.logg.Error("invalid city value")
			c.JSON(400, gin.H{"error": "Invalid city value"})
		default:
			h.logg.Errorf("internal server error: %v", err)
			c.JSON(500, gin.H{"error": "Internal server error", "details": err.Error()})
		}
		return
	}

	if user.ID == uuid.Nil {
		h.logg.Error("user created but ID not generated")
		c.JSON(500, gin.H{"error": "Internal server error - user ID not generated"})
		return
	}

	h.logg.Infof("successfully created user with id: %s", user.ID)

	// Return created user
	c.JSON(201, gin.H{
		"id":                user.ID,
		"email":             user.Email,
		"phone":             user.Phone,
		"age":               user.Age,
		"gender":            user.Gender,
		"city":              user.City,
		"registration_date": user.RegistrationDate,
		"last_activity":     user.LastActivity,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	if c == nil {
		h.logg.Error("nil context received")
		return
	}

	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		h.logg.Error("user id is required")
		c.JSON(400, gin.H{"error": "user id is required"})
		return
	}

	// Validate UUID format
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		h.logg.Errorf("invalid user id format: %v", err)
		c.JSON(400, gin.H{"error": "invalid user id format - must be UUID"})
		return
	}

	var userRequest struct {
		Email  string `json:"email" binding:"omitempty,email"`
		Phone  string `json:"phone" binding:"omitempty"`
		Age    int    `json:"age" binding:"omitempty,min=0"`
		Gender string `json:"gender" binding:"omitempty"`
		City   string `json:"city" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		h.logg.Errorf("invalid request data: %v", err)
		c.JSON(400, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Validate that at least one field is being updated
	if userRequest.Email == "" && userRequest.Phone == "" && userRequest.Age == 0 &&
		userRequest.Gender == "" && userRequest.City == "" {
		h.logg.Error("no fields to update")
		c.JSON(400, gin.H{"error": "At least one field must be provided for update"})
		return
	}

	// Create user object for update
	user := &entities.User{
		ID:     parsedID,
		Email:  userRequest.Email,
		Phone:  userRequest.Phone,
		Age:    userRequest.Age,
		Gender: userRequest.Gender,
		City:   userRequest.City,
	}

	h.logg.Debugf("updating user with id: %s", user.ID)

	// Call service to update user
	err = h.userService.UpdateUser(c.Request.Context(), user)
	if err != nil {
		switch err.Error() {
		case "user not found":
			h.logg.Errorf("user with id %s not found", user.ID)
			c.JSON(404, gin.H{"error": "User not found"})
		case "email already exists":
			h.logg.Errorf("user with email %s already exists", user.Email)
			c.JSON(409, gin.H{"error": "User with this email already exists"})
		case "invalid email format":
			h.logg.Error("invalid email format")
			c.JSON(400, gin.H{"error": "Invalid email format"})
		case "invalid phone format":
			h.logg.Error("invalid phone format")
			c.JSON(400, gin.H{"error": "Invalid phone number format"})
		case "invalid age":
			h.logg.Error("invalid age value")
			c.JSON(400, gin.H{"error": "Invalid age value"})
		case "invalid gender":
			h.logg.Error("invalid gender value")
			c.JSON(400, gin.H{"error": "Invalid gender value"})
		case "invalid city":
			h.logg.Error("invalid city value")
			c.JSON(400, gin.H{"error": "Invalid city value"})
		default:
			h.logg.Errorf("internal server error: %v", err)
			c.JSON(500, gin.H{"error": "Internal server error", "details": err.Error()})
		}
		return
	}

	h.logg.Infof("successfully updated user with id: %s", user.ID)

	c.JSON(200, gin.H{
		"id":     user.ID,
		"email":  user.Email,
		"phone":  user.Phone,
		"age":    user.Age,
		"gender": user.Gender,
		"city":   user.City,
	})
}
