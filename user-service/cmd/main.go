// cmd/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	_ "net/http"
	_ "os"
	_ "os/signal"
	_ "syscall"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"user-service/config"
	_ "user-service/internal/application"
	"user-service/internal/infrastructure/postgres"
	_ "user-service/internal/infrastructure/services"

	// httpInterface "user-service/internal/interfaces/http"
	_ "user-service/internal/interfaces/http/handlers"
	_ "user-service/internal/interfaces/kafka"
	"user-service/pkg/logger"
)

func main() {
	//	Загрузка конфигурации из файла config.yaml
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация логгера
	logg := logger.NewLogger(cfg.Logger.Level)
	logg.Info("Логгер успешно инициализирован")

	// Инициализация базы данных
	db, err := sql.Open("postgres", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	userRepo := postgres.NewUserRepository(db)
	// Создаем контекст с таймаутом в 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Преобразуем строку в UUID
	userID, err := uuid.Parse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15")
	if err != nil {
		logg.Error("Ошибка при парсинге UUID:", err)
		return
	}

	// Получаем пользователя по ID
	user, err := userRepo.Get(ctx, userID)
	if err != nil {
		logg.Error("Ошибка при получении пользователя:", err)
		return
	}

	// Проверяем, найден ли пользователь
	if user == nil {
		logg.Info("Пользователь с ID a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15 не найден")
		return
	}

	// Выводим информацию о пользователе
	logg.Info("Информация о пользователе:")
	logg.Info(fmt.Sprintf("ID: %s", user.ID))
	logg.Info(fmt.Sprintf("Email: %s", user.Email))
	logg.Info(fmt.Sprintf("Телефон: %s", user.Phone))
	logg.Info(fmt.Sprintf("Возраст: %d", user.Age))
	logg.Info(fmt.Sprintf("Пол: %s", user.Gender))
	logg.Info(fmt.Sprintf("Город: %s", user.City))
	logg.Info(fmt.Sprintf("Дата регистрации: %s", user.RegistrationDate.Format(time.RFC3339)))
	logg.Info(fmt.Sprintf("Последняя активность: %s", user.LastActivity.Format(time.RFC3339)))

	// transactionRepo := postgres.NewTransactionRepository(db)
	// segmentRepo := postgres.NewSegmentRepository(db)
	// metricsRepo := postgres.NewUserMetricsRepository(db)

	// // Инициализация сервисов
	// segmentationSvc := services.NewKMeansSegmentation()
	// survivalSvc := services.NewCoxSurvivalAnalysis()
	// transitionSvc := services.NewMarkovTransitionService()
	// clvSvc := services.NewDiscountedCLVService()

	// // Инициализация use cases
	// userService := application.NewUserService(userRepo, metricsRepo)
	// segmentationService := application.NewSegmentationService(userRepo, segmentRepo, metricsRepo, transactionRepo, segmentationSvc)
	// retentionService := application.NewRetentionService(userRepo, metricsRepo, survivalSvc, transitionSvc)
	// clvService := application.NewCLVService(userRepo, metricsRepo, clvSvc, retentionService)

	// // Инициализация handlers
	// userHandler := handlers.NewUserHandler(userService)
	// segmentHandler := handlers.NewSegmentHandler(segmentationService)

	// // Инициализация HTTP роутера
	// router := httpInterface.SetupRouter(userHandler, segmentHandler)

	// // Инициализация Kafka consumer
	// transactionConsumer := kafka.NewTransactionConsumer([]string{"localhost:9092"}, "transactions", userService, retentionService)

	// // Запуск Kafka consumer в горутине
	// go func() {
	// 	if err := transactionConsumer.Start(context.Background()); err != nil {
	// 		log.Fatalf("Failed to start Kafka consumer: %v", err)
	// 	}
	// }()

	// // Запуск HTTP сервера
	// srv := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: router,
	// }

	// go func() {
	// 	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("Failed to start server: %v", err)
	// 	}
	// }()

	// Graceful shutdown
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit

	// log.Println("Shutting down server...")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatalf("Server forced to shutdown: %v", err)
	// }

	// log.Println("Server exited properly")
}
