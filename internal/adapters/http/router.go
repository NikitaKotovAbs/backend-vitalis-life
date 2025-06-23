package http

import (
	appProduct "backend/internal/app/product"
    httpProduct "backend/internal/adapters/http/product"
	"backend/pkg/logger"
	"backend/config"
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
)


func Router(productService *appProduct.Service, cfg *config.Config) *gin.Engine {
    router := gin.Default()
    
    // Добавляем CORS middleware из конфига
    router.Use(CORSNew(cfg.CORS))
    
    // Добавляем middleware для логирования запросов
    router.Use(func(c *gin.Context) {
        logger.Debug("HTTP запрос",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path))
        c.Next()
    })
    
    productHandler := httpProduct.NewProductHandler(productService)

    public := router.Group("/api/v1/public")
    {
        product := public.Group("/product")  
        {
            product.GET("/", productHandler.GetAllProducts)
        }
    }
    
    return router
}