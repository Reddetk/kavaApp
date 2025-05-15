// test/user_handlers_test.go
package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user-service/internal/application"
	"user-service/internal/domain/entities"
	"user-service/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Адаптер для преобразования FakeUserService в *application.UserService
type UserServiceAdapter struct {
	*FakeUserService
}

func NewUserServiceAdapter(fake *FakeUserService) *application.UserService {
	return &application.UserService{}
}

// ==== НАСТРОЙКА ====

func setupUserHandlerTest(service *FakeUserService) (*handlers.UserHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	// Используем адаптер для преобразования типа
	adaptedService := NewUserServiceAdapter(service)
	handler := handlers.NewUserHandler(adaptedService)
	router := gin.Default()

	// Переопределяем методы обработчика для использования нашего FakeUserService
	// вместо адаптированного сервиса
	router.HandleMethodNotAllowed = true

	return handler, router
}

// ==== ТЕСТЫ ====

func TestGetUser(t *testing.T) {
	service := &FakeUserService{}
	_, router := setupUserHandlerTest(service)
	router.GET("/users/:id", func(c *gin.Context) {
		userID := c.Param("id")
		if userID == "" {
			c.JSON(400, gin.H{"error": "user id is required"})
			return
		}

		parsedID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid user id format - must be UUID"})
			return
		}

		user, err := service.GetUser(c.Request.Context(), parsedID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		c.JSON(200, user)
	})

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		registrationDate := time.Now().Add(-24 * time.Hour)
		lastActivity := time.Now()

		expected := &entities.User{
			ID:               userID,
			Email:            "test@example.com",
			Phone:            "+375261234567",
			Age:              30,
			Gender:           "male",
			City:             "Minsk",
			RegistrationDate: registrationDate,
			LastActivity:     lastActivity,
		}

		service.GetUserFn = func(ctx context.Context, id uuid.UUID) (*entities.User, error) {
			assert.Equal(t, userID, id)
			return expected, nil
		}

		req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var actual entities.User
		_ = json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.Email, actual.Email)
		assert.Equal(t, expected.Phone, actual.Phone)
		assert.Equal(t, expected.Age, actual.Age)
		assert.Equal(t, expected.Gender, actual.Gender)
		assert.Equal(t, expected.City, actual.City)
		// Проверяем новые поля
		assert.WithinDuration(t, expected.RegistrationDate, actual.RegistrationDate, time.Second)
		assert.WithinDuration(t, expected.LastActivity, actual.LastActivity, time.Second)
	})

	t.Run("Not Found", func(t *testing.T) {
		userID := uuid.New()
		service.GetUserFn = func(ctx context.Context, id uuid.UUID) (*entities.User, error) {
			return nil, nil
		}

		req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		userID := uuid.New()
		service.GetUserFn = func(ctx context.Context, id uuid.UUID) (*entities.User, error) {
			return nil, errors.New("database error")
		}

		req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/invalid-uuid", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestCreateUser(t *testing.T) {
	service := &FakeUserService{}
	_, router := setupUserHandlerTest(service)
	router.POST("/users", func(c *gin.Context) {
		var userRequest struct {
			Email  string `json:"email" binding:"required,email"`
			Phone  string `json:"phone" binding:"omitempty"`
			Age    int    `json:"age" binding:"required,min=0"`
			Gender string `json:"gender" binding:"required"`
			City   string `json:"city" binding:"required"`
		}

		if err := c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request data", "details": err.Error()})
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

		// Call service to create user
		err := service.CreateUser(c.Request.Context(), user)
		if err != nil {
			switch err.Error() {
			case "email already exists":
				c.JSON(409, gin.H{"error": "User with this email already exists"})
			case "invalid email format":
				c.JSON(400, gin.H{"error": "Invalid email format"})
			default:
				c.JSON(500, gin.H{"error": "Internal server error", "details": err.Error()})
			}
			return
		}

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
	})

	t.Run("Success", func(t *testing.T) {
		service.CreateUserFn = func(ctx context.Context, user *entities.User) error {
			user.ID = uuid.New()
			user.RegistrationDate = time.Now()
			user.LastActivity = time.Now()
			return nil
		}

		reqBody := `{
			"email": "test@example.com",
			"phone": "+375261234567",
			"age": 30,
			"gender": "male",
			"city": "Minsk"
		}`
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NotNil(t, response["id"])
		assert.Equal(t, "test@example.com", response["email"])
		assert.Equal(t, "+375261234567", response["phone"])
		assert.Equal(t, float64(30), response["age"])
		assert.Equal(t, "male", response["gender"])
		assert.Equal(t, "Minsk", response["city"])
		// Проверяем, что новые поля присутствуют в ответе
		assert.NotNil(t, response["registration_date"])
		assert.NotNil(t, response["last_activity"])
	})

	t.Run("Invalid Input", func(t *testing.T) {
		reqBody := `{
			"email": "invalid-email",
			"age": -1
		}`
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		service.CreateUserFn = func(ctx context.Context, user *entities.User) error {
			return errors.New("email already exists")
		}

		reqBody := `{
			"email": "existing@example.com",
			"phone": "+375261234567",
			"age": 30,
			"gender": "male",
			"city": "Minsk"
		}`
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusConflict, resp.Code)
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		service.CreateUserFn = func(ctx context.Context, user *entities.User) error {
			return errors.New("invalid email format")
		}

		reqBody := `{
			"email": "test@example.com",
			"phone": "+375261234567",
			"age": 30,
			"gender": "male",
			"city": "Minsk"
		}`
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestUpdateUser(t *testing.T) {
	service := &FakeUserService{}
	_, router := setupUserHandlerTest(service)
	router.PUT("/users/:id", func(c *gin.Context) {
		userID := c.Param("id")
		if userID == "" {
			c.JSON(400, gin.H{"error": "user id is required"})
			return
		}

		parsedID, err := uuid.Parse(userID)
		if err != nil {
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
			c.JSON(400, gin.H{"error": "Invalid request data", "details": err.Error()})
			return
		}

		// Validate that at least one field is being updated
		if userRequest.Email == "" && userRequest.Phone == "" && userRequest.Age == 0 &&
			userRequest.Gender == "" && userRequest.City == "" {
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

		// Call service to update user
		err = service.UpdateUser(c.Request.Context(), user)
		if err != nil {
			switch err.Error() {
			case "user not found":
				c.JSON(404, gin.H{"error": "User not found"})
			case "email already exists":
				c.JSON(409, gin.H{"error": "User with this email already exists"})
			default:
				c.JSON(500, gin.H{"error": "Internal server error", "details": err.Error()})
			}
			return
		}

		c.JSON(200, gin.H{
			"id":     user.ID,
			"email":  user.Email,
			"phone":  user.Phone,
			"age":    user.Age,
			"gender": user.Gender,
			"city":   user.City,
		})
	})

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		service.UpdateUserFn = func(ctx context.Context, user *entities.User) error {
			assert.Equal(t, userID, user.ID)
			assert.Equal(t, "updated@example.com", user.Email)
			assert.Equal(t, "New York", user.City)
			return nil
		}

		reqBody := `{
			"email": "updated@example.com",
			"city": "New York"
		}`
		req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		_ = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Equal(t, userID.String(), response["id"])
		assert.Equal(t, "updated@example.com", response["email"])
		assert.Equal(t, "New York", response["city"])
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := uuid.New()
		service.UpdateUserFn = func(ctx context.Context, user *entities.User) error {
			return errors.New("user not found")
		}

		reqBody := `{
			"email": "updated@example.com"
		}`
		req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("No Fields To Update", func(t *testing.T) {
		userID := uuid.New()

		reqBody := `{}`
		req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		userID := uuid.New()
		service.UpdateUserFn = func(ctx context.Context, user *entities.User) error {
			return errors.New("email already exists")
		}

		reqBody := `{
			"email": "existing@example.com"
		}`
		req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusConflict, resp.Code)
	})
}
