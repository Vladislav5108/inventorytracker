package domain

import (
	"errors"
	"time"
)

var (
	ErrProductArchived  = errors.New("product on archive")
	ErrUpdateArchived   = errors.New("not update archive product")
	ErrDuplicateName    = errors.New("product with this name already exists")
	ErrInvalidIdProduct = errors.New("ID must be positive")
	ErrProductNotFound  = errors.New("product not found")
	ErrName             = errors.New("name cannot be empty")
	ErrPrice            = errors.New("price must be positive")
	ErrQuantity         = errors.New("quantity cannot be negative")
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
