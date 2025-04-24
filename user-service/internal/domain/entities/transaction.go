// internal/domain/entities/transaction.go
package entities

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Amount          float64
	Timestamp       time.Time
	Category        string
	DiscountApplied bool
}
