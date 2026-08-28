package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/my-api/internal/repository"
	"example.com/my-api/internal/service"
)

type VehicleQRHandler struct {
	service service.VehicleQRService
}

func NewVehicleQRHandler(service service.VehicleQRService) *VehicleQRHandler {
	return &VehicleQRHandler{service: service}
}

// GetVehicleByQRCode godoc
// @Summary Ambil data kendaraan via QR code
// @Description Mengambil data kendaraan publik berdasarkan QR code
// @Tags Vehicles
// @Accept json
// @Produce json
// @Param id path int true "ID kendaraan"
// @Success 200 {object} map[string]interface{} "Data kendaraan berhasil diambil"
// @Failure 400 {object} map[string]interface{} "id tidak valid"
// @Failure 404 {object} map[string]interface{} "kendaraan tidak ditemukan"
// @Router /vehicles/qr/{id} [get]
func (h *VehicleQRHandler) GetVehicleByQRCode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	vehicle, err := h.service.GetVehicleByID(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrVehicleNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data kendaraan berhasil diambil",
		"data":    vehicle,
	})
}
