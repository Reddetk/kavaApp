// internal/interfaces/kafka/health.go
package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// HealthChecker проверяет доступность Kafka
type HealthChecker struct {
	brokers []string
	timeout time.Duration
}

// NewHealthChecker создает новый экземпляр HealthChecker
func NewHealthChecker(brokers []string, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		brokers: brokers,
		timeout: timeout,
	}
}

// Check проверяет доступность Kafka
func (h *HealthChecker) Check(ctx context.Context) error {
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// Создаем диалер с таймаутом
	dialer := &kafka.Dialer{
		Timeout:   h.timeout,
		DualStack: true,
	}

	// Проверяем подключение к каждому брокеру
	for _, broker := range h.brokers {
		conn, err := dialer.DialContext(ctx, "tcp", broker)
		if err != nil {
			return fmt.Errorf("failed to connect to Kafka broker %s: %w", broker, err)
		}
		defer conn.Close()

		// Проверяем, что соединение работает
		if _, err := conn.ApiVersions(); err != nil {
			return fmt.Errorf("failed to get API versions from Kafka broker %s: %w", broker, err)
		}

		log.Printf("Successfully connected to Kafka broker %s", broker)
	}

	return nil
}

// WaitForKafka ожидает доступности Kafka с повторными попытками
func WaitForKafka(ctx context.Context, brokers []string, maxRetries int, retryInterval time.Duration) error {
	checker := NewHealthChecker(brokers, 5*time.Second)

	for i := 0; i < maxRetries; i++ {
		err := checker.Check(ctx)
		if err == nil {
			log.Println("Successfully connected to Kafka")
			return nil
		}

		log.Printf("Failed to connect to Kafka (attempt %d/%d): %v", i+1, maxRetries, err)

		// Проверяем, не был ли отменен контекст
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryInterval):
			// Продолжаем после задержки
		}
	}

	return fmt.Errorf("failed to connect to Kafka after %d attempts", maxRetries)
}