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
