package domain

import (
	"errors"
	"time"
)

var (
	ErrProductArchived  = errors.New("product on archive")
	ErrUpdateArchived   = errors.New("not update archive product")
	ErrDuplicateName    = errors.New("товар с таким именем уже существует")
	ErrInvalidIdProduct = errors.New("id должен быть положительным")
	ErrProductNotFound  = errors.New("товара нет ")
	ErrName             = errors.New("имя не должно быть пустым")
	ErrPrice            = errors.New("цена должна быть положительной")
	ErrQuantity         = errors.New("количество не должно быть отрицательным")
)

type Product struct {
	ID         int
	Name       string
	Price      int
	Quantity   int
	CategoryID int
	CreatedAt  time.Time
	IsArchived bool
}

func (p *Product) Validate() error {

	if p.Name == "" {
		return ErrName
	}
	if p.Price <= 0 {
		return ErrPrice
	}
	if p.Quantity < 0 {
		return ErrQuantity
	}
	return nil
}
