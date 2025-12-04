package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Items     []OrderItem
	CreatedAt time.Time
}

type OrderItem struct {
	ID           uuid.UUID
	OrderID      uuid.UUID
	ProductID    uuid.UUID
	Quantity     int
	PriceAtTime  float64
	CreatedAt    time.Time
}

