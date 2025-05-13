// internal/interfaces/kafka/config.go
package kafka

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config представляет конфигурацию для Kafka
type Config struct {
	Brokers   []string
	Topic     string
	GroupID   string
	Timeout   time.Duration
	MaxRetry  int
	BatchSize int
}

// NewConfigFromEnv создает конфигурацию Kafka из переменных окружения
func NewConfigFromEnv() *Config {
	// Получаем адреса брокеров
	brokersList := os.Getenv("KAFKA_BROKERS")
	if brokersList == "" {
		brokersList = "kafka:9092" // Значение по умолчанию
	}
	brokers := strings.Split(brokersList, ",")

	// Получаем имя топика
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "user-events" // Значение по умолчанию
	}

	// Получаем идентификатор группы
	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "user-service-group" // Значение по умолчанию
	}

	// Получаем таймаут
	timeoutStr := os.Getenv("KAFKA_TIMEOUT")
	timeout := 10 * time.Second // Значение по умолчанию
	if timeoutStr != "" {
		if parsedTimeout, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = parsedTimeout
		}
	}

	// Получаем максимальное количество повторных попыток
	maxRetry := 5 // Значение по умолчанию
	if maxRetryStr := os.Getenv("KAFKA_MAX_RETRY"); maxRetryStr != "" {
		if _, err := fmt.Sscanf(maxRetryStr, "%d", &maxRetry); err != nil {
			// Если не удалось распарсить, используем значение по умолчанию
		}
	}

	// Получаем размер пакета
	batchSize := 1 // Значение по умолчанию
	if batchSizeStr := os.Getenv("KAFKA_BATCH_SIZE"); batchSizeStr != "" {
		if _, err := fmt.Sscanf(batchSizeStr, "%d", &batchSize); err != nil {
			// Если не удалось распарсить, используем значение по умолчанию
		}
	}

	return &Config{
		Brokers:   brokers,
		Topic:     topic,
		GroupID:   groupID,
		Timeout:   timeout,
		MaxRetry:  maxRetry,
		BatchSize: batchSize,
	}
}