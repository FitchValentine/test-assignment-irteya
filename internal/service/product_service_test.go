package service

import (
	"ta/internal/domain"
	"ta/internal/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepositoryForService struct {
	mock.Mock
}

func (m *MockProductRepositoryForService) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepositoryForService) GetByID(id uuid.UUID) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepositoryForService) UpdateQuantity(id uuid.UUID, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

func (m *MockProductRepositoryForService) GetByIDs(ids []uuid.UUID) ([]*domain.Product, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

var _ repository.ProductRepositoryInterface = (*MockProductRepositoryForService)(nil)

func TestProductService_CheckAvailability(t *testing.T) {
	productID := uuid.New()

	tests := []struct {
		name       string
		productID  uuid.UUID
		quantity   int
		setup      func(*MockProductRepositoryForService)
		wantResult bool
		wantErr    bool
	}{
		{
			name:      "product available",
			productID: productID,
			quantity:  5,
			setup: func(m *MockProductRepositoryForService) {
				m.On("GetByID", productID).Return(&domain.Product{ID: productID, Quantity: 10}, nil)
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name:      "product insufficient quantity",
			productID: productID,
			quantity:  15,
			setup: func(m *MockProductRepositoryForService) {
				m.On("GetByID", productID).Return(&domain.Product{ID: productID, Quantity: 10}, nil)
			},
			wantResult: false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProductRepositoryForService)
			tt.setup(mockRepo)
			service := NewProductService(mockRepo)

			result, err := service.CheckAvailability(tt.productID, tt.quantity)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, result)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

