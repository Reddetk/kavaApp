// internal/interfaces/kafka/consumer.go
package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"user-service/internal/application"
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// TransactionMessage представляет структуру сообщения о транзакции из Kafka
type TransactionMessage struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Amount          float64   `json:"amount"`
	Timestamp       time.Time `json:"timestamp"`
	Category        string    `json:"category"`
	DiscountApplied bool      `json:"discount_applied"`
}

// TransactionConsumer обрабатывает сообщения о транзакциях из Kafka
type TransactionConsumer struct {
	reader         *kafka.Reader
	userService    *application.UserService
	metricsService *application.RetentionService
	running        bool
	stopCh         chan struct{}
}

// NewTransactionConsumer создает новый экземпляр TransactionConsumer
func NewTransactionConsumer(brokers []string, topic string, us *application.UserService,
	ms *application.RetentionService) *TransactionConsumer {

	// Настройка Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        "user-service",
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        1 * time.Second,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 1 * time.Second,
		ReadBackoffMin: 100 * time.Millisecond,
		ReadBackoffMax: 1 * time.Second,
	})

	return &TransactionConsumer{
		reader:         reader,
		userService:    us,
		metricsService: ms,
		running:        false,
		stopCh:         make(chan struct{}),
	}
}

// Start запускает обработку сообщений из Kafka
func (c *TransactionConsumer) Start(ctx context.Context) error {
	if c.running {
		return errors.New("consumer is already running")
	}

	c.running = true
	log.Println("Starting Kafka consumer")

	// Создаем контекст, который можно отменить
	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Горутина для обработки сигнала остановки
	go func() {
		select {
		case <-c.stopCh:
			cancel()
		case <-ctx.Done():
			// Контекст был отменен извне
		}
	}()

	// Основной цикл обработки сообщений
	for c.running {
		message, err := c.reader.FetchMessage(consumerCtx)
		if err != nil {
			// Проверяем, был ли контекст отменен
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				log.Println("Context canceled, stopping consumer")
				c.running = false
				break
			}

			log.Printf("Error fetching message: %v", err)
			time.Sleep(1 * time.Second) // Небольшая задержка перед повторной попыткой
			continue
		}

		// Обработка сообщения
		if err := c.processMessage(consumerCtx, message); err != nil {
			log.Printf("Error processing message: %v", err)
			// Можно решить, коммитить ли сообщение с ошибкой или нет
			// В данном случае мы коммитим, чтобы не застрять на проблемном сообщении
		}

		// Подтверждае�� обработку сообщения
		if err := c.reader.CommitMessages(consumerCtx, message); err != nil {
			log.Printf("Error committing message: %v", err)
		}
	}

	return nil
}

// Stop останавливает обработку сообщений
func (c *TransactionConsumer) Stop(ctx context.Context) error {
	if !c.running {
		return nil
	}

	log.Println("Stopping Kafka consumer")
	close(c.stopCh)

	// Ждем некоторое время для завершения текущих операций
	select {
	case <-time.After(5 * time.Second):
		log.Println("Forcing consumer to stop")
	case <-ctx.Done():
		log.Println("Context canceled while stopping consumer")
	}

	c.running = false
	return c.reader.Close()
}

// processMessage обрабатывает полученное сообщение
func (c *TransactionConsumer) processMessage(ctx context.Context, message kafka.Message) error {
	var txMsg TransactionMessage
	if err := json.Unmarshal(message.Value, &txMsg); err != nil {
		return fmt.Errorf("failed to unmarshal transaction message: %w", err)
	}

	// Преобразуем строковые ID в UUID
	transactionID, err := uuid.Parse(txMsg.ID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID format: %w", err)
	}

	userID, err := uuid.Parse(txMsg.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Создаем объект транзакции
	transaction := &entities.Transaction{
		ID:              transactionID,
		UserID:          userID,
		Amount:          txMsg.Amount,
		Timestamp:       txMsg.Timestamp,
		Category:        txMsg.Category,
		DiscountApplied: txMsg.DiscountApplied,
	}

	// Обновляем метрики пользователя
	if err := c.updateUserMetrics(ctx, transaction); err != nil {
		return fmt.Errorf("failed to update user metrics: %w", err)
	}

	log.Printf("Processed transaction %s for user %s", transactionID, userID)
	return nil
}

// updateUserMetrics обновляет метрики пользователя на основе транзакции
func (c *TransactionConsumer) updateUserMetrics(ctx context.Context, transaction *entities.Transaction) error {
	// Получаем пользователя
	user, err := c.userService.GetUser(ctx, transaction.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found: %s", transaction.UserID)
	}

	// Обновляем время последней активности пользователя
	user.LastActivity = transaction.Timestamp
	if err := c.userService.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Пересчитываем метрики пользователя
	if err := c.metricsService.RecalculateUserMetrics(ctx, transaction.UserID); err != nil {
		return fmt.Errorf("failed to recalculate user metrics: %w", err)
	}

	// Обновляем вероятность оттока
	if err := c.metricsService.UpdateChurnProbability(ctx, transaction.UserID); err != nil {
		log.Printf("Warning: failed to update churn probability: %v", err)
		// Не возвращаем ошибку, так как это некритичная операция
	}

	return nil
}
