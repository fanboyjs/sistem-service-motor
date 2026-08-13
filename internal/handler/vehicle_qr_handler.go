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
