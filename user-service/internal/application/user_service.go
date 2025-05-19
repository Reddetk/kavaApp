// internal/application/user_service.go
package application

import (
	"context"
	"errors"
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
	// Проверяем, существует ли пользователь с таким email
	exists, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	// Генерируем UUID для нового пользователя
	user.ID = uuid.New()

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// Create initial metrics for new user with default values
	metrics := &entities.UserMetrics{
		UserID:    user.ID,
		Recency:   0,
		Frequency: 0,
		Monetary:  0,
		TBP:       0,
		AvgCheck:  0,
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
