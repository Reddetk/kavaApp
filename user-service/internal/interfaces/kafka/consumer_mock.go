// internal/interfaces/kafka/consumer_mock.go

//go:build !kafka
// +build !kafka

package kafka

import (
	"context"
	"log"
	"time"
	"user-service/internal/application"
)

// MockTransactionConsumer представляет мок-реализацию Kafka консьюмера
type MockTransactionConsumer struct {
	userService    *application.UserService
	metricsService *application.RetentionService
	topic          string
	groupID        string
	running        bool
	stopCh         chan struct{}
}

// NewTransactionConsumer создает новый экземпляр MockTransactionConsumer
func NewTransactionConsumer(brokers []string, topic string, groupID string, us *application.UserService,
	ms *application.RetentionService) MessageConsumer {

	log.Printf("Creating mock Kafka consumer for topic: %s, group: %s", topic, groupID)
	return &MockTransactionConsumer{
		userService:    us,
		metricsService: ms,
		topic:          topic,
		groupID:        groupID,
		running:        false,
		stopCh:         make(chan struct{}),
	}
}

// Start имитирует запуск обработки сообщений из Kafka
func (c *MockTransactionConsumer) Start(ctx context.Context) error {
	if c.running {
		return nil
	}

	c.running = true
	log.Printf("Mock Kafka consumer started for topic: %s, group: %s", c.topic, c.groupID)

	// Создаем контекст, который можно отменить
	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Горутина для обработки сигнала остановки
	go func() {
		select {
		case <-c.stopCh:
			cancel()
			log.Println("Mock consumer received stop signal")
		case <-ctx.Done():
			log.Println("Mock consumer context canceled")
		}
	}()

	// Имитируем работу консьюмера, просто ожидая завершения контекста
	<-consumerCtx.Done()
	c.running = false
	log.Println("Mock Kafka consumer stopped")

	return nil
}

// Stop останавливает имитацию обработки сообщений
func (c *MockTransactionConsumer) Stop(ctx context.Context) error {
	if !c.running {
		return nil
	}

	log.Println("Stopping mock Kafka consumer")
	close(c.stopCh)

	// Ждем некоторое время для имитации завершения
	select {
	case <-time.After(100 * time.Millisecond):
		log.Println("Mock consumer stopped")
	case <-ctx.Done():
		log.Println("Context canceled while stopping mock consumer")
	}

	c.running = false
	return nil
}