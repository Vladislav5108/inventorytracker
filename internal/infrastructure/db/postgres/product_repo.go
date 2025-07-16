package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
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

// GetALL возвращает список всех товаров
func (r *ProductRepo) GetALL() ([]domain.Product, error) {

	query := `SELECT * FROM products`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса списка: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	var p domain.Product

	for rows.Next() {

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Quantity,
			&p.CategoryID,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения: %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка получения данных: %w", err)
	}

	return products, nil
}

//  возвращает товары конкретной категории

func (r *ProductRepo) GetByCategory(CategoryID int) ([]domain.Product, error) {
	query := `SELECT *
			FROM products
			WHERE category_id = $1`

	rows, err := r.db.Query(query, CategoryID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса списка по категории: %w", err)
	}
	defer rows.Close()

	var product []domain.Product

	for rows.Next() {
		var p domain.Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Quantity,
			&p.CategoryID,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %w", err)
		}
		product = append(product, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка получения данных: %w", err)
	}
	return product, nil
}

//добавление нового товара!

func (r *ProductRepo) Add(product domain.Product) (int, error) {
	query := `INSERT INTO products(
	name,
	price,
	quantity,
	category_id,
	is_archived
	)
	VALUES(
	$1,
	$2,
	$3,
	$4,
	false
	)
	RETURNING id`

	var id int

	err := r.db.QueryRow(query,
		product.Name,
		product.Price,
		product.Quantity,
		product.CategoryID).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ошибка добавления товара: %w", err)
	}
	return id, nil
}

// обновляем данные товара

func (r *ProductRepo) UpDate(product domain.Product) error {
	query := `UPDATE products
			SET
				name = $1,
			price = $2,
			quantity = $3,
			category_id = $4
			WHERE 
			id = $5
			AND
			is_archived = false
	`
	result, err := r.db.Exec(
		query,
		product.Name,
		product.Price,
		product.Quantity,
		product.CategoryID,
		product.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка изменнения файлов: %w", err)
	}
	rowsUpdate, _ := result.RowsAffected()
	if rowsUpdate == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

// убераем товар из каталога но не из базы
func (r *ProductRepo) Archiv(id int) error {
	query := `UPDATE products
			SET is_archived = true
			WHERE id = $1
			 AND is_archived = false`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка при добавление в архив: %w", err)
	}
	rowsArchiv, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при обновлении: %w ", err)
	}
	if rowsArchiv == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

// восстановление их архива

func (r *ProductRepo) Restore(id int) error {
	query := `UPDATE products
			SET is_archived = false
			WHERE id = $1
			 AND is_archived = true`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении из архива: %w", err)
	}
	rowsArchiv, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при востонавлении из архива: %w ", err)
	}
	if rowsArchiv == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

// возвращаем список скрытых товаров

func (r *ProductRepo) GetArchived() ([]domain.Product, error) {
	query := `SELECT * 
	FROM products
	WHERE is_archived = true`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения архивного списка: %w", err)
	}
	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var p domain.Product

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.CategoryID, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}
	return products, nil
}

