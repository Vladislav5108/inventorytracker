package postgres

import (
	"database/sql"
	"errors"

	"github.com/Vladislav5108/inventorytracker/internal/domain"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

// GetByID возвращает товар по ID
func (r *ProductRepo) GetByID(id int) (domain.Product, error) {
	var product domain.Product

	query := `SELECT id, name, price, quantity, category_id, created_at
	FROM products
	WHERE id = $1 AND is_archived = false
	`
	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Quantity,
		&product.CategoryID,
		&product.CreatedAt,
	)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Product{}, domain.ErrProductNotFound
	case err != nil:
		return domain.Product{}, err
	}
	return product, nil
}
