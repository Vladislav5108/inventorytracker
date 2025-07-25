package dto

import "time"

type CreateProductRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Price    int    `json:"price" binding:"required,gt=0"`
	Quantity int    `json:"quantity" binding:"gte=0"`
}

type ProductResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Price      int       `json:"price"`
	Quantity   int       `json:"quantity"`
	CategoryID int       `json:"category_id,omitempty"`
	CreatedAt  time.Time `json:"created_at" format:"date-time"`
	IsArchived bool      `json:"is_archived"`
}
