// test/main_test.go
package test

import (
	"context"
	"testing"
	"time"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"
	"user-service/internal/infrastructure/postgres"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestMain is a test function that demonstrates how to test repositories in main.go
func TestMain_Repositories(t *testing.T) {
	// This test demonstrates how you would test repositories in your main.go file
	// In a real application, you would typically use these tests in your CI/CD pipeline

	// Setup mock database for each repository
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

	// Initialize repositories with mock DB
	userRepo := postgres.NewUserRepository(dbUser)
	transactionRepo := postgres.NewTransactionRepository(dbTransaction)
	segmentRepo := postgres.NewSegmentRepository(dbSegment)
	metricsRepo := postgres.NewUserMetricsRepository(dbMetrics)

	// Test transaction repository functions
	t.Run("Test Transaction Repository", func(t *testing.T) {
		testTransactionRepositoryInMain(t, transactionRepo, mockTransaction)
	})

	// Test user repository functions
	t.Run("Test User Repository", func(t *testing.T) {
		testUserRepositoryInMain(t, userRepo, mockUser)
	})

	// Test segment repository functions
	t.Run("Test Segment Repository", func(t *testing.T) {
		testSegmentRepositoryInMain(t, segmentRepo, mockSegment)
	})

	// Test metrics repository functions
	t.Run("Test Metrics Repository", func(t *testing.T) {
		testUserMetricsRepositoryInMain(t, metricsRepo, mockMetrics)
	})
}

func testTransactionRepositoryInMain(t *testing.T, repo repositories.TransactionRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()

	// Test GetByUserID
	t.Run("GetByUserID", func(t *testing.T) {
		userID := uuid.New()
		transactionID := uuid.New()
		timestamp := time.Now()

		rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "timestamp", "category", "discount_applied"}).
			AddRow(transactionID, userID, 50.0, timestamp, "food", false)

		mock.ExpectQuery("SELECT (.+) FROM transactions WHERE user_id = (.+)").
			WithArgs(userID).
			WillReturnRows(rows)

		transactions, err := repo.GetByUserID(ctx, userID)

		assert.NoError(t, err)
		assert.Len(t, transactions, 1)
		assert.Equal(t, transactionID, transactions[0].ID)
		assert.Equal(t, userID, transactions[0].UserID)
		assert.Equal(t, 50.0, transactions[0].Amount)
		assert.Equal(t, "food", transactions[0].Category)
		assert.False(t, transactions[0].DiscountApplied)
	})

	// Test GetByPeriod
	t.Run("GetByPeriod", func(t *testing.T) {
		startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)

		transactionID := uuid.New()
		userID := uuid.New()
		timestamp := time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC)

		rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "timestamp", "category", "discount_applied"}).
			AddRow(transactionID, userID, 100.0, timestamp, "coffee", true)

		mock.ExpectQuery("SELECT (.+) FROM transactions WHERE timestamp BETWEEN (.+) AND (.+)").
			WithArgs(startTime, endTime).
			WillReturnRows(rows)

		transactions, err := repo.GetByPeriod(ctx, startTime, endTime)

		assert.NoError(t, err)
		assert.Len(t, transactions, 1)
		assert.Equal(t, transactionID, transactions[0].ID)
		assert.Equal(t, userID, transactions[0].UserID)
		assert.Equal(t, 100.0, transactions[0].Amount)
		assert.Equal(t, timestamp, transactions[0].Timestamp)
		assert.Equal(t, "coffee", transactions[0].Category)
		assert.True(t, transactions[0].DiscountApplied)
	})

	// Test Create
	t.Run("Create", func(t *testing.T) {
		transaction := &entities.Transaction{
			ID:              uuid.New(),
			UserID:          uuid.New(),
			Amount:          75.0,
			Timestamp:       time.Now(),
			Category:        "beverage",
			DiscountApplied: true,
		}

		mock.ExpectExec("INSERT INTO transactions").
			WithArgs(transaction.ID, transaction.UserID, transaction.Amount, transaction.Timestamp, transaction.Category, transaction.DiscountApplied).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(ctx, transaction)

		assert.NoError(t, err)
	})
}

func testUserRepositoryInMain(t *testing.T, repo repositories.UserRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()

	// Test Get
	t.Run("Get", func(t *testing.T) {
		userID := uuid.New()
		registrationDate := time.Now().Add(-30 * 24 * time.Hour)
		lastActivity := time.Now().Add(-2 * 24 * time.Hour)

		rows := sqlmock.NewRows([]string{"id", "email", "phone", "age", "gender", "city", "registration_date", "last_activity"}).
			AddRow(userID, "user@example.com", "+1234567890", 30, "male", "New York", registrationDate, lastActivity)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE id = (.+)").
			WithArgs(userID).
			WillReturnRows(rows)

		user, err := repo.Get(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, "user@example.com", user.Email)
		assert.Equal(t, "+1234567890", user.Phone)
		assert.Equal(t, 30, user.Age)
		assert.Equal(t, "male", user.Gender)
		assert.Equal(t, "New York", user.City)
	})

	// Test Create
	t.Run("Create", func(t *testing.T) {
		userID := uuid.New()
		registrationDate := time.Now()
		lastActivity := time.Now()

		user := &entities.User{
			Email:  "newuser@example.com",
			Phone:  "+9876543210",
			Age:    25,
			Gender: "female",
			City:   "San Francisco",
		}

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Email, user.Phone, user.Age, user.Gender, user.City).
			WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "last_activity"}).
				AddRow(userID, registrationDate, lastActivity))

		err := repo.Create(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, userID, user.ID)
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		user := &entities.User{
			ID:     uuid.New(),
			Email:  "updated@example.com",
			Phone:  "+1122334455",
			Age:    35,
			Gender: "male",
			City:   "Chicago",
		}

		mock.ExpectExec("UPDATE users SET").
			WithArgs(user.Email, user.Phone, user.Age, user.Gender, user.City, user.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(ctx, user)

		assert.NoError(t, err)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		userID := uuid.New()

		mock.ExpectExec("DELETE FROM users WHERE id = (.+)").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(ctx, userID)

		assert.NoError(t, err)
	})

	// Test List
	t.Run("List", func(t *testing.T) {
		limit := 10
		offset := 0

		user1ID := uuid.New()
		user2ID := uuid.New()
		registrationDate1 := time.Now().Add(-60 * 24 * time.Hour)
		registrationDate2 := time.Now().Add(-30 * 24 * time.Hour)
		lastActivity1 := time.Now().Add(-5 * 24 * time.Hour)
		lastActivity2 := time.Now().Add(-2 * 24 * time.Hour)

		rows := sqlmock.NewRows([]string{"id", "email", "phone", "age", "gender", "city", "registration_date", "last_activity"}).
			AddRow(user1ID, "user1@example.com", "+1111111111", 28, "female", "Boston", registrationDate1, lastActivity1).
			AddRow(user2ID, "user2@example.com", "+2222222222", 42, "male", "Miami", registrationDate2, lastActivity2)

		mock.ExpectQuery("SELECT (.+) FROM users ORDER BY registration_date DESC LIMIT (.+) OFFSET (.+)").
			WithArgs(limit, offset).
			WillReturnRows(rows)

		users, err := repo.List(ctx, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, user1ID, users[0].ID)
		assert.Equal(t, user2ID, users[1].ID)
	})
}

func testSegmentRepositoryInMain(t *testing.T, repo repositories.SegmentRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()

	// Test Get
	t.Run("Get", func(t *testing.T) {
		segmentID := uuid.New()
		createdAt := time.Now().Add(-7 * 24 * time.Hour)

		centroidData := []byte(`{"recency": 3.5, "frequency": 2.1, "monetary": 150.75}`)

		rows := sqlmock.NewRows([]string{"id", "name", "type", "algorithm", "centroid_data", "created_at"}).
			AddRow(segmentID, "High Value", "RFM", "KMeans", centroidData, createdAt)

		mock.ExpectQuery("SELECT (.+) FROM segments WHERE id = (.+)").
			WithArgs(segmentID).
			WillReturnRows(rows)

		segment, err := repo.Get(ctx, segmentID)

		assert.NoError(t, err)
		assert.NotNil(t, segment)
		assert.Equal(t, segmentID, segment.ID)
		assert.Equal(t, "High Value", segment.Name)
		assert.Equal(t, "RFM", segment.Type)
		assert.Equal(t, "KMeans", segment.Algorithm)
	})

	// Test Create
	t.Run("Create", func(t *testing.T) {
		segmentID := uuid.New()
		createdAt := time.Now()

		segment := &entities.Segment{
			ID:        segmentID,
			Name:      "New Customers",
			Type:      "RFM",
			Algorithm: "KMeans",
			CentroidData: map[string]interface{}{
				"recency":   1.2,
				"frequency": 5.0,
				"monetary":  300.50,
			},
			CreatedAt: createdAt,
		}

		mock.ExpectExec("INSERT INTO segments").
			WithArgs(segment.ID, segment.Name, segment.Type, segment.Algorithm, sqlmock.AnyArg(), segment.CreatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(ctx, segment)

		assert.NoError(t, err)
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		segmentID := uuid.New()

		segment := &entities.Segment{
			ID:        segmentID,
			Name:      "Updated Segment",
			Type:      "behavior",
			Algorithm: "DBSCAN",
			CentroidData: map[string]interface{}{
				"recency":   2.5,
				"frequency": 3.7,
				"monetary":  220.30,
			},
		}

		mock.ExpectExec("UPDATE segments").
			WithArgs(segment.Name, segment.Type, segment.Algorithm, sqlmock.AnyArg(), segment.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(ctx, segment)

		assert.NoError(t, err)
	})

	// Test GetByType
	t.Run("GetByType", func(t *testing.T) {
		segmentType := "RFM"

		segment1ID := uuid.New()
		segment2ID := uuid.New()
		createdAt1 := time.Now().Add(-14 * 24 * time.Hour)
		createdAt2 := time.Now().Add(-7 * 24 * time.Hour)

		centroidData1 := []byte(`{"recency": 4.2, "frequency": 1.8, "monetary": 120.50}`)
		centroidData2 := []byte(`{"recency": 2.1, "frequency": 3.5, "monetary": 250.75}`)

		rows := sqlmock.NewRows([]string{"id", "name", "type", "algorithm", "centroid_data", "created_at"}).
			AddRow(segment1ID, "Low Value", "RFM", "KMeans", centroidData1, createdAt1).
			AddRow(segment2ID, "High Value", "RFM", "KMeans", centroidData2, createdAt2)

		mock.ExpectQuery("SELECT (.+) FROM segments WHERE type = (.+)").
			WithArgs(segmentType).
			WillReturnRows(rows)

		segments, err := repo.GetByType(ctx, segmentType)

		assert.NoError(t, err)
		assert.Len(t, segments, 2)
		assert.Equal(t, segment1ID, segments[0].ID)
		assert.Equal(t, segment2ID, segments[1].ID)
	})
}

func testUserMetricsRepositoryInMain(t *testing.T, repo repositories.UserMetricsRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()

	// Test Get
	t.Run("Get", func(t *testing.T) {
		userID := uuid.New()
		segmentID := uuid.New()

		rows := sqlmock.NewRows([]string{"user_id", "recency", "frequency", "monetary", "tbp", "avg_check", "last_segment_id"}).
			AddRow(userID, 5, 12, 1500.50, 7.5, 125.04, segmentID)

		mock.ExpectQuery("SELECT (.+) FROM user_metrics WHERE user_id = (.+)").
			WithArgs(userID).
			WillReturnRows(rows)

		metrics, err := repo.Get(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, metrics)
		assert.Equal(t, userID, metrics.UserID)
		assert.Equal(t, 5, metrics.Recency)
		assert.Equal(t, 12, metrics.Frequency)
		assert.Equal(t, 1500.50, metrics.Monetary)
		assert.Equal(t, 7.5, metrics.TBP)
		assert.Equal(t, 125.04, metrics.AvgCheck)
		assert.Equal(t, segmentID, metrics.LastSegmentID)
	})

	// Test Create
	t.Run("Create", func(t *testing.T) {
		userID := uuid.New()
		segmentID := uuid.New()

		metrics := &entities.UserMetrics{
			UserID:        userID,
			Recency:       3,
			Frequency:     8,
			Monetary:      950.25,
			TBP:           5.2,
			AvgCheck:      118.78,
			LastSegmentID: segmentID,
		}

		mock.ExpectExec("INSERT INTO user_metrics").
			WithArgs(
				metrics.UserID,
				metrics.Recency,
				metrics.Frequency,
				metrics.Monetary,
				metrics.TBP,
				metrics.AvgCheck,
				metrics.LastSegmentID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(ctx, metrics)

		assert.NoError(t, err)
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		userID := uuid.New()
		segmentID := uuid.New()

		metrics := &entities.UserMetrics{
			UserID:        userID,
			Recency:       2,
			Frequency:     15,
			Monetary:      2200.75,
			TBP:           4.8,
			AvgCheck:      146.72,
			LastSegmentID: segmentID,
		}

		mock.ExpectExec("INSERT INTO user_metrics (.+) ON CONFLICT").
			WithArgs(
				metrics.UserID,
				metrics.Recency,
				metrics.Frequency,
				metrics.Monetary,
				metrics.TBP,
				metrics.AvgCheck,
				metrics.LastSegmentID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Update(ctx, metrics)

		assert.NoError(t, err)
	})

	// Test CalculateMetrics
	t.Run("CalculateMetrics", func(t *testing.T) {
		userID := uuid.New()
		segmentID := uuid.New()

		// Начало транзакции
		mock.ExpectBegin()

		// Запрос recency
		mock.ExpectQuery("SELECT EXTRACT\\(DAY FROM NOW\\(\\) - MAX\\(timestamp\\)\\)::int").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"recency"}).AddRow(4))

		// Запрос frequency
		mock.ExpectQuery("SELECT COUNT\\(\\*\\)").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"frequency"}).AddRow(10))

		// Запрос monetary
		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(amount\\), 0\\)").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"monetary"}).AddRow(1250.50))

		// Запрос tbp
		mock.ExpectQuery("WITH purchase_dates AS").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"tbp"}).AddRow(6.5))

		// Запрос avgCheck
		mock.ExpectQuery("SELECT COALESCE\\(AVG\\(amount\\), 0\\)").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"avg_check"}).AddRow(125.05))

		// Запрос lastSegmentID
		mock.ExpectQuery("SELECT COALESCE\\(last_segment_id, '00000000-0000-0000-0000-000000000000'::uuid\\)").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"last_segment_id"}).AddRow(segmentID))

		// Обновление метрик
		mock.ExpectExec("INSERT INTO user_metrics (.+) ON CONFLICT").
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Коммит транзакции
		mock.ExpectCommit()

		metrics, err := repo.CalculateMetrics(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, metrics)
		assert.Equal(t, userID, metrics.UserID)
		assert.Equal(t, 4, metrics.Recency)
		assert.Equal(t, 10, metrics.Frequency)
		assert.Equal(t, 1250.50, metrics.Monetary)
		assert.Equal(t, 6.5, metrics.TBP)
		assert.Equal(t, 125.05, metrics.AvgCheck)
		assert.Equal(t, segmentID, metrics.LastSegmentID)
	})
}
