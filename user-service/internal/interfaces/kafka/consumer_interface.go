// internal/interfaces/kafka/consumer_interface.go
package kafka

import (
	"context"
)

// MessageConsumer интерфейс для потребления сообщений из Kafka
type MessageConsumer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}