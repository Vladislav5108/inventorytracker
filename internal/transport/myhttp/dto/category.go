package dto

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" binding:"required,min=2"`
	Description *string `json:"description,omitempty"`
}
