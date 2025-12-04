package service

import (
	"ta/internal/domain"
	"ta/internal/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(id uuid.UUID) (*domain.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByUserID(userID uuid.UUID) ([]*domain.Order, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

var _ repository.OrderRepositoryInterface = (*MockOrderRepository)(nil)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByIDs(ids []uuid.UUID) ([]*domain.Product, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockProductRepository) GetByID(id uuid.UUID) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateQuantity(id uuid.UUID, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

var _ repository.ProductRepositoryInterface = (*MockProductRepository)(nil)


func TestOrderService_CreateOrder(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()

	tests := []struct {
		name    string
		items   []domain.OrderItem
		setup   func(*MockOrderRepository, *MockProductRepository, *ProductService)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful order creation",
			items: []domain.OrderItem{
				{ProductID: productID, Quantity: 2},
			},
			setup: func(orderRepo *MockOrderRepository, productRepo *MockProductRepository, productSvc *ProductService) {
				products := []*domain.Product{
					{ID: productID, Quantity: 10},
				}
				productRepo.On("GetByIDs", mock.Anything).Return(products, nil)
				productRepo.On("UpdateQuantity", productID, 8).Return(nil)
				orderRepo.On("Create", mock.AnythingOfType("*domain.Order")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "insufficient product quantity",
			items: []domain.OrderItem{
				{ProductID: productID, Quantity: 15},
			},
			setup: func(orderRepo *MockOrderRepository, productRepo *MockProductRepository, productSvc *ProductService) {
				products := []*domain.Product{
					{ID: productID, Quantity: 10},
				}
				productRepo.On("GetByIDs", mock.Anything).Return(products, nil)
			},
			wantErr: true,
			errMsg:  "insufficient product quantity",
		},
		{
			name: "product not found",
			items: []domain.OrderItem{
				{ProductID: productID, Quantity: 2},
			},
			setup: func(orderRepo *MockOrderRepository, productRepo *MockProductRepository, productSvc *ProductService) {
				productRepo.On("GetByIDs", mock.Anything).Return([]*domain.Product{}, nil)
			},
			wantErr: true,
			errMsg:  "product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo := new(MockOrderRepository)
			mockProductRepo := new(MockProductRepository)
			mockProductSvc := NewProductService(mockProductRepo)
			tt.setup(mockOrderRepo, mockProductRepo, mockProductSvc)

			service := NewOrderService(mockOrderRepo, mockProductRepo, mockProductSvc)
			_, err := service.CreateOrder(userID, tt.items)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				mockOrderRepo.AssertExpectations(t)
				mockProductRepo.AssertExpectations(t)
			}
		})
	}
}

