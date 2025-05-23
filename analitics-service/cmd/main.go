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

	"analitics-service/config"
	"analitics-service/internal/infrastructure/services"
	"analitics-service/pkg/logger"
)

func main() {
	// Загрузка конфигурации из файла config.yaml
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация логгера
	logg := logger.NewLogger(cfg.Logger.Level)
	logg.Info(context.Background(), "Логгер успешно инициализирован")

	// Инициализация базы данных
	db, err := sql.Open("postgres", cfg.Database.DSN)
	if err != nil {
		logg.Error(context.Background(), "Failed to connect to database", "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logg.Error(ctx, "Failed to ping database", "error", err)
		log.Fatalf("Failed to ping database: %v", err)
	}
	logg.Info(ctx, "Successfully connected to database")

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMinutes) * time.Minute)

	// Инициализация сервисов
	aprioriService := services.NewAprioriService(logg)
	logg.Info(ctx, "Services initialized successfully")

	// Инициализация HTTP роутера
	// TODO: Добавить HTTP обработчики
	router := http.NewServeMux()
	logg.Info(ctx, "HTTP router setup completed")

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
		logg.Info(ctx, "Starting HTTP server", "address", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Error(ctx, "Failed to start server", "error", err)
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logg.Info(ctx, "Shutting down server...")

	// Создаем контекст с таймаутом для graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Останавливаем HTTP сервер
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logg.Error(shutdownCtx, "Server forced to shutdown", "error", err)
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logg.Info(shutdownCtx, "Server exited properly")
}
