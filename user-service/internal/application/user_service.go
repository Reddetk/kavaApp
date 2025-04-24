// internal/application/user_service.go
package application

import (
	"context"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo    repositories.UserRepository
	metricsRepo repositories.UserMetricsRepository
}

func NewUserService(ur repositories.UserRepository, mr repositories.UserMetricsRepository) *UserService {
	return &UserService{
		userRepo:    ur,
		metricsRepo: mr,
	}
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	// Implementation
}

func (s *UserService) CreateUser(ctx context.Context, user *entities.User) error {
	// Implementation
}

func (s *UserService) UpdateUser(ctx context.Context, user *entities.User) error {
	// Implementation
}
