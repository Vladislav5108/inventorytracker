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
		return domain.Product{}, domain.ErrProductNotFound
	}
	return uc.repo.GetByID(id)
}

func (uc *ProductUseCase) GetAll() ([]domain.Product, error) {
	products, err := uc.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting products: %w", err)
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("the list is empty: %w", err)
	}
	return products, nil
}
func (uc *ProductUseCase) GetByCategory(CategoryID int) ([]domain.Product, error) {
	if CategoryID <= 0 {
		return nil, domain.ErrNotFoundCategory
	}
	products, err := uc.repo.GetByCategory(CategoryID)
	if err != nil {
		return nil, fmt.Errorf("error getting product on category:%w", err)
	}
	return products, nil
}
