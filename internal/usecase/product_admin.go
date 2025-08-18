package usecase

import (
	"errors"
	"fmt"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/Vladislav5108/inventorytracker/internal/domain/repository"
)

type AdminProductUseCase struct {
	productRepo repository.ProductAdminRepository
	catalogRepo repository.ProductRepository
}

func NewAdminProductUseCase(
	productRepo repository.ProductAdminRepository,
	catalogRepo repository.ProductRepository,
) *AdminProductUseCase {
	return &AdminProductUseCase{
		productRepo: productRepo,
		catalogRepo: catalogRepo,
	}
}

func (uc *AdminProductUseCase) Add(product domain.Product) (int, error) {
	if product.Name == "" {
		return 0, domain.ErrName
	}
	if product.Price <= 0 {
		return 0, domain.ErrPrice
	}
	product.IsArchived = false
	return uc.productRepo.Add(product)
}

func (uc *AdminProductUseCase) Update(product domain.Product) error {
	if product.ID <= 0 {
		return domain.ErrInvalidIdProduct
	}
	if product.Name == "" {
		return domain.ErrName
	}
	if product.Price <= 0 {
		return domain.ErrPrice
	}

	existingProduct, err := uc.catalogRepo.GetByID(product.ID)
	if err != nil {
		return fmt.Errorf("producn not found: %w", err)
	}
	if existingProduct.IsArchived {
		return domain.ErrUpdateArchived
	}
	return uc.productRepo.Update(product)
}
func (uc *AdminProductUseCase) Archiv(id int) error {
	if id <= 0 {
		return domain.ErrInvalidIdProduct

	}
	product, err := uc.catalogRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return domain.ErrInvalidIdProduct
		}
		return fmt.Errorf("product not found: %w", err)
	}
	if product.IsArchived {
		return domain.ErrProductArchived
	}
	return uc.productRepo.Archiv(id)
}
func (uc *AdminProductUseCase) Restore(id int) error {
	if id <= 0 {
		return domain.ErrInvalidIdProduct
	}
	product, err := uc.catalogRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return domain.ErrInvalidIdProduct
		}
		return fmt.Errorf("product not found: %w", err)
	}
	if !product.IsArchived {
		return domain.ErrProductArchived
	}
	return uc.productRepo.Restore(id)
}
func (uc *AdminProductUseCase) GetArchived() ([]domain.Product, error) {
	products, err := uc.productRepo.GetArchived()
	if err != nil {
		return nil, fmt.Errorf("error getting archive product: %w", err)
	}
	return products, nil
}
