package domain

import "errors"

var ErrNameCategory = errors.New("имя категории не должно быть пустым")
var ErrNotId = errors.New("ID не существует")

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

