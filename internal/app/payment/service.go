package payment

import (
    domainPayment "backend/internal/domain/payment"
)

type Service struct {
    repo domainPayment.PaymentRepository // Используем интерфейс из domain
}

// NewService теперь принимает репозиторий
func NewService(repo domainPayment.PaymentRepository) *Service {
    return &Service{
        repo: repo,
    }
}

func (s *Service) CreatePayment(req *domainPayment.PaymentRequest) (*domainPayment.PaymentResponse, error) {
    return s.repo.CreatePayment(req)
}

func (s *Service) GetPaymentStatus(paymentID string) (*domainPayment.PaymentResponse, error) {
    return s.repo.GetPaymentStatus(paymentID)
}

func (s *Service) CancelPayment(paymentID string) error {
    return s.repo.CancelPayment(paymentID)
}