// internal/interfaces/kafka/producer_mock.go

//go:build !kafka
// +build !kafka

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// MockProducer представляет мок-реализацию Kafka продюсера
type MockProducer struct {
	messages []MockMessage
	topic    string
}

// MockMessage представляет сообщение, сохраненное мок-продюсером
type MockMessage struct {
	Key   string
	Value []byte
}

// NewProducer создает новый экземпляр MockProducer
func NewProducer(brokers []string, topic string) MessageProducer {
	log.Printf("Creating mock Kafka producer for topic: %s", topic)
	return &MockProducer{
		messages: make([]MockMessage, 0),
		topic:    topic,
	}
}

// SendMessage имитирует отправку сообщения в Kafka
func (p *MockProducer) SendMessage(ctx context.Context, key string, value interface{}) error {
	// Сериализуем значение в JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message value: %w", err)
	}

	// Сохраняем сообщение в памяти
	p.messages = append(p.messages, MockMessage{
		Key:   key,
		Value: jsonValue,
	})

	log.Printf("Mock message sent to topic %s: key=%s", p.topic, key)
	return nil
}

// Close имитирует закрытие соединения с Kafka
func (p *MockProducer) Close() error {
	log.Println("Mock Kafka producer closed")
	return nil
}

// GetMessages возвращает все сохраненные сообщения
func (p *MockProducer) GetMessages() []MockMessage {
	return p.messages
}