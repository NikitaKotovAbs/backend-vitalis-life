package handlers

import (
	appPayment "backend/internal/app/payment"
	domainPayment "backend/internal/domain/payment"
	"backend/pkg/logger"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	service *appPayment.Service
}

func NewPaymentHandler(service *appPayment.Service) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var paymentRequest struct {
		Amount          float64    `json:"amount" binding:"required,gt=0"`
		Description     string     `json:"description"`
		Currency        string     `json:"currency" binding:"required,oneof=RUB USD EUR"`
		ReturnURL       string     `json:"returnUrl" binding:"required,url"`
		Email           string     `json:"email" binding:"required,email"`
		Phone           string     `json:"phone" binding:"required"`
		CustomerName    string     `json:"customerName" binding:"required"`
		DeliveryType    string     `json:"deliveryType" binding:"required,oneof=delivery pickup"`
		DeliveryAddress string     `json:"deliveryAddress"`
		Comment         string     `json:"comment"`
		CartItems       []CartItem `json:"cartItems" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		logger.Error("Invalid payment request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверные данные запроса",
			"details": err.Error(),
		})
		return
	}

	// Дополнительная валидация: если доставка, то адрес обязателен
	if paymentRequest.DeliveryType == "delivery" && paymentRequest.DeliveryAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Адрес доставки обязателен",
		})
		return
	}

	// Валидация телефона
	if len(paymentRequest.Phone) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректный номер телефона",
		})
		return
	}

	// Валидация корзины
	for _, item := range paymentRequest.CartItems {
		if item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Количество товара должно быть больше 0",
			})
			return
		}
		if item.Price <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Цена товара должна быть больше 0",
			})
			return
		}
	}

	// Преобразуем cartItems в JSON строку
	cartItemsJSON, err := json.Marshal(paymentRequest.CartItems)
	if err != nil {
		logger.Error("Failed to marshal cart items", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки корзины"})
		return
	}

	metadata := map[string]interface{}{
		"email":           paymentRequest.Email,
		"phone":           paymentRequest.Phone,
		"customerName":    paymentRequest.CustomerName,
		"deliveryType":    paymentRequest.DeliveryType,
		"deliveryAddress": paymentRequest.DeliveryAddress,
		"comment":         paymentRequest.Comment,
		"cartItems":       string(cartItemsJSON),
	}

	domainPaymentReq := &domainPayment.PaymentRequest{
		Amount:      paymentRequest.Amount,
		Description: paymentRequest.Description,
		Currency:    paymentRequest.Currency,
		ReturnURL:   paymentRequest.ReturnURL,
		Email:       paymentRequest.Email,
		Metadata:    metadata,
	}

	paymentResp, err := h.service.CreatePayment(domainPaymentReq)
	if err != nil {
		logger.Error("Failed to create payment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания платежа: " + err.Error()})
		return
	}


	c.JSON(http.StatusOK, paymentResp)
}

// GetStatus - получение статуса платежа
func (h *PaymentHandler) GetStatus(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	paymentResp, err := h.service.GetPaymentStatus(paymentID)
	if err != nil {
		logger.Error("Failed to get payment status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment status"})
		return
	}

	c.JSON(http.StatusOK, paymentResp)
}

// Cancel - отмена платежа
func (h *PaymentHandler) Cancel(c *gin.Context) {
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	err := h.service.CancelPayment(paymentID)
	if err != nil {
		logger.Error("Failed to cancel payment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment cancelled successfully"})
}

type CartItem struct {
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Name      string  `json:"name"`
}
