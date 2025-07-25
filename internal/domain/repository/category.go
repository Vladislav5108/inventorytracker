package repository

import domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"

type CategoryRepository interface {
	CreateCategory(category domain.Category) (int, error) //Создание новой категории
	GetByIDCategory(id int) (domain.Category, error)      //получение категории по ID
	GetAllCategories() ([]domain.Category, error)         //получение списка всех категорий
	UpdateCategory(category domain.Category) error        //обновление категории
	DeleteCategory(id int) error                          //удаление кактегории
}
