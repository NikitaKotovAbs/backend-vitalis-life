package main

import (
    "backend/internal/adapters/yookassa"
    appPayment "backend/internal/app/payment"
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
    logger.Init("info")
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := godotenv.Load(".env")
    if err != nil {
        logger.Error("Ошибка загрузки переменных окружения", zap.Error(err))
        os.Exit(1)
    }

    cfg, err := config.Load()
    if err != nil {
        logger.Error("Ошибка загрузки конфигурации", zap.Error(err))
        os.Exit(1)
    }

    logger.Info("Конфигурация успешно загружена",
        zap.String("config_level", cfg.Logger.Level),
        zap.String("version", cfg.Server.Version),
    )
    
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

    productRepo := db.NewUserRepository(connDb)
    productService := product.NewService(productRepo)

    // Получение переменных окружения для ЮKassa
    yookassaShopID := os.Getenv("YOOKASSA_SHOP_ID")
    yookassaSecretKey := os.Getenv("YOOKASSA_SECRET_KEY")

    // Валидация
    if yookassaShopID == "" {
        logger.Fatal("YOOKASSA_SHOP_ID не установлен в .env")
    }
    if yookassaSecretKey == "" {
        logger.Fatal("YOOKASSA_SECRET_KEY не установлен в .env")
    }

    // Инициализация репозитория ЮKassa
    yookassaRepo := yookassa.NewPaymentRepository(yookassaShopID, yookassaSecretKey)
    
    // Инициализация сервиса платежей - передаем репозиторий
    paymentService := appPayment.NewService(yookassaRepo)

    router := adaptersHttp.Router(productService, paymentService, cfg)
    
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

    logger.Info("Приложение завершает работу",
        zap.String("version", cfg.Server.Version))
}