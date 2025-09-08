package handlers

import (
	appPayment "backend/internal/app/payment"
	appProduct "backend/internal/app/product" // ПРАВИЛЬНЫЙ ИМПОРТ
	domainPayment "backend/internal/domain/payment"
	"backend/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ОБНОВЛЕННАЯ СТРУКТУРА С ДОБАВЛЕНИЕM PRODUCT SERVICE
type PaymentHandler struct {
	service        *appPayment.Service
	productService *appProduct.Service
}

func NewPaymentHandler(service *appPayment.Service, productService *appProduct.Service) *PaymentHandler {
	return &PaymentHandler{
		service:        service,
		productService: productService,
	}
}

// НОВАЯ СТРУКТУРА ДЛЯ ВХОДЯЩЕГО ЗАПРОСА (только ID и quantity)
type CartItemRequest struct {
	ProductID int `json:"productId" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,gt=0"`
}

// НОВАЯ СТРУКТУРА ДЛЯ ОБОГАЩЕННЫХ ДАННЫХ (с названием и ценой)
type CartItemResponse struct {
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Name      string  `json:"name"`
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var paymentRequest struct {
		Amount          float64           `json:"amount" binding:"required,gt=0"`
		Description     string            `json:"description"`
		Currency        string            `json:"currency" binding:"required,oneof=RUB USD EUR"`
		ReturnURL       string            `json:"returnUrl" binding:"required,url"`
		Email           string            `json:"email" binding:"required,email"`
		Phone           string            `json:"phone" binding:"required"`
		CustomerName    string            `json:"customerName" binding:"required"`
		DeliveryType    string            `json:"deliveryType" binding:"required,oneof=delivery pickup"`
		DeliveryAddress string            `json:"deliveryAddress"`
		Comment         string            `json:"comment"`
		CartItems       []CartItemRequest `json:"cartItems" binding:"required,min=1"` // ИСПОЛЬЗУЕМ НОВУЮ СТРУКТУРУ
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

	// 1. ПОЛУЧАЕМ ПОЛНЫЕ ДАННЫЕ О ТОВАРАХ ИЗ БАЗЫ
	enrichedItems, err := h.enrichCartItems(paymentRequest.CartItems)
	if err != nil {
		logger.Error("Failed to get product details", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных о товарах: " + err.Error()})
		return
	}

	// 2. ПЕРЕСЧИТЫВАЕМ СУММУ НА ОСНОВЕ РЕАЛЬНЫХ ЦЕН
	itemsTotal := 0.0
	for _, item := range enrichedItems {
		itemsTotal += item.Price * float64(item.Quantity)
	}

	// 3. РАССЧИТЫВАЕМ ДОСТАВКУ
	deliveryCost := calculateDeliveryCost(itemsTotal, paymentRequest.DeliveryType)
	totalAmount := itemsTotal + deliveryCost

	// 4. СОЗДАЕМ ОПИСАНИЕ ЗАКАЗА
	description := fmt.Sprintf("Заказ из %d товаров: ", len(enrichedItems))
	for _, item := range enrichedItems {
		description += fmt.Sprintf("%s (x%d), ", item.Name, item.Quantity)
	}
	description = description[:len(description)-2] // Убираем последнюю запятую

	// 5. ПОДГОТАВЛИВАЕМ МЕТАДАННЫЕ С ПОЛНЫМИ ДАННЫМИ
	cartItemsJSON, err := json.Marshal(enrichedItems)
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
		"cartItems":       string(cartItemsJSON), // ТЕПЕРЬ С НАЗВАНИЯМИ И ЦЕНАМИ
		"itemsTotal":      itemsTotal,
		"deliveryCost":    deliveryCost,
	}

	// 6. СОЗДАЕМ ПЛАТЕЖ С ПРАВИЛЬНОЙ СУММОЙ
	domainPaymentReq := &domainPayment.PaymentRequest{
		Amount:      totalAmount, // ПЕРЕСЧИТАННАЯ СУММА
		Description: description,
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

// НОВАЯ ФУНКЦИЯ: ПОЛУЧЕНИЕ ДАННЫХ О ТОВАРАХ ИЗ БАЗЫ
func (h *PaymentHandler) enrichCartItems(cartItems []CartItemRequest) ([]CartItemResponse, error) {
	enrichedItems := make([]CartItemResponse, len(cartItems))

	for i, item := range cartItems {
		// Получаем данные о товаре из базы
		product, err := h.productService.GetProductByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product %d: %w", item.ProductID, err)
		}

		// Используем актуальную цену (со скидкой если есть)
		price := product.Price
		if product.Discount > 0 {
			price = product.Price * (1 - float64(product.Discount)/100)
		}

		enrichedItems[i] = CartItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
			Name:      product.Name, // ТЕПЕРЬ ЗДЕСЬ БУДЕТ НАЗВАНИЕ
		}
	}

	return enrichedItems, nil
}

// ФУНКЦИЯ РАСЧЕТА ДОСТАВКИ
func calculateDeliveryCost(itemsTotal float64, deliveryType string) float64 {
	if deliveryType == "pickup" {
		return 0
	}

	switch {
	case itemsTotal >= 5000:
		return 0
	case itemsTotal >= 3000:
		return itemsTotal * 0.10
	case itemsTotal >= 1000:
		return itemsTotal * 0.15
	default:
		return itemsTotal * 0.20
	}
}

// GetStatus - получение статуса платежа (без изменений)
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

// Cancel - отмена платежа (без изменений)
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