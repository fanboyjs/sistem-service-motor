package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/middleware"
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

// GetUserInfo godoc
// @Summary Ambil data user saat ini
// @Description Mengambil data user berdasarkan token JWT
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Data user berhasil diambil"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "user tidak ditemukan"
// @Router /user-info [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetInt64(middleware.UserIDKey)

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Data user berhasil diambil",
		"data":    user,
	})
}

// CreateUser godoc
// @Summary Buat user baru
// @Description Membuat data user baru (admin)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "Data user"
// @Success 201 {object} map[string]interface{} "Data user berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Router /users [post]
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

// GetUsers godoc
// @Summary Ambil semua user
// @Description Mengambil daftar semua user (admin)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Data user berhasil diambil"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Router /users [get]
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

// GetUserById godoc
// @Summary Ambil user berdasarkan ID
// @Description Mengambil data user berdasarkan ID (admin)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID user"
// @Success 200 {object} map[string]interface{} "Data user berhasil diambil"
// @Failure 400 {object} map[string]interface{} "id tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "user tidak ditemukan"
// @Router /users/{id} [get]
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

// UpdateUser godoc
// @Summary Update user
// @Description Memperbarui data user berdasarkan ID (admin)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID user"
// @Param request body dto.UpdateUserRequest true "Data user"
// @Success 200 {object} map[string]interface{} "Data user berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "user tidak ditemukan"
// @Router /users/{id} [put]
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

// DeleteUser godoc
// @Summary Hapus user
// @Description Menghapus data user berdasarkan ID (admin)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID user"
// @Success 200 {object} map[string]interface{} "berhasil hapus data user"
// @Failure 400 {object} map[string]interface{} "id tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "user tidak ditemukan"
// @Router /users/{id} [delete]
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
