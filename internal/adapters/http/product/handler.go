package product

import (
    "net/http"
    "backend/internal/app/product"
    "github.com/gin-gonic/gin"
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, products)
}