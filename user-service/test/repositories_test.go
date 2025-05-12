// test/repositories_test.go
package test

import (
	"testing"
	"user-service/internal/infrastructure/postgres"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestRepositories is a test function that demonstrates how to test all repositories
func TestRepositories(t *testing.T) {
	// Setup mock database for each repository separately
	// This is important because each repository needs its own mock
	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	dbTransaction, mockTransaction, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbTransaction.Close()

	dbSegment, mockSegment, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbSegment.Close()

	dbMetrics, mockMetrics, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbMetrics.Close()

	// Initialize repositories with the mocks we just created
	userRepo := postgres.NewUserRepository(dbUser)
	transactionRepo := postgres.NewTransactionRepository(dbTransaction)
	segmentRepo := postgres.NewSegmentRepository(dbSegment)
	metricsRepo := postgres.NewUserMetricsRepository(dbMetrics)

	// Test transaction repository functions
	t.Run("Test Transaction Repository", func(t *testing.T) {
		// Test GetByUserID
		t.Run("GetByUserID", func(t *testing.T) {
			TestGetByUserIDHelper(t, transactionRepo, mockTransaction)
		})

		// Test GetByPeriod
		t.Run("GetByPeriod", func(t *testing.T) {
			TestGetByPeriodHelper(t, transactionRepo, mockTransaction)
		})

		// Test Create
		t.Run("Create", func(t *testing.T) {
			TestCreateHelper(t, transactionRepo, mockTransaction)
		})
	})

	// Test user repository functions
	t.Run("Test User Repository", func(t *testing.T) {
		// Test Get
		t.Run("Get", func(t *testing.T) {
			TestUserGetHelper(t, userRepo, mockUser)
		})

		// Test Create
		t.Run("Create", func(t *testing.T) {
			TestUserCreateHelper(t, userRepo, mockUser)
		})

		// Test Update
		t.Run("Update", func(t *testing.T) {
			TestUserUpdateHelper(t, userRepo, mockUser)
		})

		// Test Delete
		t.Run("Delete", func(t *testing.T) {
			TestUserDeleteHelper(t, userRepo, mockUser)
		})

		// Test List
		t.Run("List", func(t *testing.T) {
			TestUserListHelper(t, userRepo, mockUser)
		})
	})

	// Test segment repository functions
	t.Run("Test Segment Repository", func(t *testing.T) {
		// Test Get
		t.Run("Get", func(t *testing.T) {
			TestSegmentGetHelper(t, segmentRepo, mockSegment)
		})

		// Test Create
		t.Run("Create", func(t *testing.T) {
			TestSegmentCreateHelper(t, segmentRepo, mockSegment)
		})

		// Test Update
		t.Run("Update", func(t *testing.T) {
			TestSegmentUpdateHelper(t, segmentRepo, mockSegment)
		})

		// Test GetByType
		t.Run("GetByType", func(t *testing.T) {
			TestSegmentGetByTypeHelper(t, segmentRepo, mockSegment)
		})
	})

	// Test metrics repository functions
	t.Run("Test Metrics Repository", func(t *testing.T) {
		// Test Get
		t.Run("Get", func(t *testing.T) {
			TestUserMetricsGetHelper(t, metricsRepo, mockMetrics)
		})

		// Test Create
		t.Run("Create", func(t *testing.T) {
			TestUserMetricsCreateHelper(t, metricsRepo, mockMetrics)
		})

		// Test Update
		t.Run("Update", func(t *testing.T) {
			TestUserMetricsUpdateHelper(t, metricsRepo, mockMetrics)
		})

		// Test CalculateMetrics
		t.Run("CalculateMetrics", func(t *testing.T) {
			TestUserMetricsCalculateMetricsHelper(t, metricsRepo, mockMetrics)
		})
	})
}
