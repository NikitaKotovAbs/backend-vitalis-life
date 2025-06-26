package product

// UserRepository определяет контракт для работы с хранилищем пользователей.
type ProductRepository interface {
    GetAll() ([]*Product, error)
    GetByID(id int) (*Product, error)
}