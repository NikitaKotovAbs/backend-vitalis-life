package db

import (
	"database/sql"
    "backend/internal/domain/product"
)

// UserRepository реализует интерфейс domain/user.UserRepository.
type ProductRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// FindByID возвращает пользователя по ID (заглушка для примера).
func (r *ProductRepository) GetAll() ([]*product.Product, error) {
    query := `SELECT id, name, price, description, discount, img, created_at FROM product`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []*product.Product
    for rows.Next() {
        var p product.Product
        var img sql.Null[[]byte] // Для обработки NULL в img
        
        err := rows.Scan(
            &p.ID,
            &p.Name,
            &p.Price,
            &p.Description,
            &p.Discount,
            &img,
            &p.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        
        if img.Valid {
            p.Image = img.V
        }
        
        products = append(products, &p)
    }
    
    return products, nil
}
