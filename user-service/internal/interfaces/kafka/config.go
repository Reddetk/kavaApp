// internal/interfaces/kafka/config.go
package kafka

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Config содержит настройки для Kafka
type Config struct {
	Brokers   []string
	Topic     string
	GroupID   string
	Timeout   time.Duration
	MaxRetry  int
	BatchSize int
}

// LoadConfig загружает конфигурацию Kafka из переменных окружения
func LoadConfig() *Config {
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		log.Println("KAFKA_BROKERS not set, using default: localhost:9092")
		brokersStr = "localhost:9092"
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		log.Println("KAFKA_TOPIC not set, using default: user-events")
		topic = "user-events"
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		log.Println("KAFKA_GROUP_ID not set, using default: user-service-group")
		groupID = "user-service-group"
	}

	timeoutStr := os.Getenv("KAFKA_TIMEOUT")
	timeout := 10 * time.Second
	if timeoutStr != "" {
		var err error
		timeout, err = time.ParseDuration(timeoutStr)
		if err != nil {
			log.Printf("Invalid KAFKA_TIMEOUT: %v, using default: 10s", err)
			timeout = 10 * time.Second
		}
	}

	maxRetryStr := os.Getenv("KAFKA_MAX_RETRY")
	maxRetry := 5
	if maxRetryStr != "" {
		_, err := fmt.Sscanf(maxRetryStr, "%d", &maxRetry)
		if err != nil {
			log.Printf("Invalid KAFKA_MAX_RETRY: %v, using default: 5", err)
			maxRetry = 5
		}
	}

	batchSizeStr := os.Getenv("KAFKA_BATCH_SIZE")
	batchSize := 1
	if batchSizeStr != "" {
		_, err := fmt.Sscanf(batchSizeStr, "%d", &batchSize)
		if err != nil {
			log.Printf("Invalid KAFKA_BATCH_SIZE: %v, using default: 1", err)
			batchSize = 1
		}
	}

	return &Config{
		Brokers:   strings.Split(brokersStr, ","),
		Topic:     topic,
		GroupID:   groupID,
		Timeout:   timeout,
		MaxRetry:  maxRetry,
		BatchSize: batchSize,
	}
}