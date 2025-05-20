package repositories

import (
	"context"
	"time"

	"analitics-service/internal/domain/enteties"
)

// TransactionRepository определяет интерфейс для работы с транзакциями
type TransactionRepository interface {
	// GetTransactionsByPeriod возвращает транзакции за указанный период
	GetTransactionsByPeriod(ctx context.Context, startDate, endDate time.Time) ([]entities.Transaction, error)

	// GetTransactionByID возвращает транзакцию по её ID
	GetTransactionByID(ctx context.Context, transactionID string) (entities.Transaction, error)

	// GetTransactionsByCustomerID возвращает транзакции конкретного клиента
	GetTransactionsByCustomerID(ctx context.Context, customerID string, startDate, endDate time.Time) ([]entities.Transaction, error)

	// CreateTransaction создает новую транзакцию
	CreateTransaction(ctx context.Context, transaction entities.Transaction) error

	// GetTransactionsWithProduct возвращает транзакции, содержащие указанный продукт
	GetTransactionsWithProduct(ctx context.Context, productID string, startDate, endDate time.Time) ([]entities.Transaction, error)

	// GetTransactionCount возвращает количество транзакций за период
	GetTransactionCount(ctx context.Context, startDate, endDate time.Time) (int, error)
}
