package yookassa

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    "math/rand"
    domainPayment "backend/internal/domain/payment"
    "backend/pkg/logger"
    "go.uber.org/zap"
)

type PaymentRepository struct {
    shopID    string
    secretKey string
    baseURL   string
}

func NewPaymentRepository(shopID, secretKey string) *PaymentRepository {
    return &PaymentRepository{
        shopID:    shopID,
        secretKey: secretKey,
        baseURL:   "https://api.yookassa.ru/v3",
    }
}

func (r *PaymentRepository) CreatePayment(request *domainPayment.PaymentRequest) (*domainPayment.PaymentResponse, error) {
    amountValue := fmt.Sprintf("%.2f", request.Amount)

    paymentData := map[string]interface{}{
        "amount": map[string]string{
            "value":    amountValue,
            "currency": request.Currency,
        },
        "capture": true,
        "confirmation": map[string]string{
            "type":       "redirect",
            "return_url": request.ReturnURL,
        },
        "description": request.Description,
        "metadata":    request.Metadata,
    }

    // ДОБАВЛЯЕМ ЧЕК 54-ФЗ ИЗ ОТДЕЛЬНОГО ПОЛЯ
    if len(request.ReceiptItems) > 0 {
        paymentData["receipt"] = map[string]interface{}{
            "customer": map[string]string{
                "email": request.Email,
            },
            "items": request.ReceiptItems,
        }
    }

    jsonData, err := json.Marshal(paymentData)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal payment data: %w", err)
    }

    // Логируем запрос для отладки
    logger.Debug("Sending request to YooKassa", 
        zap.String("request", string(jsonData)))

    httpReq, err := http.NewRequest(
        "POST", 
        r.baseURL+"/payments", 
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    auth := base64.StdEncoding.EncodeToString([]byte(r.shopID + ":" + r.secretKey))
    httpReq.Header.Set("Authorization", "Basic "+auth)
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Idempotence-Key", fmt.Sprintf("%d", time.Now().UnixNano()))

    client := &http.Client{}
    resp, err := client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    if resp.StatusCode != http.StatusOK {
        logger.Error("YooKassa API error", 
            zap.String("status", resp.Status), 
            zap.String("response", string(body)))
        return nil, fmt.Errorf("YooKassa error: %s - %s", resp.Status, string(body))
    }

    var response domainPayment.PaymentResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    return &response, nil
}

func (r *PaymentRepository) GetPaymentStatus(paymentID string) (*domainPayment.PaymentResponse, error) {
    httpReq, err := http.NewRequest("GET", r.baseURL+"/payments/"+paymentID, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    auth := base64.StdEncoding.EncodeToString([]byte(r.shopID + ":" + r.secretKey))
    httpReq.Header.Set("Authorization", "Basic "+auth)

    client := &http.Client{}
    resp, err := client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    if resp.StatusCode != http.StatusOK {
        logger.Error("YooKassa API error", 
            zap.String("status", resp.Status), 
            zap.String("response", string(body)))
        return nil, fmt.Errorf("YooKassa error: %s - %s", resp.Status, string(body))
    }

    var response domainPayment.PaymentResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    return &response, nil
}

func (r *PaymentRepository) CancelPayment(paymentID string) error {
    httpReq, err := http.NewRequest("POST", r.baseURL+"/payments/"+paymentID+"/cancel", nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    auth := base64.StdEncoding.EncodeToString([]byte(r.shopID + ":" + r.secretKey))
    httpReq.Header.Set("Authorization", "Basic "+auth)
    httpReq.Header.Set("Idempotence-Key", fmt.Sprintf("%d", time.Now().UnixNano()))

    client := &http.Client{}
    resp, err := client.Do(httpReq)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        logger.Error("YooKassa API error", 
            zap.String("status", resp.Status), 
            zap.String("response", string(body)))
        return fmt.Errorf("YooKassa error: %s - %s", resp.Status, string(body))
    }

    return nil
}

func generateOrderID() string {
    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, 8)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return fmt.Sprintf("ORDER-%d-%s", time.Now().Unix(), string(b))
}