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

type ServiceRecordHandler struct {
	service service.ServiceRecordService
}

func NewServiceRecordHandler(service service.ServiceRecordService) *ServiceRecordHandler {
	return &ServiceRecordHandler{service: service}
}

func (h *ServiceRecordHandler) CreateServiceRecord(c *gin.Context) {
	var req dto.CreateServiceRecordRequest
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

	userID := c.GetInt64(middleware.UserIDKey)
	record, err := h.service.CreateServiceRecord(c.Request.Context(), userID, req)
	if err != nil {
		if err == repository.ErrVehicleNotFoundFK {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "kendaraan tidak ditemukan",
			})
			return
		}
		if err == repository.ErrServiceTypeNotFoundFK {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "jenis servis tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal membuat data record servis",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data record servis berhasil dibuat",
		"data":    record,
	})
}

func (h *ServiceRecordHandler) GetServiceRecords(c *gin.Context) {
	userID := c.GetInt64(middleware.UserIDKey)
	records, err := h.service.GetServiceRecords(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data record servis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data record servis berhasil diambil",
		"data":    records,
	})
}

func (h *ServiceRecordHandler) GetServiceRecordById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	userID := c.GetInt64(middleware.UserIDKey)
	record, err := h.service.GetServiceRecordByID(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrServiceRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "record servis tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data record servis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data record servis berhasil diambil, dengan id : %d", id),
		"data":    record,
	})
}

func (h *ServiceRecordHandler) UpdateServiceRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	var req dto.UpdateServiceRecordRequest
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

	userID := c.GetInt64(middleware.UserIDKey)
	record, err := h.service.UpdateServiceRecord(c.Request.Context(), id, userID, req)
	if err != nil {
		if err == repository.ErrServiceRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "record servis tidak ditemukan",
			})
			return
		}
		if err == repository.ErrVehicleNotFoundFK {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "kendaraan tidak ditemukan",
			})
			return
		}
		if err == repository.ErrServiceTypeNotFoundFK {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "jenis servis tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal update data record servis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data record servis berhasil diupdate, dengan id : %d", id),
		"data":    record,
	})
}

func (h *ServiceRecordHandler) DeleteServiceRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	userID := c.GetInt64(middleware.UserIDKey)
	err = h.service.DeleteServiceRecord(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrServiceRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "record servis tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal hapus data record servis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("berhasil hapus data record servis, dengan id : %d", id),
	})
}
