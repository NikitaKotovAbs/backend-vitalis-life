package payment

// PaymentRepository определяет контракт для работы с платежами
type PaymentRepository interface {
    CreatePayment(request *PaymentRequest) (*PaymentResponse, error)
    GetPaymentStatus(paymentID string) (*PaymentResponse, error)
    CancelPayment(paymentID string) error
}