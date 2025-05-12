// test/segment_repository_helpers.go
package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"
	"user-service/internal/infrastructure/postgres"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// SetupSegmentRepositoryTest создает мок базы данных и репозиторий для тестирования
func SetupSegmentRepositoryTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, repositories.SegmentRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo := postgres.NewSegmentRepository(db)
	return db, mock, repo
}

// TestSegmentGetHelper тестирует метод Get
func TestSegmentGetHelper(t *testing.T, repo repositories.SegmentRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
	segmentID := uuid.New()
	createdAt := time.Now().Add(-7 * 24 * time.Hour)

	centroidData := map[string]interface{}{
		"recency":   3.5,
		"frequency": 2.1,
		"monetary":  150.75,
	}

	centroidDataJSON, _ := json.Marshal(centroidData)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "algorithm", "centroid_data", "created_at"}).
		AddRow(segmentID, "High Value", "RFM", "KMeans", centroidDataJSON, createdAt)

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
	assert.Equal(t, createdAt.Unix(), segment.CreatedAt.Unix())
	assert.Equal(t, 3.5, segment.CentroidData["recency"])
	assert.Equal(t, 2.1, segment.CentroidData["frequency"])
	assert.Equal(t, 150.75, segment.CentroidData["monetary"])
}

// TestSegmentCreateHelper тестирует метод Create
func TestSegmentCreateHelper(t *testing.T, repo repositories.SegmentRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
	segmentID := uuid.New()
	createdAt := time.Now()

	centroidData := map[string]interface{}{
		"recency":   1.2,
		"frequency": 5.0,
		"monetary":  300.50,
	}

	segment := &entities.Segment{
		ID:           segmentID,
		Name:         "New Customers",
		Type:         "RFM",
		Algorithm:    "KMeans",
		CentroidData: centroidData,
		CreatedAt:    createdAt,
	}

	centroidDataJSON, _ := json.Marshal(centroidData)

	mock.ExpectExec("INSERT INTO segments").
		WithArgs(segment.ID, segment.Name, segment.Type, segment.Algorithm, centroidDataJSON, segment.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(ctx, segment)

	assert.NoError(t, err)
}

// TestSegmentUpdateHelper тестирует метод Update
func TestSegmentUpdateHelper(t *testing.T, repo repositories.SegmentRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
	segmentID := uuid.New()

	centroidData := map[string]interface{}{
		"recency":   2.5,
		"frequency": 3.7,
		"monetary":  220.30,
	}

	segment := &entities.Segment{
		ID:           segmentID,
		Name:         "Updated Segment",
		Type:         "behavior",
		Algorithm:    "DBSCAN",
		CentroidData: centroidData,
	}

	centroidDataJSON, _ := json.Marshal(centroidData)

	mock.ExpectExec("UPDATE segments").
		WithArgs(segment.Name, segment.Type, segment.Algorithm, centroidDataJSON, segment.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(ctx, segment)

	assert.NoError(t, err)
}

// TestSegmentGetByTypeHelper тестирует метод GetByType
func TestSegmentGetByTypeHelper(t *testing.T, repo repositories.SegmentRepository, mock sqlmock.Sqlmock) {
	ctx := context.Background()
	segmentType := "RFM"

	segment1ID := uuid.New()
	segment2ID := uuid.New()
	createdAt1 := time.Now().Add(-14 * 24 * time.Hour)
	createdAt2 := time.Now().Add(-7 * 24 * time.Hour)

	centroidData1 := map[string]interface{}{
		"recency":   4.2,
		"frequency": 1.8,
		"monetary":  120.50,
	}

	centroidData2 := map[string]interface{}{
		"recency":   2.1,
		"frequency": 3.5,
		"monetary":  250.75,
	}

	centroidData1JSON, _ := json.Marshal(centroidData1)
	centroidData2JSON, _ := json.Marshal(centroidData2)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "algorithm", "centroid_data", "created_at"}).
		AddRow(segment1ID, "Low Value", "RFM", "KMeans", centroidData1JSON, createdAt1).
		AddRow(segment2ID, "High Value", "RFM", "KMeans", centroidData2JSON, createdAt2)

	mock.ExpectQuery("SELECT (.+) FROM segments WHERE type = (.+)").
		WithArgs(segmentType).
		WillReturnRows(rows)

	segments, err := repo.GetByType(ctx, segmentType)

	assert.NoError(t, err)
	assert.Len(t, segments, 2)

	// First segment
	assert.Equal(t, segment1ID, segments[0].ID)
	assert.Equal(t, "Low Value", segments[0].Name)
	assert.Equal(t, "RFM", segments[0].Type)
	assert.Equal(t, "KMeans", segments[0].Algorithm)
	assert.Equal(t, createdAt1.Unix(), segments[0].CreatedAt.Unix())

	// Second segment
	assert.Equal(t, segment2ID, segments[1].ID)
	assert.Equal(t, "High Value", segments[1].Name)
	assert.Equal(t, "RFM", segments[1].Type)
	assert.Equal(t, "KMeans", segments[1].Algorithm)
	assert.Equal(t, createdAt2.Unix(), segments[1].CreatedAt.Unix())
}
