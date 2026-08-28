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

type VehicleTaxHandler struct {
	service service.VehicleTaxService
}

func NewVehicleTaxHandler(service service.VehicleTaxService) *VehicleTaxHandler {
	return &VehicleTaxHandler{service: service}
}

// CreateVehicleTax godoc
// @Summary Buat pajak kendaraan
// @Description Membuat data pajak kendaraan baru untuk user yang login
// @Tags Vehicle Tax
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateVehicleTaxRequest true "Data pajak kendaraan"
// @Success 201 {object} map[string]interface{} "Data pajak kendaraan berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "kendaraan tidak ditemukan"
// @Router /vehicle-taxes [post]
func (h *VehicleTaxHandler) CreateVehicleTax(c *gin.Context) {
	var req dto.CreateVehicleTaxRequest
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
	tax, err := h.service.CreateVehicleTax(c.Request.Context(), userID, req)
	if err != nil {
		if err == repository.ErrVehicleNotFoundFK {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal membuat data pajak kendaraan",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data pajak kendaraan berhasil dibuat",
		"data":    tax,
	})
}

// GetVehicleTaxes godoc
// @Summary Ambil semua pajak kendaraan user
// @Description Mengambil daftar pajak kendaraan milik user yang login
// @Tags Vehicle Tax
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Data pajak kendaraan berhasil diambil"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Router /vehicle-taxes [get]
func (h *VehicleTaxHandler) GetVehicleTaxes(c *gin.Context) {
	userID := c.GetInt64(middleware.UserIDKey)
	taxes, err := h.service.GetVehicleTaxes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data pajak kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data pajak kendaraan berhasil diambil",
		"data":    taxes,
	})
}

// GetVehicleTaxById godoc
// @Summary Ambil pajak kendaraan berdasarkan ID
// @Description Mengambil data pajak kendaraan milik user yang login berdasarkan ID
// @Tags Vehicle Tax
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID pajak"
// @Success 200 {object} map[string]interface{} "Data pajak kendaraan berhasil diambil"
// @Failure 400 {object} map[string]interface{} "id tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "pajak kendaraan tidak ditemukan"
// @Router /vehicle-taxes/{id} [get]
func (h *VehicleTaxHandler) GetVehicleTaxById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	userID := c.GetInt64(middleware.UserIDKey)
	tax, err := h.service.GetVehicleTaxByID(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrVehicleTaxNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "pajak kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data pajak kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data pajak kendaraan berhasil diambil, dengan id : %d", id),
		"data":    tax,
	})
}

// UpdateVehicleTax godoc
// @Summary Update pajak kendaraan
// @Description Memperbarui data pajak kendaraan milik user yang login berdasarkan ID
// @Tags Vehicle Tax
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID pajak"
// @Param request body dto.UpdateVehicleTaxRequest true "Data pajak kendaraan"
// @Success 200 {object} map[string]interface{} "Data pajak kendaraan berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "pajak kendaraan tidak ditemukan"
// @Router /vehicle-taxes/{id} [put]
func (h *VehicleTaxHandler) UpdateVehicleTax(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	var req dto.UpdateVehicleTaxRequest
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
	tax, err := h.service.UpdateVehicleTax(c.Request.Context(), id, userID, req)
	if err != nil {
		if err == repository.ErrVehicleTaxNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "pajak kendaraan tidak ditemukan",
			})
			return
		}
		if err == repository.ErrVehicleNotFoundFK {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal update data pajak kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data pajak kendaraan berhasil diupdate, dengan id : %d", id),
		"data":    tax,
	})
}

// DeleteVehicleTax godoc
// @Summary Hapus pajak kendaraan
// @Description Menghapus data pajak kendaraan milik user yang login berdasarkan ID
// @Tags Vehicle Tax
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID pajak"
// @Success 200 {object} map[string]interface{} "berhasil hapus data pajak kendaraan"
// @Failure 400 {object} map[string]interface{} "id tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "pajak kendaraan tidak ditemukan"
// @Router /vehicle-taxes/{id} [delete]
func (h *VehicleTaxHandler) DeleteVehicleTax(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	userID := c.GetInt64(middleware.UserIDKey)
	err = h.service.DeleteVehicleTax(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrVehicleTaxNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "pajak kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal hapus data pajak kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("berhasil hapus data pajak kendaraan, dengan id : %d", id),
	})
}
