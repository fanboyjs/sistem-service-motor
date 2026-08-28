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

// CreateVehicleBrand godoc
// @Summary Buat brand kendaraan
// @Description Membuat data brand kendaraan baru dengan upload logo
// @Tags Vehicle Brand
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Nama brand"
// @Param logo_url formData file false "File logo (png, jpg, jpeg, webp, max 5MB)"
// @Success 201 {object} map[string]interface{} "Data brand kendaraan berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 409 {object} map[string]interface{} "nama brand sudah terdaftar"
// @Router /vehicle-brands [post]
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

// GetVehicleBrands godoc
// @Summary Ambil semua brand kendaraan
// @Description Mengambil daftar semua brand kendaraan
// @Tags Vehicle Brand
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Data brand kendaraan berhasil diambil"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Router /vehicle-brands [get]
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

// GetVehicleBrandById godoc
// @Summary Ambil brand kendaraan berdasarkan ID
// @Description Mengambil data brand kendaraan berdasarkan ID
// @Tags Vehicle Brand
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID brand"
// @Success 200 {object} map[string]interface{} "Data brand kendaraan berhasil diambil"
// @Failure 400 {object} map[string]interface{} "id tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "brand kendaraan tidak ditemukan"
// @Router /vehicle-brands/{id} [get]
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

// UpdateVehicleBrand godoc
// @Summary Update brand kendaraan
// @Description Memperbarui data brand kendaraan berdasarkan ID (opsional upload logo baru)
// @Tags Vehicle Brand
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID brand"
// @Param name formData string true "Nama brand"
// @Param logo_url formData file false "File logo (png, jpg, jpeg, webp, max 5MB)"
// @Success 200 {object} map[string]interface{} "Data brand kendaraan berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "request body tidak valid"
// @Failure 401 {object} map[string]interface{} "token tidak valid"
// @Failure 404 {object} map[string]interface{} "brand kendaraan tidak ditemukan"
// @Failure 409 {object} map[string]interface{} "nama brand sudah terdaftar"
// @Router /vehicle-brands/{id} [put]
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
