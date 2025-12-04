package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Fullname  string    `json:"fullname"`
	Age       int       `json:"age"`
	IsMarried bool      `json:"is_married"`
}

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Quantity    int       `json:"quantity"`
}

type OrderItemResponse struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int       `json:"quantity"`
	PriceAtTime float64   `json:"price_at_time"`
}

type OrderResponse struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"user_id"`
	Items     []OrderItemResponse `json:"items"`
	CreatedAt string             `json:"created_at"`
}

