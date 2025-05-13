// internal/interfaces/kafka/producer.go

//go:build kafka
// +build kafka

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer представляет Kafka продюсера
type Producer struct {
	writer *kafka.Writer
}

// NewProducer создает новый экземпляр Producer
func NewProducer(brokers []string, topic string) MessageProducer {
	// Настройка Kafka writer
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		// Настройки для повышения надежности
		Async:       false,
		BatchSize:   1,
		ReadTimeout: 10 * time.Second,
		// Добавляем диалер с таймаутом для предотвращения бесконечного ожидания
		Transport: &kafka.Transport{
			Dial: (&kafka.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true, // Поддержка IPv4 и IPv6
			}).DialFunc,
		},
	}

	return &Producer{
		writer: writer,
	}
}

// SendMessage отправляет сообщение в Kafka
func (p *Producer) SendMessage(ctx context.Context, key string, value interface{}) error {
	// Сериализуем значение в JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message value: %w", err)
	}

	// Создаем сообщение
	message := kafka.Message{
		Key:   []byte(key),
		Value: jsonValue,
		Time:  time.Now(),
	}

	// Отправляем сообщение
	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	log.Printf("Message sent to Kafka: key=%s", key)
	return nil
}

// Close закрывает соединение с Kafka
func (p *Producer) Close() error {
	return p.writer.Close()
}
