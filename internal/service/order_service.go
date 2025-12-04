package service

import (
	"errors"
	"ta/internal/domain"
	"ta/internal/repository"
	"time"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo   repository.OrderRepositoryInterface
	productRepo repository.ProductRepositoryInterface
	productSvc  *ProductService
}

func NewOrderService(orderRepo repository.OrderRepositoryInterface, productRepo repository.ProductRepositoryInterface, productSvc *ProductService) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		productSvc:  productSvc,
	}
}

func (s *OrderService) CreateOrder(userID uuid.UUID, items []domain.OrderItem) (*domain.Order, error) {
	productIDs := make([]uuid.UUID, len(items))
	for i, item := range items {
		productIDs[i] = item.ProductID
	}

	products, err := s.productRepo.GetByIDs(productIDs)
	if err != nil {
		return nil, err
	}

	productMap := make(map[uuid.UUID]*domain.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	for _, item := range items {
		product, exists := productMap[item.ProductID]
		if !exists {
			return nil, errors.New("product not found")
		}
		if product.Quantity < item.Quantity {
			return nil, errors.New("insufficient product quantity")
		}
	}

	order := &domain.Order{
		ID:        uuid.New(),
		UserID:    userID,
		Items:     make([]domain.OrderItem, len(items)),
		CreatedAt: time.Now(),
	}

	for i, item := range items {
		product := productMap[item.ProductID]
		order.Items[i] = domain.OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			PriceAtTime: 0.0,
			CreatedAt:   time.Now(),
		}

		newQuantity := product.Quantity - item.Quantity
		if err := s.productSvc.UpdateQuantity(product.ID, newQuantity); err != nil {
			return nil, err
		}
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetByID(id string) (*domain.Order, error) {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.orderRepo.GetByID(orderID)
}

func (s *OrderService) GetByUserID(userID string) ([]*domain.Order, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.orderRepo.GetByUserID(uid)
}

