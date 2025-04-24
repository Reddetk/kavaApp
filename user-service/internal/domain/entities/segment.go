// internal/domain/entities/segment.go
package entities

import (
	"time"

	"github.com/google/uuid"
)

type Segment struct {
	ID           uuid.UUID
	Name         string
	Type         string // "RFM", "behavior", "demographic", "promo"
	Algorithm    string // "KMeans", "DBSCAN"
	CentroidData map[string]interface{}
	CreatedAt    time.Time
}
