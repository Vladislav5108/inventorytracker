package domain

import "errors"

var ErrNameCategory = errors.New("имя категории не должно быть пустым")
var ErrNotFoundCategory = errors.New("ID не существует")
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
