package dto

import "time"

type CreateProductRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type ProductResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Quantity   int       `json:"quantity"`
	CategoryID int       `json:"category_id,omitempty"`
	CreatedAt  time.Time `json:"created_at" format:"date-time"`
}


