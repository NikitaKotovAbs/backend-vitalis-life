package handlers

import (
    "net"
    "net/http"
    "os"
	"encoding/json"
    "strings"
    "backend/pkg/logger"
    "backend/pkg/smtp_sender"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type WebhookHandler struct {
    // Можно добавить зависимости если нужно
}

func NewWebhookHandler() *WebhookHandler {
    return &WebhookHandler{}
}

func (h *WebhookHandler) HandlePaymentWebhook(c *gin.Context) {
    // Проверяем IP адрес отправителя
    clientIP := c.ClientIP()
    
    // Разрешенные IP адреса ЮKassa (официальные из документации)
    allowedIPs := []string{
        "185.71.76.0/27",    // ЮKassa диапазон 1
        "185.71.77.0/27",    // ЮKassa диапазон 2  
        "77.75.153.0/25",    // ЮKassa диапазон 3
        "77.75.154.128/25",  // ЮKassa диапазон 4
        "2a02:5180::/32",    // IPv6 диапазон
    }
    
    if !isIPAllowed(clientIP, allowedIPs) {
        logger.Warn("Webhook from unauthorized IP", 
            zap.String("ip", clientIP),
            zap.String("path", c.Request.URL.Path),
            zap.String("method", c.Request.Method))
        c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
        return
    }

    var notification struct {
        Event  string                 `json:"event"`
        Object map[string]interface{} `json:"object"`
    }

    if err := c.BindJSON(&notification); err != nil {
        logger.Error("Invalid webhook data", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook data"})
        return
    }

    logger.Info("Webhook received", 
        zap.String("event", notification.Event),
        zap.String("payment_id", notification.Object["id"].(string)),
        zap.String("source_ip", clientIP))

    // Обрабатываем только успешные платежи
    if notification.Event == "payment.succeeded" {
        h.handleSuccessfulPayment(notification.Object)
    }

    // Всегда отвечаем 200 OK на вебхуки
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *WebhookHandler) handleSuccessfulPayment(paymentData map[string]interface{}) {
    // Извлекаем данные из metadata
    metadata, ok := paymentData["metadata"].(map[string]interface{})
    if !ok {
        logger.Error("Invalid metadata format in webhook")
        return
    }

    

    // Извлекаем данные
    email, _ := metadata["email"].(string)
    customerName, _ := metadata["customerName"].(string)
    phone, _ := metadata["phone"].(string)
    deliveryType, _ := metadata["deliveryType"].(string)
    deliveryAddress, _ := metadata["deliveryAddress"].(string)
    comment, _ := metadata["comment"].(string)
    cartItemsJSON, _ := metadata["cartItems"].(string)

    amountData, _ := paymentData["amount"].(map[string]interface{})
    amount, _ := amountData["value"].(string)
    currency, _ := amountData["currency"].(string)
    description, _ := paymentData["description"].(string)
    paymentID, _ := paymentData["id"].(string)

    var cartItems []smtp_sender.CartItem
    if cartItemsJSON != "" {
        if err := json.Unmarshal([]byte(cartItemsJSON), &cartItems); err != nil {
            logger.Error("Failed to parse cart items", 
                zap.Error(err),
                zap.String("cart_items_json", cartItemsJSON))
            cartItems = []smtp_sender.CartItem{}
        } else {
            logger.Info("Cart items parsed successfully",
                zap.Int("items_count", len(cartItems)),
                zap.Any("items", cartItems))
        }
    } else {
        logger.Warn("Empty cart items JSON in webhook")
    }

    

    // Формируем данные заказа
    order := smtp_sender.OrderData{
        CustomerName:    customerName,
        Email:           email,
        Phone:           phone,
        DeliveryType:    deliveryType,
        DeliveryAddress: deliveryAddress,
        Comment:         comment,
        PaymentID:       paymentID,
        Amount:          amount,
        Currency:        currency,
        Description:     description,
        CartItems:       cartItems, // Теперь это массив товаров, а не строка
    }

    // Получаем email менеджера из переменных окружения
    managerEmail := os.Getenv("MANAGER_EMAIL")
    if managerEmail == "" {
        managerEmail = "orders@vitalis-life.ru" // email по умолчанию
    }

    // Отправляем письма в горутине (асинхронно)
    go func() {
        err := smtp_sender.SendOrderEmails(order, managerEmail)
        if err != nil {
            logger.Error("Failed to send order emails", 
                zap.Error(err),
                zap.String("client_email", email),
                zap.String("payment_id", paymentID))
        } else {
            logger.Info("Order emails sent successfully",
                zap.String("client_email", email),
                zap.String("payment_id", paymentID))
        }
    }()
}

// isIPAllowed проверяет, разрешен ли IP адрес
func isIPAllowed(ipStr string, allowedIPs []string) bool {
    // Пропускаем локальные адреса для тестирования
    if ipStr == "::1" || ipStr == "127.0.0.1" || strings.HasPrefix(ipStr, "192.168.") {
        logger.Debug("Local IP allowed for testing", zap.String("ip", ipStr))
        return true
    }
    
    clientIP := net.ParseIP(ipStr)
    if clientIP == nil {
        logger.Warn("Invalid IP address", zap.String("ip", ipStr))
        return false
    }

    for _, allowedIP := range allowedIPs {
        // Проверяем CIDR диапазон
        if strings.Contains(allowedIP, "/") {
            _, ipNet, err := net.ParseCIDR(allowedIP)
            if err != nil {
                logger.Error("Invalid CIDR format", 
                    zap.String("cidr", allowedIP),
                    zap.Error(err))
                continue
            }
            
            if ipNet.Contains(clientIP) {
                return true
            }
        } else {
            // Простая проверка точного IP
            if allowedIP == ipStr {
                return true
            }
        }
    }

    return false
}

// Дополнительная функция для логирования всех входящих запросов (для отладки)
func (h *WebhookHandler) logWebhookRequest(c *gin.Context) {
    logger.Debug("Webhook request",
        zap.String("ip", c.ClientIP()),
        zap.String("method", c.Request.Method),
        zap.String("path", c.Request.URL.Path),
        zap.String("user_agent", c.Request.UserAgent()),
        zap.Any("headers", c.Request.Header))
}