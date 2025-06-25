package basket

import (
	"backend/internal/domain/basket"
)

type Service struct {
	repo basket.BasketRepository
}

func NewService(repo basket.BasketRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(data basket.Basket) error {
	return s.repo.Add(data)
}