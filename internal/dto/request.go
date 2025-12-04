package dto

import "github.com/google/uuid"

type RegisterUserRequest struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Age       int    `json:"age" binding:"required,min=18"`
	IsMarried bool   `json:"is_married"`
	Password  string `json:"password" binding:"required,min=8"`
}

type CreateProductRequest struct {
	Description string   `json:"description" binding:"required"`
	Tags        []string `json:"tags"`
	Quantity    int      `json:"quantity" binding:"required,min=0"`
}

type CreateOrderRequest struct {
	UserID uuid.UUID           `json:"user_id" binding:"required"`
	Items  []OrderItemRequest  `json:"items" binding:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

