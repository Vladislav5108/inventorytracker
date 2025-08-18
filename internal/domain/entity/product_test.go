package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductValidate_Success(t *testing.T) {
	p := Product{
		Name:     "iPhpne 15 pro",
		Price:    300,
		Quantity: 3,
	}
	err := p.Validate()
	assert.NoError(t, err)
}
func TestProductValidate_EmptyName(t *testing.T) {
	p := Product{
		Name:     "",
		Price:    300,
		Quantity: 3,
	}
	err := p.Validate()
	assert.ErrorIs(t, err, ErrName)
}
func TestProductValidate_InvalidPrice(t *testing.T) {
	p := Product{
		Name:     "iPhpne 15 pro",
		Price:    -100,
		Quantity: 3,
	}
	err := p.Validate()
	assert.ErrorIs(t, err, ErrPrice)

}
func TestProductValidate_ZeroPrice(t *testing.T) {
	p := Product{
		Name:     "iphone 15 pro",
		Price:    0,
		Quantity: 3,
	}
	err := p.Validate()

	if err != ErrPrice {
		t.Errorf("Expected error ErrPrice, got %v", err)
	}
}
func TestProductValidate_InvalidQuantity(t *testing.T) {
	p := Product{
		Name:     "iPyjne 15p pro",
		Price:    300,
		Quantity: -3,
	}
	err := p.Validate()
	assert.ErrorIs(t, err, ErrQuantity)
}
