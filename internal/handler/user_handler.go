package handler

import (
	"net/http"
	"ta/internal/dto"
	"ta/internal/mapper"
	"ta/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := mapper.ToUserDomain(req)
	if err := h.service.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapper.ToUserResponse(user))
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, mapper.ToUserResponse(user))
}

