package handlers

import (
	"backend/internal/app/product"
	"backend/pkg/logger"
	"net/http"
    "strconv"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductHandler struct {
    service *product.Service
}

func NewProductHandler(service *product.Service) *ProductHandler {
    return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
    products, err := h.service.GetAllProducts()
    if err != nil {
        logger.Error("Ошибка при получении данных",
            zap.Error(err),
        )
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
        return
    }
    
    c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetByIdProducts(c *gin.Context) {
    idParam := c.Param("id")

    id, err := strconv.Atoi(idParam)
    if err != nil {
        logger.Error("Ошибка конвертации id",
            zap.Error(err),
        )
        c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сервера"})
        return
    }

    products, err := h.service.GetByIdProducts(id)
    if err != nil {
        logger.Error("Ошибка при получении данных",
            zap.Error(err),
        )
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
        return
    }
    
    c.JSON(http.StatusOK, products)
}