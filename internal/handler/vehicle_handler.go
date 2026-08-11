package handler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/middleware"
	"example.com/my-api/internal/repository"
	"example.com/my-api/internal/service"
	"example.com/my-api/pkg/validator"
)

const maxVehicleImageSize = 5 << 20

var allowedVehicleImageExts = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
}

type VehicleHandler struct {
	service service.VehicleService
}

func NewVehicleHandler(service service.VehicleService) *VehicleHandler {
	return &VehicleHandler{service: service}
}

func (h *VehicleHandler) CreateVehicle(c *gin.Context) {
	req, ok := parseVehicleForm(c)
	if !ok {
		return
	}

	var image io.Reader
	var imageExt string
	if fileHeader, err := c.FormFile("image_url"); err == nil {
		if fileHeader.Size > maxVehicleImageSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"message": "ukuran file gambar melebihi batas 5MB",
			})
			return
		}
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !allowedVehicleImageExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "format file gambar tidak valid, gunakan png, jpg, jpeg, atau webp",
			})
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "gagal membaca file gambar",
			})
			return
		}
		defer file.Close()
		image = file
		imageExt = ext
	}

	userID := c.GetInt64(middleware.UserIDKey)
	vehicle, err := h.service.CreateVehicle(c.Request.Context(), userID, req, image, imageExt)
	if err != nil {
		vehicleErrorResponse(c, err, "gagal membuat data kendaraan")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data kendaraan berhasil dibuat",
		"data":    vehicle,
	})
}

func (h *VehicleHandler) GetVehicles(c *gin.Context) {
	userID := c.GetInt64(middleware.UserIDKey)
	vehicles, err := h.service.GetVehicles(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data kendaraan berhasil diambil",
		"data":    vehicles,
	})
}

func (h *VehicleHandler) GetVehicleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	userID := c.GetInt64(middleware.UserIDKey)
	vehicle, err := h.service.GetVehicleByID(c.Request.Context(), id, userID)
	if err != nil {
		vehicleErrorResponse(c, err, "gagal ambil data kendaraan")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data kendaraan berhasil diambil, dengan id : %d", id),
		"data":    vehicle,
	})
}

func (h *VehicleHandler) UpdateVehicle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	req, ok := parseVehicleForm(c)
	if !ok {
		return
	}

	var image io.Reader
	var imageExt string
	if fileHeader, err := c.FormFile("image_url"); err == nil {
		if fileHeader.Size > maxVehicleImageSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"message": "ukuran file gambar melebihi batas 5MB",
			})
			return
		}
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !allowedVehicleImageExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "format file gambar tidak valid, gunakan png, jpg, jpeg, atau webp",
			})
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "gagal membaca file gambar",
			})
			return
		}
		defer file.Close()
		image = file
		imageExt = ext
	}

	userID := c.GetInt64(middleware.UserIDKey)
	vehicle, err := h.service.UpdateVehicle(c.Request.Context(), id, userID, dto.UpdateVehicleRequest(req), image, imageExt)
	if err != nil {
		vehicleErrorResponse(c, err, "gagal update data kendaraan")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data kendaraan berhasil diupdate, dengan id : %d", id),
		"data":    vehicle,
	})
}

func (h *VehicleHandler) DeleteVehicle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	userID := c.GetInt64(middleware.UserIDKey)
	err = h.service.DeleteVehicle(c.Request.Context(), id, userID)
	if err != nil {
		if err == repository.ErrVehicleNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal hapus data kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("berhasil hapus data kendaraan, dengan id : %d", id),
	})
}

func parseVehicleForm(c *gin.Context) (dto.CreateVehicleRequest, bool) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxVehicleImageSize)
	if err := c.Request.ParseMultipartForm(maxVehicleImageSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request body tidak valid",
		})
		return dto.CreateVehicleRequest{}, false
	}

	brandID, err := strconv.ParseInt(c.PostForm("brand_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "brand_id tidak valid",
		})
		return dto.CreateVehicleRequest{}, false
	}
	modelID, err := strconv.ParseInt(c.PostForm("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "model_id tidak valid",
		})
		return dto.CreateVehicleRequest{}, false
	}
	manufacturingYear, err := strconv.Atoi(strings.TrimSpace(c.PostForm("manufacturing_year")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "manufacturing_year tidak valid",
		})
		return dto.CreateVehicleRequest{}, false
	}
	currentMileage, err := strconv.Atoi(strings.TrimSpace(c.PostForm("current_mileage")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "current_mileage tidak valid",
		})
		return dto.CreateVehicleRequest{}, false
	}

	req := dto.CreateVehicleRequest{
		BrandID:           brandID,
		ModelID:           modelID,
		LicensePlate:      strings.TrimSpace(c.PostForm("license_plate")),
		ManufacturingYear: manufacturingYear,
		Color:             strings.TrimSpace(c.PostForm("color")),
		PurchaseDate:      strings.TrimSpace(c.PostForm("purchase_date")),
		EngineNumber:      strings.TrimSpace(c.PostForm("engine_number")),
		CurrentMileage:    currentMileage,
		Status:            strings.TrimSpace(c.PostForm("status")),
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return dto.CreateVehicleRequest{}, false
	}

	return req, true
}

func vehicleErrorResponse(c *gin.Context, err error, defaultMessage string) {
	if err == repository.ErrVehicleNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "kendaraan tidak ditemukan",
		})
		return
	}
	if err == repository.ErrVehicleBrandNotFoundFK {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "brand kendaraan tidak ditemukan",
		})
		return
	}
	if err == repository.ErrVehicleModelNotFoundFK {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "model kendaraan tidak ditemukan",
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": defaultMessage,
	})
}
