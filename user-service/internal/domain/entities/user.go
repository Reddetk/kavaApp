// internal/domain/entities/user.go
package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID
	Email            string
	Phone            string
	Age              int
	Gender           string
	City             string
	RegistrationDate time.Time
	LastActivity     time.Time
}
