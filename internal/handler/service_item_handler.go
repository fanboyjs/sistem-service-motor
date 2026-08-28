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

type ServiceItemHandler struct {
	service service.ServiceItemService
}

func NewServiceItemHandler(service service.ServiceItemService) *ServiceItemHandler {
	return &ServiceItemHandler{service: service}
}

// CreateServiceItem godoc
// @Summary Buat item servis
// @Description Membuat data item servis baru
// @Tags Service Item
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateServiceItemRequest true "Data item servis"
// @Success 201 {object} map[string]interface{} "Data item servis berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "record servis tidak ditemukan"
// @Router /service-items [post]
func (h *ServiceItemHandler) CreateServiceItem(c *gin.Context) {
	var req dto.CreateServiceItemRequest
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

	serviceItem, err := h.service.CreateServiceItem(c.Request.Context(), req)
	if err != nil {
		if err == repository.ErrServiceRecordNotFoundFK {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "record servis tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal membuat item servis",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data item servis berhasil dibuat",
		"data":    serviceItem,
	})
}

// GetServiceItems godoc
// @Summary Ambil semua item servis
// @Description Mengambil daftar semua item servis
// @Tags Service Item
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Data item servis berhasil diambil"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Router /service-items [get]
func (h *ServiceItemHandler) GetServiceItems(c *gin.Context) {
	serviceItems, err := h.service.GetServiceItems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data item servis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data item servis berhasil diambil",
		"data":    serviceItems,
	})
}

// GetServiceItemById godoc
// @Summary Ambil item servis berdasarkan ID
// @Description Mengambil data item servis berdasarkan ID
// @Tags Service Item
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID item servis"
// @Success 200 {object} map[string]interface{} "Data item servis berhasil diambil"
// @Failure 400 {object} map[string]interface{} "id tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "item servis tidak ditemukan"
// @Router /service-items/{id} [get]
func (h *ServiceItemHandler) GetServiceItemById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	serviceItem, err := h.service.GetServiceItemByID(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrServiceItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item servis tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data item servis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data item servis berhasil diambil, dengan id : %d", id),
		"data":    serviceItem,
	})
}

// UpdateServiceItem godoc
// @Summary Update item servis
// @Description Memperbarui data item servis berdasarkan ID
// @Tags Service Item
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID item servis"
// @Param request body dto.UpdateServiceItemRequest true "Data item servis"
// @Success 200 {object} map[string]interface{} "Data item servis berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "item servis tidak ditemukan"
// @Router /service-items/{id} [put]
func (h *ServiceItemHandler) UpdateServiceItem(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	var req dto.UpdateServiceItemRequest
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

	serviceItem, err := h.service.UpdateServiceItem(c.Request.Context(), id, req)
	if err != nil {
		if err == repository.ErrServiceItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item servis tidak ditemukan",
			})
			return
		}
		if err == repository.ErrServiceRecordNotFoundFK {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "record servis tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal update data item servis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data item servis berhasil diupdate, dengan id : %d", id),
		"data":    serviceItem,
	})
}
