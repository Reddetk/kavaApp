// test/transaction_repository_helpers.go
package test

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"
	"user-service/internal/infrastructure/postgres"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestTransactionRepositoryHelpers содержит общие функции для тестирования TransactionRepository
// Эти функции могут быть использованы в разных тестовых файлах

// SetupTransactionRepositoryTest создает мок базы данных и репозиторий для тестирования
func SetupTransactionRepositoryTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repositories.TransactionRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo := postgres.NewTransactionRepository(db)
	return db, mock, repo
}

// TestGetByUserIDHelper тестирует метод GetByUserID
func TestGetByUserIDHelper(t *testing.T, repo repositories.TransactionRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
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
}

// TestGetByPeriodHelper тестирует метод GetByPeriod
func TestGetByPeriodHelper(t *testing.T, repo repositories.TransactionRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
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
}

// TestCreateHelper тестирует метод Create
func TestCreateHelper(t *testing.T, repo repositories.TransactionRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
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
}
