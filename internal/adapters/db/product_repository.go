package db

import (
	"database/sql"
    "backend/internal/domain/product"
    "fmt"
    "errors"
)

type ProductRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll возвращает все продукты из базы данных
func (r *ProductRepository) GetAll() ([]*product.Product, error) {
    query := `SELECT id, title, price, description, discount, img, created_at FROM product`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("ошибка при получении товаров: %w", err)
    }
    defer rows.Close()

    var products []*product.Product
    for rows.Next() {
        var p product.Product
        var img sql.NullString
        
        if err := rows.Scan(
            &p.ID,
            &p.Title,
            &p.Price,
            &p.Description,
            &p.Discount,
            &img,
            &p.CreatedAt,
        ); err != nil {
            return nil, fmt.Errorf("ошибка при сканировании товаров: %w", err)
        }
        
        if img.Valid {
            p.Image = img.String
        }
        
        products = append(products, &p)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
    }
    
    return products, nil
}

func (r *ProductRepository) GetByID(id int) (*product.Product, error) {
    query := `SELECT id, title, price, description, discount, img, created_at FROM product WHERE id = $1`
    
    var p product.Product
    err := r.db.QueryRow(query, id).Scan(
        &p.ID,
        &p.Title,
        &p.Price,
        &p.Description,
        &p.Discount,
        &p.Image,
        &p.CreatedAt,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("товар с id %d не найден", id)
        }
        return nil, fmt.Errorf("ошибка при получении товара: %w", err)
    }
    
    return &p, nil
}
