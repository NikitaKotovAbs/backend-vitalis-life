package handlers

import (
	// "backend/internal/app/basket"
	// "backend/pkg/logger"
	// "net/http"
	// "backend/internal/domain/basket"
	// "github.com/gin-gonic/gin"
	// "go.uber.org/zap"
)

// type BasketHandler struct {
// 	service *basket.Service
// }

// func NewBasketHandler(service *basket.Service) *BasketHandler{
// 	return &BasketHandler{service: service}
// }

// func (h *BasketHandler) AddAllProducts(c *gin.Context) {
// 	var data b

//     products, err := h.service.Add()
//     if err != nil {
//         logger.Debug("error",
//             zap.Error(err),
//         )
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }
    
//     c.JSON(http.StatusOK, products)
// }