package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/repository"
	"example.com/my-api/internal/service"
	"example.com/my-api/pkg/validator"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
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

	user, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal membuat user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
	    "message": "Data user berhasil dibuat",
	    "data":    user,
	})

}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
	    "message": "Data user berhasil diambil",
	    "data":    users,
	})
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "user tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
	    "message": fmt.Sprintf("Data user berhasil diambil, dengan id : %d", id),
	    "data":    user,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	var req dto.UpdateUserRequest
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

	user, err := h.service.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "user tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal update data user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data user berhasil diupdate, dengan id : %d", id),
		"data":    user,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	err = h.service.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "user tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal hapus data user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("berhasil hapus data user, dengan id : %d", id),
	})
}
