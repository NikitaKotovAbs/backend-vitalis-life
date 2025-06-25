package db

import (
	"backend/internal/domain/basket"
	"backend/pkg/logger"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type BasketRepository struct {
	db *sql.DB
}

func NewBasketRepository(db *sql.DB) *BasketRepository {
	return &BasketRepository{db: db}
}

func (r *BasketRepository) Add(data basket.Basket) error {
	logger.Debug("Начало добавления товара в корзину",
		zap.Int("product_id", data.ProductID),
		zap.Int("user_id", data.UID),
	)

	tx, err := r.db.Begin()
	if err != nil {
		logger.Error("Ошибка начала транзакции",
			zap.Error(err))
		return fmt.Errorf("ошибка сервера")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query := `
		INSERT INTO baskets (product_id, uid, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (product_id, uid) 
		DO UPDATE SET quantity = baskets.quantity + EXCLUDED.quantity
	`

	result, err := tx.Exec(query,
		data.ProductID,
		data.UID,
		data.Quantity,
	)
	if err != nil {
		logger.Error("Ошибка добавления товара в корзину",
			zap.Int("product_id", data.ProductID),
			zap.Int("user_id", data.UID),
			zap.Error(err))
		return fmt.Errorf("ошибка добавления товара")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Ошибка проверки измененных строк",
			zap.Error(err))
		return fmt.Errorf("ошибка сервера")
	}

	if rowsAffected == 0 {
		logger.Debug("Товар не был добавлен в корзину",
			zap.Int("product_id", data.ProductID))
		return fmt.Errorf("товар не добавлен")
	}

	logger.Debug("Товар успешно добавлен в корзину",
		zap.Int("product_id", data.ProductID),
		zap.Int64("rows_affected", rowsAffected))

	return nil
}