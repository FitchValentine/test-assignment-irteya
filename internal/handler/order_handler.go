package handler

import (
	"net/http"
	"ta/internal/domain"
	"ta/internal/dto"
	"ta/internal/mapper"
	"ta/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items := make([]domain.OrderItem, len(req.Items))
	for i, itemReq := range req.Items {
		items[i] = domain.OrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
		}
	}

	order, err := h.service.CreateOrder(req.UserID, items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapper.ToOrderResponse(order))
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, mapper.ToOrderResponse(order))
}

func (h *OrderHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	orders, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = mapper.ToOrderResponse(order)
	}
	c.JSON(http.StatusOK, responses)
}

