package handlers

import (
	"backend/internal/app/product"
	"backend/pkg/logger"
	"net/http"

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
        logger.Debug("error",
            zap.Error(err),
        )
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, products)
}