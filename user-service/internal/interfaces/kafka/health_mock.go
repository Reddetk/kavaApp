// internal/interfaces/kafka/health_mock.go

//go:build !kafka
// +build !kafka

package kafka

import (
	"log"
	"time"
)

// HealthCheck имитирует проверку доступности Kafka
func HealthCheck(brokers []string, timeout time.Duration) error {
	log.Println("Mock Kafka health check passed")
	return nil
}