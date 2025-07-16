package usecase

import (
	"errors"
	"fmt"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/Vladislav5108/inventorytracker/internal/domain/repository"
)

type AdminProductUseCase struct {
	repo repository.ProductAdminRepository
}

func NewAdminProductUseCase(repo repository.ProductAdminRepository) *AdminProductUseCase {
	return &AdminProductUseCase{repo: repo}
}

func (uc *AdminProductUseCase) Add(product domain.Product) (int, error) {
	if product.Name == "" {
		return 0, domain.ErrName
	}
	if product.Price <= 0 {
		return 0, domain.ErrPrice
	}
	return uc.repo.Add(product)
}

func (uc *AdminProductUseCase) UpDate(product domain.Product) error {
	if product.ID <= 0 {
		return domain.ErrInvalidIdProduct
	}
	if product.Name == "" {
		return domain.ErrName
	}
	if product.Price <= 0 {
		return domain.ErrPrice
	}
	return uc.repo.UpDate(product)
}
func (uc *AdminProductUseCase) Archiv(id int) error {
	if id <= 0 {
		return domain.ErrInvalidIdProduct
	}
	_, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return domain.ErrInvalidIdProduct
		}
		return fmt.Errorf("продукт не найден: %w", err)
	}
	return uc.repo.Archiv(id)
}
func (uc *AdminProductUseCase) Restore(id int) error {
	if id <= 0 {
		return domain.ErrInvalidIdProduct
	}
	_, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return domain.ErrInvalidIdProduct
		}
		return fmt.Errorf("продукт не найден: %w", err)
	}
	return uc.repo.Restore(id)
}
func (uc *AdminProductUseCase) GetArchived() ([]domain.Product, error) {
	products, err := uc.repo.GetArchived()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка азхивных товаров: %w", err)
	}
	return products, nil
}

