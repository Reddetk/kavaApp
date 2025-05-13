// internal/interfaces/kafka/producer_interface.go
package kafka

import (
	"context"
)

// MessageProducer интерфейс для отправки сообщений
type MessageProducer interface {
	SendMessage(ctx context.Context, key string, value interface{}) error
	Close() error
}