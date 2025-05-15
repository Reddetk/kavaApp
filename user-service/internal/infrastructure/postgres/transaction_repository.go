// internal/infrastructure/postgres/transaction_repository.go
package postgres

import (
	"context"
	"database/sql"
	"time"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"

	"github.com/google/uuid"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repositories.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// GetByPeriod implements repositories.TransactionRepository.
func (r *TransactionRepository) GetByPeriod(ctx context.Context, start time.Time, end time.Time) ([]*entities.Transaction, error) {
	query := `SELECT id, user_id, amount, timestamp, category, discount_applied 
			  FROM public.transactions 
			  WHERE timestamp BETWEEN $1 AND $2
			  ORDER BY timestamp DESC`

	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		var t entities.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Timestamp, &t.Category, &t.DiscountApplied); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}

// GetByUserID implements repositories.TransactionRepository.
func (r *TransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Transaction, error) {
	query := `SELECT id, user_id, amount, timestamp, category, discount_applied 
			  FROM public.transactions 
			  WHERE user_id = $1
			  ORDER BY timestamp DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		var t entities.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Timestamp, &t.Category, &t.DiscountApplied); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}

func (r *TransactionRepository) Create(ctx context.Context, t *entities.Transaction) error {
	query := `INSERT INTO public.transactions (id, user_id, amount, timestamp, category, discount_applied)
				VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, t.ID, t.UserID, t.Amount, t.Timestamp, t.Category, t.DiscountApplied)
	return err
}

func (r *TransactionRepository) Get(ctx context.Context, id uuid.UUID) (*entities.Transaction, error) {
	query := `SELECT id, user_id, amount, timestamp, category, discount_applied FROM public.transactions WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var t entities.Transaction
	if err := row.Scan(&t.ID, &t.UserID, &t.Amount, &t.Timestamp, &t.Category, &t.DiscountApplied); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepository) Update(ctx context.Context, t *entities.Transaction) error {
	query := `UPDATE public.transactions SET user_id = $1, amount = $2, timestamp = $3, category = $4, discount_applied = $5 WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, t.UserID, t.Amount, t.Timestamp, t.Category, t.DiscountApplied, t.ID)
	return err
}

func (r *TransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM public.transactions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *TransactionRepository) List(ctx context.Context, limit, offset int) ([]*entities.Transaction, error) {
	query := `SELECT id, user_id, amount, timestamp, category, discount_applied FROM public.transactions ORDER BY timestamp DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		var t entities.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Timestamp, &t.Category, &t.DiscountApplied); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}
