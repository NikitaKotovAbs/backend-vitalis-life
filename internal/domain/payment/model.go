package payment

import "time"

// PaymentRequest - запрос на создание платежа
type PaymentRequest struct {
    Amount      float64                `json:"amount" binding:"required"`
    Description string                 `json:"description"`
    Currency    string                 `json:"currency" binding:"required"`
    ReturnURL   string                 `json:"returnUrl" binding:"required"`
    Email       string                 `json:"email" binding:"required"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// PaymentResponse - ответ от платежной системы
type PaymentResponse struct {
    ID           string                 `json:"id"`
    Status       string                 `json:"status"`
    Amount       Amount                 `json:"amount"`
    Description  string                 `json:"description"`
    Confirmation Confirmation           `json:"confirmation"`
    Metadata     map[string]interface{} `json:"metadata"`
    CreatedAt    time.Time              `json:"created_at"`
}

type Amount struct {
    Value    string `json:"value"`
    Currency string `json:"currency"`
}

type Confirmation struct {
    ConfirmationURL string `json:"confirmation_url"`
    Type            string `json:"type"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details,omitempty"`
}