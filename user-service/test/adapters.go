// test/adapters.go
package test

import (
	"context"
	"time"
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

// FakeUserService реализует интерфейс, аналогичный application.UserService
type FakeUserService struct {
	CreateUserFn func(ctx context.Context, user *entities.User) error
	GetUserFn    func(ctx context.Context, id uuid.UUID) (*entities.User, error)
	UpdateUserFn func(ctx context.Context, user *entities.User) error
	DeleteUserFn func(ctx context.Context, id uuid.UUID) error
	ListUsersFn  func(ctx context.Context, limit, offset int) ([]*entities.User, error)
}

func (f *FakeUserService) CreateUser(ctx context.Context, user *entities.User) error {
	if f.CreateUserFn != nil {
		return f.CreateUserFn(ctx, user)
	}
	return nil
}

func (f *FakeUserService) GetUser(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	if f.GetUserFn != nil {
		return f.GetUserFn(ctx, id)
	}
	return nil, nil
}

func (f *FakeUserService) UpdateUser(ctx context.Context, user *entities.User) error {
	if f.UpdateUserFn != nil {
		return f.UpdateUserFn(ctx, user)
	}
	return nil
}

func (f *FakeUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if f.DeleteUserFn != nil {
		return f.DeleteUserFn(ctx, id)
	}
	return nil
}

func (f *FakeUserService) ListUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	if f.ListUsersFn != nil {
		return f.ListUsersFn(ctx, limit, offset)
	}
	return nil, nil
}

// FakeSegmentationService реализует интерфейс, аналогичный application.SegmentationService
type FakeSegmentationService struct {
	CreateSegmentFn               func(ctx context.Context, segment *entities.Segment) error
	UpdateSegmentFn               func(ctx context.Context, segment *entities.Segment) error
	GetSegmentFn                  func(ctx context.Context, id uuid.UUID) (*entities.Segment, error)
	GetAllSegmentsByTypeFn        func(ctx context.Context, segmentType string) ([]*entities.Segment, error)
	PerformRFMSegmentationFn      func(ctx context.Context) error
	PerformBehaviorSegmentationFn func(ctx context.Context) error
	AssignUserToSegmentFn         func(ctx context.Context, userID uuid.UUID) error
	GetUserSegmentFn              func(ctx context.Context, userID uuid.UUID) (*entities.Segment, error)
}

func (f *FakeSegmentationService) CreateSegment(ctx context.Context, segment *entities.Segment) error {
	if f.CreateSegmentFn != nil {
		return f.CreateSegmentFn(ctx, segment)
	}
	return nil
}

func (f *FakeSegmentationService) UpdateSegment(ctx context.Context, segment *entities.Segment) error {
	if f.UpdateSegmentFn != nil {
		return f.UpdateSegmentFn(ctx, segment)
	}
	return nil
}

func (f *FakeSegmentationService) GetSegment(ctx context.Context, id uuid.UUID) (*entities.Segment, error) {
	if f.GetSegmentFn != nil {
		return f.GetSegmentFn(ctx, id)
	}
	return nil, nil
}

func (f *FakeSegmentationService) GetAllSegmentsByType(ctx context.Context, segmentType string) ([]*entities.Segment, error) {
	if f.GetAllSegmentsByTypeFn != nil {
		return f.GetAllSegmentsByTypeFn(ctx, segmentType)
	}
	return nil, nil
}

func (f *FakeSegmentationService) PerformRFMSegmentation(ctx context.Context) error {
	if f.PerformRFMSegmentationFn != nil {
		return f.PerformRFMSegmentationFn(ctx)
	}
	return nil
}

func (f *FakeSegmentationService) PerformBehaviorSegmentation(ctx context.Context) error {
	if f.PerformBehaviorSegmentationFn != nil {
		return f.PerformBehaviorSegmentationFn(ctx)
	}
	return nil
}

func (f *FakeSegmentationService) AssignUserToSegment(ctx context.Context, userID uuid.UUID) error {
	if f.AssignUserToSegmentFn != nil {
		return f.AssignUserToSegmentFn(ctx, userID)
	}
	return nil
}

func (f *FakeSegmentationService) GetUserSegment(ctx context.Context, userID uuid.UUID) (*entities.Segment, error) {
	if f.GetUserSegmentFn != nil {
		return f.GetUserSegmentFn(ctx, userID)
	}
	return nil, nil
}

// FakeCLVService реализует интерфейс, аналогичный application.CLVService
type FakeCLVService struct {
	CalculateUserCLVFn func(ctx context.Context, userID uuid.UUID) (*entities.CLV, error)
	BatchUpdateCLVFn   func(ctx context.Context, batchSize int) error
	EstimateCLVFn      func(ctx context.Context, userID uuid.UUID, scenario string) (*entities.CLV, error)
	GetHistoricalCLVFn func(ctx context.Context, userID uuid.UUID) ([]*entities.CLVDataPoint, error)
}

func (f *FakeCLVService) CalculateUserCLV(ctx context.Context, userID uuid.UUID) (*entities.CLV, error) {
	if f.CalculateUserCLVFn != nil {
		return f.CalculateUserCLVFn(ctx, userID)
	}
	return &entities.CLV{
		UserID:       userID,
		Value:        1000.0,
		Currency:     "USD",
		CalculatedAt: time.Now(),
		Scenario:     "default",
	}, nil
}

func (f *FakeCLVService) BatchUpdateCLV(ctx context.Context, batchSize int) error {
	if f.BatchUpdateCLVFn != nil {
		return f.BatchUpdateCLVFn(ctx, batchSize)
	}
	return nil
}

func (f *FakeCLVService) EstimateCLV(ctx context.Context, userID uuid.UUID, scenario string) (*entities.CLV, error) {
	if f.EstimateCLVFn != nil {
		return f.EstimateCLVFn(ctx, userID, scenario)
	}
	return &entities.CLV{
		UserID:       userID,
		Value:        1200.0,
		Currency:     "USD",
		CalculatedAt: time.Now(),
		Scenario:     scenario,
	}, nil
}

func (f *FakeCLVService) GetHistoricalCLV(ctx context.Context, userID uuid.UUID) ([]*entities.CLVDataPoint, error) {
	if f.GetHistoricalCLVFn != nil {
		return f.GetHistoricalCLVFn(ctx, userID)
	}

	now := time.Now()
	return []*entities.CLVDataPoint{
		{
			UserID:   userID,
			Value:    800.0,
			Date:     now.AddDate(0, -6, 0),
			Scenario: "default",
		},
		{
			UserID:   userID,
			Value:    900.0,
			Date:     now.AddDate(0, -3, 0),
			Scenario: "default",
		},
		{
			UserID:   userID,
			Value:    1000.0,
			Date:     now,
			Scenario: "default",
		},
	}, nil
}

// FakeRetentionService реализует интерфейс, аналогичный application.RetentionService
type FakeRetentionService struct {
	PredictChurnProbabilityFn func(ctx context.Context, userID uuid.UUID) (float64, error)
	PredictTimeToEventFn      func(ctx context.Context, userID uuid.UUID) (time.Duration, error)
	RecalculateUserMetricsFn  func(ctx context.Context, userID uuid.UUID) error
	UpdateChurnProbabilityFn  func(ctx context.Context, userID uuid.UUID) error
}

func (f *FakeRetentionService) PredictChurnProbability(ctx context.Context, userID uuid.UUID) (float64, error) {
	if f.PredictChurnProbabilityFn != nil {
		return f.PredictChurnProbabilityFn(ctx, userID)
	}
	return 0.5, nil
}

func (f *FakeRetentionService) PredictTimeToEvent(ctx context.Context, userID uuid.UUID) (time.Duration, error) {
	if f.PredictTimeToEventFn != nil {
		return f.PredictTimeToEventFn(ctx, userID)
	}
	return 30 * 24 * time.Hour, nil
}

func (f *FakeRetentionService) RecalculateUserMetrics(ctx context.Context, userID uuid.UUID) error {
	if f.RecalculateUserMetricsFn != nil {
		return f.RecalculateUserMetricsFn(ctx, userID)
	}
	return nil
}

func (f *FakeRetentionService) UpdateChurnProbability(ctx context.Context, userID uuid.UUID) error {
	if f.UpdateChurnProbabilityFn != nil {
		return f.UpdateChurnProbabilityFn(ctx, userID)
	}
	return nil
}
