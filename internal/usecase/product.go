package usecase

import (
	"fmt"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/Vladislav5108/inventorytracker/internal/domain/repository"
)

type ProductUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

func (uc *ProductUseCase) GetByID(id int) (domain.Product, error) {
	if id <= 0 {
		return domain.Product{}, domain.ErrNotId
	}
	product, err := uc.repo.GetByID(id)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}
func (uc *ProductUseCase) GetALL() ([]domain.Product, error) {
	products, err := uc.repo.GetALL()
	if err != nil {
		return nil, fmt.Errorf("ошибка, не удалось получить список продуктов: %w", err)
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("список пуст: %w", err)
	}
	return products, nil
}
func (uc *ProductUseCase) GetByCategory(CategoryID int) ([]domain.Product, error) {
	if CategoryID <= 0 {
		return nil, domain.ErrNotId
	}
	products, err := uc.repo.GetByCategory(CategoryID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении продуктов по категории:%w", err)
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("ошибка, пустой список: %w", err)
	}
	return products, nil
}

