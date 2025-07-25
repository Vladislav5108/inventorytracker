package domain

import "errors"

var ErrNameCategory = errors.New("category name cannot be empty")
var ErrNotFoundCategory = errors.New("ID not found")
var ErrAlreadyCategory = errors.New("category already exists")

type Category struct {
	ID          int
	Name        string
	Description string
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return ErrNameCategory
	}
	return nil
}
