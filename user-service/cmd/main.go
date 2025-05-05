// cmd/main.go
package main

import (
	_ "context"
	"database/sql"
	"log"
	_ "net/http"
	_ "os"
	_ "os/signal"
	_ "syscall"
	_ "time"

	_ "github.com/lib/pq"

	"user-service/config"
	_ "user-service/internal/application"
	_ "user-service/internal/infrastructure/postgres"
	_ "user-service/internal/infrastructure/services"

	// httpInterface "user-service/internal/interfaces/http"
	_ "user-service/internal/interfaces/http/handlers"
	_ "user-service/internal/interfaces/kafka"
	"user-service/pkg/logger"
)

func main() {
	//	Загрузка конфигурации из файла config.yaml
	cfg, err := config.LoadConfig("config.yaml")
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

	// // Инициализация репозиториев
	// userRepo := postgres.NewUserRepository(db)
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
