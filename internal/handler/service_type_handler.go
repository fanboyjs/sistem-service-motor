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

type ServiceTypeHandler struct {
	service service.ServiceTypeService
}

func NewServiceTypeHandler(service service.ServiceTypeService) *ServiceTypeHandler {
	return &ServiceTypeHandler{service: service}
}

// CreateServiceType godoc
// @Summary Buat tipe servis
// @Description Membuat data tipe servis baru
// @Tags Service Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateServiceTypeRequest true "Data tipe servis"
// @Success 201 {object} map[string]interface{} "Data tipe servis berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 409 {object} map[string]interface{} "nama tipe servis sudah terdaftar"
// @Router /service-types [post]
func (h *ServiceTypeHandler) CreateServiceType(c *gin.Context) {
	var req dto.CreateServiceTypeRequest
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

	serviceType, err := h.service.CreateServiceType(c.Request.Context(), req)
	if err != nil {
		if err == repository.ErrServiceTypeNameExists {
			c.JSON(http.StatusConflict, gin.H{
				"message": "nama tipe servis sudah terdaftar",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal membuat tipe servis",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data tipe servis berhasil dibuat",
		"data":    serviceType,
	})
}

// GetServiceTypes godoc
// @Summary Ambil semua tipe servis
// @Description Mengambil daftar semua tipe servis
// @Tags Service Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Data tipe servis berhasil diambil"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Router /service-types [get]
func (h *ServiceTypeHandler) GetServiceTypes(c *gin.Context) {
	serviceTypes, err := h.service.GetServiceTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data tipe servis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data tipe servis berhasil diambil",
		"data":    serviceTypes,
	})
}

// GetServiceTypeById godoc
// @Summary Ambil tipe servis berdasarkan ID
// @Description Mengambil data tipe servis berdasarkan ID
// @Tags Service Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID tipe servis"
// @Success 200 {object} map[string]interface{} "Data tipe servis berhasil diambil"
// @Failure 400 {object} map[string]interface{} "id tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "tipe servis tidak ditemukan"
// @Router /service-types/{id} [get]
func (h *ServiceTypeHandler) GetServiceTypeById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	serviceType, err := h.service.GetServiceTypeByID(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrServiceTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "tipe servis tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data tipe servis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data tipe servis berhasil diambil, dengan id : %d", id),
		"data":    serviceType,
	})
}

// UpdateServiceType godoc
// @Summary Update tipe servis
// @Description Memperbarui data tipe servis berdasarkan ID
// @Tags Service Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID tipe servis"
// @Param request body dto.UpdateServiceTypeRequest true "Data tipe servis"
// @Success 200 {object} map[string]interface{} "Data tipe servis berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "tipe servis tidak ditemukan"
// @Failure 409 {object} map[string]interface{} "nama tipe servis sudah terdaftar"
// @Router /service-types/{id} [put]
func (h *ServiceTypeHandler) UpdateServiceType(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	var req dto.UpdateServiceTypeRequest
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

	serviceType, err := h.service.UpdateServiceType(c.Request.Context(), id, req)
	if err != nil {
		if err == repository.ErrServiceTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "tipe servis tidak ditemukan",
			})
			return
		}
		if err == repository.ErrServiceTypeNameExists {
			c.JSON(http.StatusConflict, gin.H{
				"message": "nama tipe servis sudah terdaftar",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal update data tipe servis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data tipe servis berhasil diupdate, dengan id : %d", id),
		"data":    serviceType,
	})
}
