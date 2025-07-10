package repository

import "github.com/Vladislav5108/inventorytracker/internal/domain"

type ProductRepository interface {
	GetByID(id int) (domain.Product, error)                 // возвращает товар по ID
	GetALL() ([]domain.Product, error)                      // возвращает список всех товаров
	GetByCategory(CategoryID int) ([]domain.Product, error) // возращает товары конкретной категории
	SearchByName(keyword string) ([]domain.Product, error)  // возвращает товары по части названия без учета регистра
}

type ProductAdminRepository interface {
	ProductRepository

	Add(product domain.Product) error       // добавляем новый товар
	UpDate(product domain.Product) error    // обновляем данные товара
	Archiv(id int) error                    //убераем товар из католога но не из базы
	Restore(id int) error                   // возвращаем товар после архивации
	GetArchived() ([]domain.Product, error) //возвращаем список скрытых товаров
}
