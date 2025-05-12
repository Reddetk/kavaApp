// cmd/main.go
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // Postgres driver

	"user-service/config"
	"user-service/internal/application"
	"user-service/internal/infrastructure/postgres"
	"user-service/internal/infrastructure/services"

	httpInterface "user-service/internal/interfaces/http"
	"user-service/internal/interfaces/http/handlers"
	"user-service/internal/interfaces/kafka"
	"user-service/pkg/logger"
)

func main() {
	// Загрузка конфигурации из файла config.yaml
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
		logg.Error("Failed to connect to database", "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logg.Error("Failed to ping database", "error", err)
		log.Fatalf("Failed to ping database: %v", err)
	}
	logg.Info("Successfully connected to database")

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMinutes) * time.Minute)

	// Инициализация репозиториев
	userRepo := postgres.NewUserRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	segmentRepo := postgres.NewSegmentRepository(db)
	metricsRepo := postgres.NewUserMetricsRepository(db)

	// Проверка соединения с репозиториями
	if err := userRepo.Ping(context.Background()); err != nil {
		logg.Error("Failed to ping user repository", "error", err)
		log.Fatalf("Failed to ping user repository: %v", err)
	}
	logg.Info("User repository connection test passed")

	// Инициализация сервисов
	segmentationSvc := services.NewKMeansSegmentation(cfg.Segmentation.RFMClustering.Clusters)
	survivalSvc := services.NewCoxSurvivalAnalysis()
	transitionSvc := services.NewMarkovTransitionService()
	clvSvc := services.NewDiscountedCLVService(userRepo, transactionRepo)
	logg.Info("Services initialized successfully")

	// Инициализация use cases
	userService := application.NewUserService(userRepo, metricsRepo)
	retentionService := application.NewRetentionService(userRepo, metricsRepo, survivalSvc, transitionSvc)
	clvService := application.NewCLVService(userRepo, metricsRepo, clvSvc, retentionService)
	segmentationService := application.NewSegmentationService(userRepo, segmentRepo, metricsRepo, transactionRepo, segmentationSvc)
	logg.Info("Application services initialized successfully")

	// Инициализация handlers
	userHandler := handlers.NewUserHandler(userService)
	segmentHandler := handlers.NewSegmentHandler(segmentationService)
	clvHandler := handlers.NewCLVHandler(clvService)
	retentionHandler := handlers.NewRetentionHandler(retentionService)
	logg.Info("HTTP handlers initialized successfully")

	// Инициализация HTTP роутера
	router := httpInterface.SetupRouter(userHandler, segmentHandler, clvHandler, retentionHandler)
	logg.Info("HTTP router setup completed")

	// Инициализация Kafka consumer
	kafkaBrokers := cfg.Kafka.Brokers
	if len(kafkaBrokers) == 0 {
		kafkaBrokers = []string{"localhost:9092"}
	}

	transactionConsumer := kafka.NewTransactionConsumer(
		kafkaBrokers,
		cfg.Kafka.Topic,
		userService,
		retentionService,
	)
	logg.Info("Kafka consumer initialized", "brokers", kafkaBrokers, "topic", cfg.Kafka.Topic)

	// Запуск Kafka consumer в горутине
	go func() {
		logg.Info("Starting Kafka consumer")
		if err := transactionConsumer.Start(context.Background()); err != nil {
			logg.Error("Failed to start Kafka consumer", "error", err)
			log.Fatalf("Failed to start Kafka consumer: %v", err)
		}
	}()

	// Запуск HTTP сервера
	serverAddr := cfg.Server.Address
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeoutSeconds) * time.Second,
	}

	go func() {
		logg.Info("Starting HTTP server", "address", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Error("Failed to start server", "error", err)
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logg.Info("Shutting down server...")

	// Создаем контекст с таймаутом для graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Сначала останавли��аем Kafka consumer
	if err := transactionConsumer.Stop(shutdownCtx); err != nil {
		logg.Error("Failed to stop Kafka consumer", "error", err)
	} else {
		logg.Info("Kafka consumer stopped successfully")
	}

	// Затем останавливаем HTTP сервер
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logg.Error("Server forced to shutdown", "error", err)
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logg.Info("Server exited properly")
}
