// internal/interfaces/kafka/health.go

//go:build kafka
// +build kafka

package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// HealthCheck проверяет доступность Kafka
func HealthCheck(brokers []string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dialer := &kafka.Dialer{
		Timeout:   timeout,
		DualStack: true,
	}

	// Пытаемся подключиться к каждому брокеру
	for _, broker := range brokers {
		conn, err := dialer.DialContext(ctx, "tcp", broker)
		if err != nil {
			return fmt.Errorf("failed to connect to Kafka broker %s: %w", broker, err)
		}
		conn.Close()
	}

	log.Println("Kafka health check passed")
	return nil
}