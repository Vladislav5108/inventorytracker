package repository

import domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"

type CategoryRepository interface {
	CreateCategory(category domain.Category) (int, error) 
	GetByIDCategory(id int) (domain.Category, error)     
	GetAllCategories() ([]domain.Category, error)      
	UpdateCategory(category domain.Category) error        
	DeleteCategory(id int) error                          
}
