package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"example.com/my-api/internal/middleware"
	"example.com/my-api/internal/repository"
	"example.com/my-api/internal/service"
)

type VehicleQRHandler struct {
	service service.VehicleQRService
}

func NewVehicleQRHandler(service service.VehicleQRService) *VehicleQRHandler {
	return &VehicleQRHandler{service: service}
}

func (h *VehicleQRHandler) GetVehicleByQRToken(c *gin.Context) {
	token := c.Param("token")
	vehicle, err := h.service.GetVehicleByToken(c.Request.Context(), token)
	if err != nil {
		if err == repository.ErrQRCodeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "qr code tidak ditemukan",
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

func (h *VehicleQRHandler) RefreshVehicleQR(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	userID := c.GetInt64(middleware.UserIDKey)
	qr, err := h.service.RefreshForVehicle(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrVehicleNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal membuat ulang qr code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "QR code berhasil dibuat ulang",
		"data":    qr,
	})
}
