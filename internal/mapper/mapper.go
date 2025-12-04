package mapper

import (
	"ta/internal/domain"
	"ta/internal/dto"
	"time"

	"github.com/google/uuid"
)

func ToUserDomain(req dto.RegisterUserRequest) *domain.User {
	return &domain.User{
		ID:        uuid.New(),
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Fullname:  req.Firstname + " " + req.Lastname,
		Age:       req.Age,
		IsMarried: req.IsMarried,
		Password:  req.Password,
	}
}

func ToUserResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Fullname:  user.Fullname,
		Age:       user.Age,
		IsMarried: user.IsMarried,
	}
}

func ToProductDomain(req dto.CreateProductRequest) *domain.Product {
	return &domain.Product{
		ID:          uuid.New(),
		Description: req.Description,
		Tags:        req.Tags,
		Quantity:    req.Quantity,
	}
}

func ToProductResponse(product *domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          product.ID,
		Description: product.Description,
		Tags:        product.Tags,
		Quantity:    product.Quantity,
	}
}

func ToOrderResponse(order *domain.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.OrderItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			PriceAtTime: item.PriceAtTime,
		}
	}
	return dto.OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Items:     items,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
	}
}

