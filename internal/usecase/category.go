package usecase

import (
	"fmt"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/Vladislav5108/inventorytracker/internal/domain/repository"
)

type CategoryUseCase struct {
	repo repository.CategoryRepository
}

func NewCategoryUseCase(repo repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{repo: repo}
}
func (uc *CategoryUseCase) CreateCategory(category domain.Category) (int, error) {
	if category.Name == "" {
		return 0, domain.ErrNameCategory
	}
	if category.ID != 0 {
		return 0, fmt.Errorf("the id must not be specified during creation ")
	}
	return uc.repo.CreateCategory(category)
}
func (uc *CategoryUseCase) GetByIDCategory(id int) (domain.Category, error) {
	if id <= 0 {
		return domain.Category{}, domain.ErrNotFoundCategory
	}
	category, err := uc.repo.GetByIDCategory(id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("repository error: %w", err)
	}
	return category, nil
}
func (uc *CategoryUseCase) GetAllCategories() ([]domain.Category, error) {
	categories, err := uc.repo.GetAllCategories()
	if err != nil {
		return categories, fmt.Errorf("error getting all categories: %w", err)
	}
	return categories, nil
}
func (uc *CategoryUseCase) UpDateCategory(category domain.Category) error {
	if category.ID <= 0 {
		return domain.ErrNotFoundCategory
	}
	if category.Name == "" {
		return domain.ErrNameCategory
	}
	err := uc.repo.UpdateCategory(category)
	if err != nil {
		return fmt.Errorf("category update error: %w", err)
	}
	return nil
}
func (uc *CategoryUseCase) DeleteCategory(id int) error {
	if id <= 0 {
		return domain.ErrNotFoundCategory
	}
	err := uc.repo.DeleteCategory(id)
	if err != nil {
		return fmt.Errorf("error delete category: %w", err)
	}
	return nil
}
