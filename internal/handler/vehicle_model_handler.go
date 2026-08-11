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

type VehicleModelHandler struct {
	service service.VehicleModelService
}

func NewVehicleModelHandler(service service.VehicleModelService) *VehicleModelHandler {
	return &VehicleModelHandler{service: service}
}

func (h *VehicleModelHandler) CreateVehicleModel(c *gin.Context) {
	var req dto.CreateVehicleModelRequest
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

	modelData, err := h.service.CreateVehicleModel(c.Request.Context(), req)
	if err != nil {
		if err == repository.ErrModelNameExists {
			c.JSON(http.StatusConflict, gin.H{
				"message": "nama model kendaraan sudah terdaftar",
			})
			return
		}
		if err == repository.ErrVehicleBrandNotFoundFK {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "brand kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal membuat model kendaraan",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data model kendaraan berhasil dibuat",
		"data":    modelData,
	})
}

func (h *VehicleModelHandler) GetVehicleModels(c *gin.Context) {
	models, err := h.service.GetVehicleModels(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data model kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data model kendaraan berhasil diambil",
		"data":    models,
	})
}

func (h *VehicleModelHandler) GetVehicleModelById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	modelData, err := h.service.GetVehicleModelByID(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrVehicleModelNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "model kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data model kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data model kendaraan berhasil diambil, dengan id : %d", id),
		"data":    modelData,
	})
}

func (h *VehicleModelHandler) UpdateVehicleModel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	var req dto.UpdateVehicleModelRequest
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

	modelData, err := h.service.UpdateVehicleModel(c.Request.Context(), id, req)
	if err != nil {
		if err == repository.ErrVehicleModelNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "model kendaraan tidak ditemukan",
			})
			return
		}
		if err == repository.ErrModelNameExists {
			c.JSON(http.StatusConflict, gin.H{
				"message": "nama model kendaraan sudah terdaftar",
			})
			return
		}
		if err == repository.ErrVehicleBrandNotFoundFK {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "brand kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal update data model kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data model kendaraan berhasil diupdate, dengan id : %d", id),
		"data":    modelData,
	})
}
