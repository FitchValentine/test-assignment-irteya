package repository

import (
	"ta/internal/domain"

	"github.com/google/uuid"
)

type UserRepositoryInterface interface {
	Create(user *domain.User) error
	GetByID(id uuid.UUID) (*domain.User, error)
}

type ProductRepositoryInterface interface {
	Create(product *domain.Product) error
	GetByID(id uuid.UUID) (*domain.Product, error)
	UpdateQuantity(id uuid.UUID, quantity int) error
	GetByIDs(ids []uuid.UUID) ([]*domain.Product, error)
}

type OrderRepositoryInterface interface {
	Create(order *domain.Order) error
	GetByID(id uuid.UUID) (*domain.Order, error)
	GetByUserID(userID uuid.UUID) ([]*domain.Order, error)
}

