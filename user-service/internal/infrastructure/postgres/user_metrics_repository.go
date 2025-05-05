// internal/infrastructure/postgres/user_metrics_repository.go
package postgres

import (
	"context"
	"database/sql"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"

	"github.com/google/uuid"
)

type UserMetricsRepository struct {
	db *sql.DB
}

func NewUserMetricsRepository(db *sql.DB) repositories.UserMetricsRepository {
	return &UserMetricsRepository{db: db}
}

// Get retrieves user metrics from the database for a given user ID
func (r *UserMetricsRepository) Get(ctx context.Context, userID uuid.UUID) (*entities.UserMetrics, error) {
	// SQL query to select all metrics fields for a specific user
	query := `
        SELECT user_id, recency, frequency, monetary, tbp, avg_check, last_segment_id
        FROM user_metrics
        WHERE user_id = $1`

	// Initialize metrics struct to store the result
	var metrics entities.UserMetrics
	// Execute query and scan results into metrics struct
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&metrics.UserID,
		&metrics.Recency,
		&metrics.Frequency,
		&metrics.Monetary,
		&metrics.TBP,
		&metrics.AvgCheck,
		&metrics.LastSegmentID,
	)

	// Return nil if no metrics found for user
	if err == sql.ErrNoRows {
		return nil, nil
	}
	// Return error if query failed
	if err != nil {
		return nil, err
	}

	// Return pointer to metrics struct
	return &metrics, nil
}

func (r *UserMetricsRepository) Create(ctx context.Context, metrics *entities.UserMetrics) error {
	query := `
        INSERT INTO user_metrics (
            user_id, recency, frequency, monetary, tbp, avg_check, last_segment_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		metrics.UserID,
		metrics.Recency,
		metrics.Frequency,
		metrics.Monetary,
		metrics.TBP,
		metrics.AvgCheck,
		metrics.LastSegmentID,
	)

	return err
}

func (r *UserMetricsRepository) Update(ctx context.Context, metrics *entities.UserMetrics) error {
	query := `
        INSERT INTO user_metrics (
            user_id, recency, frequency, monetary, tbp, avg_check, last_segment_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (user_id) DO UPDATE SET
            recency = EXCLUDED.recency,
            frequency = EXCLUDED.frequency,
            monetary = EXCLUDED.monetary,
            tbp = EXCLUDED.tbp,
            avg_check = EXCLUDED.avg_check,
            last_segment_id = EXCLUDED.last_segment_id`

	_, err := r.db.ExecContext(ctx, query,
		metrics.UserID,
		metrics.Recency,
		metrics.Frequency,
		metrics.Monetary,
		metrics.TBP,
		metrics.AvgCheck,
		metrics.LastSegmentID,
	)

	return err
}

func (r *UserMetricsRepository) CalculateMetrics(ctx context.Context, userID uuid.UUID) (*entities.UserMetrics, error) {
	// Start a transaction to ensure consistency
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Calculate recency (days since last purchase)
	recencyQuery := `
        SELECT EXTRACT(DAY FROM NOW() - MAX(timestamp))::int
        FROM transactions
        WHERE user_id = $1`

	var recency int
	err = tx.QueryRowContext(ctx, recencyQuery, userID).Scan(&recency)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Calculate frequency (total number of purchases)
	frequencyQuery := `
        SELECT COUNT(*)
        FROM transactions
        WHERE user_id = $1`

	var frequency int
	err = tx.QueryRowContext(ctx, frequencyQuery, userID).Scan(&frequency)
	if err != nil {
		return nil, err
	}

	// Calculate monetary (total amount spent)
	monetaryQuery := `
        SELECT COALESCE(SUM(amount), 0)
        FROM transactions
        WHERE user_id = $1`

	var monetary float64
	err = tx.QueryRowContext(ctx, monetaryQuery, userID).Scan(&monetary)
	if err != nil {
		return nil, err
	}

	// Calculate average time between purchases (TBP)
	tbpQuery := `
        WITH purchase_dates AS (
            SELECT timestamp,
                   LAG(timestamp) OVER (ORDER BY timestamp) as prev_timestamp
            FROM transactions
            WHERE user_id = $1
            ORDER BY timestamp
        )
        SELECT COALESCE(
            EXTRACT(EPOCH FROM AVG(timestamp - prev_timestamp))/86400, 
            0
        )
        FROM purchase_dates
        WHERE prev_timestamp IS NOT NULL`

	var tbp float64
	err = tx.QueryRowContext(ctx, tbpQuery, userID).Scan(&tbp)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Calculate average check
	avgCheckQuery := `
        SELECT COALESCE(AVG(amount), 0)
        FROM transactions
        WHERE user_id = $1`

	var avgCheck float64
	err = tx.QueryRowContext(ctx, avgCheckQuery, userID).Scan(&avgCheck)
	if err != nil {
		return nil, err
	}

	// Get last segment ID
	lastSegmentQuery := `
        SELECT COALESCE(last_segment_id, '00000000-0000-0000-0000-000000000000'::uuid)
        FROM user_metrics
        WHERE user_id = $1`

	var lastSegmentID uuid.UUID
	err = tx.QueryRowContext(ctx, lastSegmentQuery, userID).Scan(&lastSegmentID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	metrics := &entities.UserMetrics{
		UserID:        userID,
		Recency:       recency,
		Frequency:     frequency,
		Monetary:      monetary,
		TBP:           tbp,
		AvgCheck:      avgCheck,
		LastSegmentID: lastSegmentID,
	}

	// Update the metrics in the database
	err = r.Update(ctx, metrics)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
