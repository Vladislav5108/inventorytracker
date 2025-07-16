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
		return 0, fmt.Errorf("id не должен указываться при создании")
	}
	return uc.repo.CreateCategory(category)
}
func (uc *CategoryUseCase) GetByIDCategory(id int) (domain.Category, error) {
	if id <= 0 {
		return domain.Category{}, domain.ErrNotId
	}
	category, err := uc.repo.GetByIDCategory(id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("ошибка репозитория: %w", err)
	}
	return category, nil
}
func (uc *CategoryUseCase) GetALLCategory() ([]domain.Category, error) {
	categories, err := uc.repo.GetALLCategory()
	if err != nil {
		return categories, fmt.Errorf("ошибка получения всех категорий: %w", err)
	}
	return categories, nil
}
func (uc *CategoryUseCase) UpDateCategory(category domain.Category) error {
	if category.ID <= 0 {
		return domain.ErrNotId
	}
	if category.Name == "" {
		return domain.ErrNameCategory
	}
	err := uc.repo.UpDateCategory(category)
	if err != nil {
		return fmt.Errorf("ошибка обновления категории: %w", err)
	}
	return nil
}
func (uc *CategoryUseCase) DeleteCategory(id int) error {
	if id <= 0 {
		return domain.ErrNotId
	}
	err := uc.repo.DeleteCategory(id)
	if err != nil {
		return fmt.Errorf("ошибка удаления: %w", err)
	}
	return nil
}

