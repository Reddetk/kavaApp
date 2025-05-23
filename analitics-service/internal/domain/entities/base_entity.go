// internal/domain/entities/base_entity.go
package entities

import "time"

// BaseEntity содержит общие поля для всех сущностей
type BaseEntity struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
