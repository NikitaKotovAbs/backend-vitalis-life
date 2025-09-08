package product

import 
(
	"backend/internal/domain/product"
	// "errors"
)

// Service содержит бизнес-логику работы с пользователями.
type Service struct {
    repo product.ProductRepository
}

func NewService(repo product.ProductRepository) *Service {
    return &Service{repo: repo}
}

func (s *Service) GetAllProducts() ([]*product.Product, error) {
    return s.repo.GetAll()
}

func (s *Service) GetByIdProducts(id int) (*product.Product, error) {
    return s.repo.GetByID(id)
}

func (s *Service) GetProductByID(id int) (*product.Product, error) {
    return s.repo.GetByID(id)
}