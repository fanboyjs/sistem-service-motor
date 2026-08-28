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

// CreateVehicleModel godoc
// @Summary Buat model kendaraan
// @Description Membuat data model kendaraan baru yang terhubung ke brand
// @Tags Vehicle Model
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateVehicleModelRequest true "Data model kendaraan"
// @Success 201 {object} map[string]interface{} "Data model kendaraan berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "brand kendaraan tidak ditemukan"
// @Failure 409 {object} map[string]interface{} "nama model kendaraan sudah terdaftar"
// @Router /vehicle-models [post]
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

// GetVehicleModels godoc
// @Summary Ambil semua model kendaraan
// @Description Mengambil daftar semua model kendaraan
// @Tags Vehicle Model
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Data model kendaraan berhasil diambil"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Router /vehicle-models [get]
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

// GetVehicleModelById godoc
// @Summary Ambil model kendaraan berdasarkan ID
// @Description Mengambil data model kendaraan berdasarkan ID
// @Tags Vehicle Model
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID model"
// @Success 200 {object} map[string]interface{} "Data model kendaraan berhasil diambil"
// @Failure 400 {object} map[string]interface{} "id tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "model kendaraan tidak ditemukan"
// @Router /vehicle-models/{id} [get]
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

// UpdateVehicleModel godoc
// @Summary Update model kendaraan
// @Description Memperbarui data model kendaraan berdasarkan ID
// @Tags Vehicle Model
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID model"
// @Param request body dto.UpdateVehicleModelRequest true "Data model kendaraan"
// @Success 200 {object} map[string]interface{} "Data model kendaraan berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "model kendaraan tidak ditemukan"
// @Failure 409 {object} map[string]interface{} "nama model kendaraan sudah terdaftar"
// @Router /vehicle-models/{id} [put]
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
