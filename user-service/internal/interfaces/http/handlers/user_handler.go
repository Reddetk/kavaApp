// internal/interfaces/http/handlers/user_handler.go
package handlers

import (
	"user-service/internal/application"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *application.UserService
}

func NewUserHandler(us *application.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	// Implementation
}
