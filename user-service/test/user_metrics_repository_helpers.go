// test/user_metrics_repository_helpers.go
package test

import (
	"context"
	"database/sql"
	"testing"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"
	"user-service/internal/infrastructure/postgres"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// SetupUserMetricsRepositoryTest создает мок базы данных и репозиторий для тестирования
func SetupUserMetricsRepositoryTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repositories.UserMetricsRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo := postgres.NewUserMetricsRepository(db)
	return db, mock, repo
}

// TestUserMetricsGetHelper тестирует метод Get
func TestUserMetricsGetHelper(t *testing.T, repo repositories.UserMetricsRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
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
}

// TestUserMetricsCreateHelper тестирует метод Create
func TestUserMetricsCreateHelper(t *testing.T, repo repositories.UserMetricsRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
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
}

// TestUserMetricsUpdateHelper тестирует метод Update
func TestUserMetricsUpdateHelper(t *testing.T, repo repositories.UserMetricsRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
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
}

// TestUserMetricsCalculateMetricsHelper тестирует метод CalculateMetrics
func TestUserMetricsCalculateMetricsHelper(t *testing.T, repo repositories.UserMetricsRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
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
}
