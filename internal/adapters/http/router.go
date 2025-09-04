package http

import (
    appPayment "backend/internal/app/payment"
    appProduct "backend/internal/app/product"
    "backend/internal/adapters/http/handlers"
    "backend/pkg/logger"
    "backend/config"
    "go.uber.org/zap"
    "github.com/gin-gonic/gin"
)

func Router(
    productService *appProduct.Service, 
    paymentService *appPayment.Service, 
    cfg *config.Config,
) *gin.Engine {
    router := gin.Default()
    
    router.Use(CORSNew(cfg.CORS))
    
    router.Use(func(c *gin.Context) {
        logger.Debug("HTTP запрос",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path))
        c.Next()
    })
    
    productHandler := handlers.NewProductHandler(productService)
    paymentHandler := handlers.NewPaymentHandler(paymentService)
    webhookHandler := handlers.NewWebhookHandler()

    public := router.Group("/api/v1/public")
    {
        product := public.Group("/product")  
        {
            product.GET("/", productHandler.GetAllProducts)
            product.GET("/:id", productHandler.GetByIdProducts)
        }

        payment := public.Group("/payment")
        {
            payment.POST("/create", paymentHandler.CreatePayment)
            payment.GET("/:id/status", paymentHandler.GetStatus)      // Исправлено на GetStatus
            payment.POST("/:id/cancel", paymentHandler.Cancel)        // Исправлено на Cancel
        }
    }
    router.POST("/webhook/payment", webhookHandler.HandlePaymentWebhook)
    return router
}