package handlers

import (
	"net/http"

	"github.com/Secure-Website-Builder/Backend/internal/services/auth"
	"github.com/Secure-Website-Builder/Backend/internal/types"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

type RegisterRequest struct {
	Name     string         `json:"name" binding:"required"`
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=6"`
	Role     string         `json:"role" binding:"required,oneof=store_owner customer"`
	StoreID  *int64         `json:"store_id"`             // required only for customers
	Phone    *string        `json:"phone,omitempty"`      // optional
	Address  *types.Address `json:"address,omitempty"`    // optional
}

type LoginRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=store_owner customer"`
	StoreID  *int64 `json:"store_id"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var phone *string
	if req.Phone != nil {
		phone = req.Phone
	}

	var addr *types.Address
	if req.Address != nil {
		addr = req.Address
	}

	token, err := h.service.Register(
		c.Request.Context(),
		req.Name,
		req.Email,
		req.Password,
		req.Role,
		req.StoreID,
		phone,
		addr,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}


func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(
		c.Request.Context(),
		req.Email,
		req.Password,
		req.Role,
		req.StoreID,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}
