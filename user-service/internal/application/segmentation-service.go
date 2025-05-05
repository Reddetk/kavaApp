package application

import (
	"context"
	"errors"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"
	"user-service/internal/domain/services"

	"github.com/google/uuid"
)

// SegmentationType defines the type of segmentation algorithm to use
type SegmentationType string

const (
	SegmentationTypeRFM      SegmentationType = "rfm"
	SegmentationTypeBehavior SegmentationType = "behavior"
)

// SegmentationService orchestrates the segmentation process
type SegmentationService struct {
	userRepo        repositories.UserRepository
	segmentRepo     repositories.SegmentRepository
	metricsRepo     repositories.UserMetricsRepository
	transactionRepo repositories.TransactionRepository
	segmentationSvc services.SegmentationService
	config          SegmentationConfig
}

// SegmentationConfig holds configuration for the segmentation service
type SegmentationConfig struct {
	BatchSize          int
	DefaultSegmentType SegmentationType
}

// NewSegmentationService creates a new segmentation service instance
func NewSegmentationService(
	ur repositories.UserRepository,
	sr repositories.SegmentRepository,
	mr repositories.UserMetricsRepository,
	tr repositories.TransactionRepository,
	ss services.SegmentationService,
	config SegmentationConfig,
) *SegmentationService {
	// Set default batch size if not provided
	if config.BatchSize <= 0 {
		config.BatchSize = 100 // Default batch size
	}

	// Set default segmentation type if not provided
	if config.DefaultSegmentType == "" {
		config.DefaultSegmentType = SegmentationTypeRFM
	}

	return &SegmentationService{
		userRepo:        ur,
		segmentRepo:     sr,
		metricsRepo:     mr,
		transactionRepo: tr,
		segmentationSvc: ss,
		config:          config,
	}
}

// PerformRFMSegmentation performs RFM segmentation for all users in batches
func (s *SegmentationService) PerformRFMSegmentation(ctx context.Context) error {
	// Fetch all user metrics in batches
	var allUserMetrics []entities.UserMetrics
	offset := 0

	for {
		// Get users in batches
		users, err := s.userRepo.List(ctx, s.config.BatchSize, offset)
		if err != nil {
			return err
		}

		// If no more users, we're done collecting metrics
		if len(users) == 0 {
			break
		}

		// Collect metrics for each user
		for _, user := range users {
			metrics, err := s.metricsRepo.Get(ctx, user.ID)
			if err != nil {
				// Log error but continue with other users
				continue
			}
			allUserMetrics = append(allUserMetrics, *metrics)
		}

		// Move to the next batch
		offset += len(users)

		// If we got fewer users than the batch size, we've reached the end
		if len(users) < s.config.BatchSize {
			break
		}
	}

	// Perform RFM clustering on all collected metrics
	segments, err := s.segmentationSvc.PerformRFMClustering(allUserMetrics)
	if err != nil {
		return err
	}

	// Store the new segment definitions
	for _, segment := range segments {
		// Check if segment already exists
		existingSegments, err := s.segmentRepo.GetByType(ctx, string(segment.Type))
		if err != nil {
			return err
		}

		var exists bool
		for _, existing := range existingSegments {
			if existing.Name == segment.Name {
				// Update existing segment
				segment.ID = existing.ID // Preserve ID
				err = s.segmentRepo.Update(ctx, &segment)
				if err != nil {
					return err
				}
				exists = true
				break
			}
		}

		if !exists {
			// Create new segment
			err = s.segmentRepo.Create(ctx, &segment)
			if err != nil {
				return err
			}
		}
	}

	// Now that segments are defined, assign each user to their segment
	offset = 0
	for {
		// Get users in batches
		users, err := s.userRepo.List(ctx, s.config.BatchSize, offset)
		if err != nil {
			return err
		}

		// If no more users, we're done
		if len(users) == 0 {
			break
		}

		// Assign each user to their segment
		for _, user := range users {
			metrics, err := s.metricsRepo.Get(ctx, user.ID)
			if err != nil {
				continue
			}

			// Use the domain service to determine the appropriate segment
			segment, err := s.segmentationSvc.AssignUserToSegment(user.ID, *metrics)
			if err != nil {
				continue
			}

			// Update user metrics with the segment ID
			metrics.LastSegmentID = segment.ID
			err = s.metricsRepo.Update(ctx, metrics)
			if err != nil {
				continue
			}
		}

		// Move to the next batch
		offset += len(users)

		// If we got fewer users than the batch size, we've reached the end
		if len(users) < s.config.BatchSize {
			break
		}
	}

	return nil
}

// PerformBehaviorSegmentation performs behavioral segmentation for all users based on transactions
func (s *SegmentationService) PerformBehaviorSegmentation(ctx context.Context) error {
	// Fetch all transactions in batches
	var allTransactions []entities.Transaction
	offset := 0

	for {
		// Get users in batches for their transactions
		users, err := s.userRepo.List(ctx, s.config.BatchSize, offset)
		if err != nil {
			return err
		}

		// If no more users, we're done collecting transactions
		if len(users) == 0 {
			break
		}

		// Collect transactions for each user
		for _, user := range users {
			transactions, err := s.transactionRepo.GetByUserID(ctx, user.ID)
			if err != nil {
				// Log error but continue with other users
				continue
			}

			for _, t := range transactions {
				allTransactions = append(allTransactions, *t)
			}
		}

		// Move to the next batch
		offset += len(users)

		// If we got fewer users than the batch size, we've reached the end
		if len(users) < s.config.BatchSize {
			break
		}
	}

	// Perform behavior clustering on all collected transactions
	segments, err := s.segmentationSvc.PerformBehaviorClustering(allTransactions)
	if err != nil {
		return err
	}
	// Store the new segment definitions
	for _, segment := range segments {
		// Check if segment already exists
		existingSegments, err := s.segmentRepo.GetByType(ctx, string(segment.Type))
		if err != nil {
			return err
		}

		var exists bool
		for _, existing := range existingSegments {
			if existing.Name == segment.Name {
				// Update existing segment
				segment.ID = existing.ID // Preserve ID
				err = s.segmentRepo.Update(ctx, &segment)
				if err != nil {
					return err
				}
				exists = true
				break
			}
		}

		if !exists {
			// Create new segment
			err = s.segmentRepo.Create(ctx, &segment)
			if err != nil {
				return err
			}
		}
	}

	// Now that segments are defined, assign each user to their segment
	offset = 0
	for {
		// Get users in batches
		users, err := s.userRepo.List(ctx, s.config.BatchSize, offset)
		if err != nil {
			return err
		}

		// If no more users, we're done
		if len(users) == 0 {
			break
		}

		// Assign each user to their segment
		for _, user := range users {
			metrics, err := s.metricsRepo.Get(ctx, user.ID)
			if err != nil {
				continue
			}

			// Use the domain service to determine the appropriate segment
			segment, err := s.segmentationSvc.AssignUserToSegment(user.ID, *metrics)
			if err != nil {
				continue
			}

			// Update user metrics with the segment ID
			metrics.LastSegmentID = segment.ID
			err = s.metricsRepo.Update(ctx, metrics)
			if err != nil {
				continue
			}
		}

		// Move to the next batch
		offset += len(users)

		// If we got fewer users than the batch size, we've reached the end
		if len(users) < s.config.BatchSize {
			break
		}
	}

	return nil
}

// AssignUserToSegment assigns a single user to their appropriate segment
func (s *SegmentationService) AssignUserToSegment(ctx context.Context, userID uuid.UUID) error {
	// Get user metrics
	metrics, err := s.metricsRepo.Get(ctx, userID)
	if err != nil {
		return err
	}

	// Use domain service to determine the appropriate segment
	segment, err := s.segmentationSvc.AssignUserToSegment(userID, *metrics)
	if err != nil {
		return err
	}

	// Update user metrics with the segment ID
	metrics.LastSegmentID = segment.ID
	return s.metricsRepo.Update(ctx, metrics)
}

// GetUserSegment retrieves the current segment for a user
func (s *SegmentationService) GetUserSegment(ctx context.Context, userID uuid.UUID) (*entities.Segment, error) {
	// Get user metrics to find their segment ID
	metrics, err := s.metricsRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// If user has no segment assigned
	if metrics.LastSegmentID == uuid.Nil {
		return nil, errors.New("user has no segment assigned")
	}

	// Get the segment by ID
	return s.segmentRepo.Get(ctx, metrics.LastSegmentID)
}

// GetAllSegmentsByType retrieves all segments of a specific type
func (s *SegmentationService) GetAllSegmentsByType(ctx context.Context, segmentType string) ([]*entities.Segment, error) {
	return s.segmentRepo.GetByType(ctx, segmentType)
}

// GetSegment retrieves a segment by ID
func (s *SegmentationService) GetSegment(ctx context.Context, segmentID uuid.UUID) (*entities.Segment, error) {
	return s.segmentRepo.Get(ctx, segmentID)
}

// CreateSegment creates a new segment definition
func (s *SegmentationService) CreateSegment(ctx context.Context, segment *entities.Segment) error {
	return s.segmentRepo.Create(ctx, segment)
}

// UpdateSegment updates an existing segment definition
func (s *SegmentationService) UpdateSegment(ctx context.Context, segment *entities.Segment) error {
	return s.segmentRepo.Update(ctx, segment)
}
