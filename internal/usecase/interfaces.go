package usecase

import domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"

type ProductCatalog interface {
	GetByID(id int) (domain.Product, error)                 // возвращает товар по ID
	GetAll() ([]domain.Product, error)                      // возвращает список всех товаров
	GetByCategory(CategoryID int) ([]domain.Product, error) // возращает товары конкретной категории
}

type ProductAdmin interface {
	Add(product domain.Product) (int, error) // добавляем новый товар
	Update(product domain.Product) error     // обновляем данные товара
	Archiv(id int) error                     //убераем товар из католога но не из базы
	Restore(id int) error                    // возвращаем товар после архивации
	GetArchived() ([]domain.Product, error)  //возвращаем список скрытых товаров
}

type CategoryCatalog interface {
	CreateCategory(category domain.Category) (int, error) //Создание новой категории
	GetByIDCategory(id int) (domain.Category, error)      //получение категории по ID
	GetAllCategories() ([]domain.Category, error)         //получение списка всех категорий
	UpDateCategory(category domain.Category) error        //обновление категории
	DeleteCategory(id int) error                          //удаление кактегории
}
