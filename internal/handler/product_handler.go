package handler

import (
	"net/http"
	"ta/internal/dto"
	"ta/internal/mapper"
	"ta/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := mapper.ToProductDomain(req)
	if err := h.service.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapper.ToProductResponse(product))
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, mapper.ToProductResponse(product))
}

