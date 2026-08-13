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
	"example.com/my-api/internal/repository"
	"example.com/my-api/internal/service"
	"example.com/my-api/pkg/validator"
)

const maxLogoSize = 5 << 20

var allowedLogoExts = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
}

type VehicleBrandHandler struct {
	service service.VehicleBrandService
}

func NewVehicleBrandHandler(service service.VehicleBrandService) *VehicleBrandHandler {
	return &VehicleBrandHandler{service: service}
}

func (h *VehicleBrandHandler) CreateVehicleBrand(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxLogoSize)
	if err := c.Request.ParseMultipartForm(maxLogoSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request body tidak valid",
		})
		return
	}

	req := dto.CreateVehicleBrandRequest{
		Name: strings.TrimSpace(c.PostForm("name")),
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var logo io.Reader
	var logoExt string
	if fileHeader, err := c.FormFile("logo_url"); err == nil {
		if fileHeader.Size > maxLogoSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"message": "ukuran file logo melebihi batas 5MB",
			})
			return
		}
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !allowedLogoExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "format file logo tidak valid, gunakan png, jpg, jpeg, atau webp",
			})
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "gagal membaca file logo",
			})
			return
		}
		defer file.Close()
		logo = file
		logoExt = ext
	}

	brand, err := h.service.CreateVehicleBrand(c.Request.Context(), req, logo, logoExt)
	if err != nil {
		if err == repository.ErrBrandNameExists {
			c.JSON(http.StatusConflict, gin.H{
				"message": "nama brand sudah terdaftar",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal membuat brand kendaraan",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Data brand kendaraan berhasil dibuat",
		"data":    brand,
	})
}

func (h *VehicleBrandHandler) GetVehicleBrands(c *gin.Context) {
	brands, err := h.service.GetVehicleBrands(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data brand kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data brand kendaraan berhasil diambil",
		"data":    brands,
	})
}

func (h *VehicleBrandHandler) GetVehicleBrandById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	brand, err := h.service.GetVehicleBrandByID(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrVehicleBrandNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "brand kendaraan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal ambil data brand kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data brand kendaraan berhasil diambil, dengan id : %d", id),
		"data":    brand,
	})
}

func (h *VehicleBrandHandler) UpdateVehicleBrand(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id tidak valid",
		})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxLogoSize)
	if err := c.Request.ParseMultipartForm(maxLogoSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request body tidak valid",
		})
		return
	}

	req := dto.UpdateVehicleBrandRequest{
		Name: strings.TrimSpace(c.PostForm("name")),
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var logo io.Reader
	var logoExt string
	if fileHeader, err := c.FormFile("logo_url"); err == nil {
		if fileHeader.Size > maxLogoSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"message": "ukuran file logo melebihi batas 5MB",
			})
			return
		}
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !allowedLogoExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "format file logo tidak valid, gunakan png, jpg, jpeg, atau webp",
			})
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "gagal membaca file logo",
			})
			return
		}
		defer file.Close()
		logo = file
		logoExt = ext
	}

	brand, err := h.service.UpdateVehicleBrand(c.Request.Context(), id, req, logo, logoExt)
	if err != nil {
		if err == repository.ErrVehicleBrandNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "brand kendaraan tidak ditemukan",
			})
			return
		}
		if err == repository.ErrBrandNameExists {
			c.JSON(http.StatusConflict, gin.H{
				"message": "nama brand sudah terdaftar",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal update data brand kendaraan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Data brand kendaraan berhasil diupdate, dengan id : %d", id),
		"data":    brand,
	})
}
