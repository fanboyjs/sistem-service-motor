package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/service"
	"example.com/my-api/pkg/validator"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request body tidak valid",
		})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, email, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal register",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "register berhasil",
		"data": gin.H{
			"token": token,
			"email": email,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request body tidak valid",
		})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal login",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login berhasil",
		"data": gin.H{
			"token": token,
		},
	})
}
