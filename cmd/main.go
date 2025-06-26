package main

import (
	"backend/config"
	"backend/internal/adapters/db"
	adaptersHttp "backend/internal/adapters/http"
	"backend/internal/app/product"
	"backend/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	_ "github.com/lib/pq"
)

func main() {

	// Инициализация логгера
	logger.Init("info")

	// Канал для перехвата сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Загрузка переменных окружения
	// err := godotenv.Load(".env")
	// if err != nil {
    //     logger.Error("Ошибка загрузки переменных окружения",
	// 	zap.Error(err),
	// 	zap.String("source", "main"))
	// 	os.Exit(1)
	// }

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Ошибка загрузки конфигурации",
			zap.Error(err),
			zap.String("source", "main"))
		os.Exit(1)
	}

	logger.Info("Конфигурация успешно загружена",
		zap.Any("config_level", cfg.Logger.Level),
        zap.Any("version", cfg.Server.Version),
        )
	
	// Переинициализация логгера
	logger.Init(cfg.Logger.Level)
	
	connData := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
	)

	connDb, err := sql.Open("postgres", connData)
	if err != nil {
    	log.Fatal(err)
	}
	defer connDb.Close()

	// Инициализация репозитория (адаптер БД)
	productRepo := db.NewUserRepository(connDb) // В реальности тут будет подключение к PostgreSQL/MySQL и т.д.

	// Инициализация сервиса (бизнес-логика)
	productService := product.NewService(productRepo)

	// Инициализация HTTP-роутера (Gin)
	router := adaptersHttp.Router(productService, cfg)
	
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	// Запуск сервера
	go func() {
		logger.Info("Запуск сервера",
			zap.String("port", cfg.Server.Port),
			zap.String("mode", gin.Mode()),
		)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Ошибка запуска сервера",
				zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	<-quit
	logger.Info("Получен сигнал завершения, остановка сервера...")

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Ошибка при остановке сервера",
			zap.Error(err))
	} else {
		logger.Info("Сервер корректно остановлен")
	}

	// Дополнительные операции перед выходом
	logger.Info("Приложение завершает работу",
		zap.String("version", cfg.Server.Version))
}