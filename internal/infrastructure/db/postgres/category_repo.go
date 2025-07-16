package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

//Создание новой категории

func (c *CategoryRepo) CreateCategory(category domain.Category) (int, error) {
	query := `INSERT INTO category(
	 name,description
	 )
	 VALUES(
	 $1,
	 $2
	 )
	 RETURNING id`

	var id int

	err := c.db.QueryRow(query, category.Name, category.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания категории: %w", err)
	}
	return id, nil
}

// получениу категории
func (c *CategoryRepo) GetByIDCategory(id int) (domain.Category, error) {
	query := `SELECT  id, name,description
	FROM category
	WHERE id = $1`

	var category domain.Category

	row := c.db.QueryRow(query, id)
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Category{}, domain.ErrNotId
		}
		return domain.Category{}, fmt.Errorf("ошибка чтения данных: %w", err)
	}
	return category, nil
}

func (c *CategoryRepo) GetALLCategory() ([]domain.Category, error) {
	query := `SELECT * FROM category`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса списка категорий:%w", err)
	}
	defer rows.Close()

	var categories []domain.Category

	for rows.Next() {
		var ca domain.Category
		err := rows.Scan(
			&ca.ID,
			&ca.Name,
			&ca.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %w", err)
		}
		categories = append(categories, ca)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка аолучения спика :%w", err)
	}
	return categories, nil
}

func (c *CategoryRepo) UpDateCategory(category domain.Category) error {
	query := `UPDATE category
	SET
	name = $1,
	description = $2
	WHERE id = $3`

	result, err := c.db.Exec(
		query,
		category.Name,
		category.Description,
		category.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка запроса обновления категории: %w", err)
	}
	rowsUpdate, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки на возврат количества строк: %w", err)
	}
	if rowsUpdate == 0 {
		return domain.ErrNotId
	}
	return nil
}

func (c *CategoryRepo) DeleteCategory(id int) error {
	query := `DELETE FROM category WHERE id = $1`

	result, err := c.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления:%w", err)
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки удаления: %w", err)
	}
	if rowsDeleted == 0 {
		return domain.ErrNotId
	}
	return nil
}

