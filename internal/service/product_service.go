package service

import (
	"ta/internal/domain"
	"ta/internal/repository"

	"github.com/google/uuid"
)

type ProductService struct {
	repo repository.ProductRepositoryInterface
}

func NewProductService(repo repository.ProductRepositoryInterface) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(product *domain.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) GetByID(id string) (*domain.Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(productID)
}

func (s *ProductService) CheckAvailability(productID uuid.UUID, quantity int) (bool, error) {
	product, err := s.repo.GetByID(productID)
	if err != nil {
		return false, err
	}
	return product.Quantity >= quantity, nil
}

func (s *ProductService) UpdateQuantity(productID uuid.UUID, quantity int) error {
	return s.repo.UpdateQuantity(productID, quantity)
}

