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
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *entities.User) error {
	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// Create initial metrics for new user
	metrics := &entities.UserMetrics{
		UserID: user.ID,
	}
	if err := s.metricsRepo.Create(ctx, metrics); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *entities.User) error {
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}
