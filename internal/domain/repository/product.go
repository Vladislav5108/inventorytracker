package repository

import domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"

type ProductRepository interface {
	GetByID(id int) (domain.Product, error)                 
	GetAll() ([]domain.Product, error)                      
	GetByCategory(CategoryID int) ([]domain.Product, error) 
}

type ProductAdminRepository interface {
	Add(product domain.Product) (int, error) 
	Update(product domain.Product) error     
	Archiv(id int) error                     
	Restore(id int) error                    
	GetArchived() ([]domain.Product, error)  
}
