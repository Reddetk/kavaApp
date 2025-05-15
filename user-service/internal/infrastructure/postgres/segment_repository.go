// user-service/internal/infrastructure/postgres/segment_repository.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"

	"github.com/google/uuid"
)

type SegmentRepository struct {
	db *sql.DB
}

func NewSegmentRepository(db *sql.DB) repositories.SegmentRepository {
	return &SegmentRepository{db: db}
}

func (r *SegmentRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *SegmentRepository) Create(ctx context.Context, s *entities.Segment) error {
	centroidData, err := json.Marshal(s.CentroidData)
	if err != nil {
		return err
	}

	query := `INSERT INTO public.user_segments (id, name, type, algorithm, centroid_data, created_at)
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = r.db.ExecContext(ctx, query, s.ID, s.Name, s.Type, s.Algorithm, centroidData, s.CreatedAt)
	return err
}

func (r *SegmentRepository) Get(ctx context.Context, id uuid.UUID) (*entities.Segment, error) {
	query := `SELECT id, name, type, algorithm, centroid_data, created_at
              FROM public.user_segments 
              WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var s entities.Segment
	var centroidData []byte
	if err := row.Scan(&s.ID, &s.Name, &s.Type, &s.Algorithm, &centroidData, &s.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(centroidData, &s.CentroidData); err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *SegmentRepository) Update(ctx context.Context, s *entities.Segment) error {
	centroidData, err := json.Marshal(s.CentroidData)
	if err != nil {
		return err
	}

	query := `UPDATE public.user_segments 
              SET name = $1, type = $2, algorithm = $3, centroid_data = $4
              WHERE id = $5`
	_, err = r.db.ExecContext(ctx, query, s.Name, s.Type, s.Algorithm, centroidData, s.ID)
	return err
}

func (r *SegmentRepository) GetByType(ctx context.Context, segmentType string) ([]*entities.Segment, error) {
	query := `SELECT id, name, type, algorithm, centroid_data, created_at
              FROM public.user_segments 
              WHERE type = $1
              ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, segmentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segments []*entities.Segment
	for rows.Next() {
		var s entities.Segment
		var centroidData []byte
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.Algorithm, &centroidData, &s.CreatedAt); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(centroidData, &s.CentroidData); err != nil {
			return nil, err
		}

		segments = append(segments, &s)
	}
	return segments, nil
}