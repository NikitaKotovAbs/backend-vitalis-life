package payment

import "time"

// PaymentRequest - запрос на создание платежа
type PaymentRequest struct {
    Amount      float64                   `json:"amount" binding:"required"`
    Description string                    `json:"description"`
    Currency    string                    `json:"currency" binding:"required"`
    ReturnURL   string                    `json:"returnUrl" binding:"required"`
    Email       string                    `json:"email" binding:"required"`
    Phone       string                    `json:"phone" binding:"required"`
    Metadata    map[string]interface{}    `json:"metadata"`
    ReceiptItems []map[string]interface{} `json:"receipt_items"` // Добавляем поле для чека
}

type Receipt struct {
    Customer ReceiptCustomer `json:"customer"`
    Items    []ReceiptItem   `json:"items"`
}

type ReceiptCustomer struct {
    Email string `json:"email"`
}

type ReceiptItem struct {
    Description    string         `json:"description"`
    Quantity       string         `json:"quantity"`
    Amount         Amount         `json:"amount"`
    VatCode        string         `json:"vat_code"`
    PaymentMode    string         `json:"payment_mode"`
    PaymentSubject string         `json:"payment_subject"`
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