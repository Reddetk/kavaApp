// internal/interfaces/kafka/consumer.go
package kafka

import (
	"context"
	"user-service/internal/application"

	"github.com/segmentio/kafka-go"
)

type TransactionConsumer struct {
	reader         *kafka.Reader
	userService    *application.UserService
	metricsService *application.RetentionService
}

func NewTransactionConsumer(brokers []string, topic string, us *application.UserService,
	ms *application.RetentionService) *TransactionConsumer {
	return &TransactionConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
		userService:    us,
		metricsService: ms,
	}
}

func (c *TransactionConsumer) Start(ctx context.Context) error {
	//implementation
	return nil
}
