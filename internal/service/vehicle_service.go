package service

import (
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
	"example.com/my-api/internal/storage"
	"example.com/my-api/pkg/helper"
)

type VehicleService interface {
	CreateVehicle(ctx context.Context, userID int64, req dto.CreateVehicleRequest, image io.Reader, imageExt string) (dto.VehicleResponse, error)
	GetVehicles(ctx context.Context, userID int64) ([]dto.VehicleResponse, error)
	GetVehicleByID(ctx context.Context, id int64, userID int64) (dto.VehicleResponse, error)
	UpdateVehicle(ctx context.Context, id int64, userID int64, req dto.UpdateVehicleRequest, image io.Reader, imageExt string) (dto.VehicleResponse, error)
	DeleteVehicle(ctx context.Context, id int64, userID int64) error
}

type vehicleService struct {
	repository repository.VehicleRepository
	storage    storage.Storage
	qrService  VehicleQRService
}

func NewVehicleService(repository repository.VehicleRepository, storage storage.Storage, qrService VehicleQRService) VehicleService {
	return &vehicleService{repository: repository, storage: storage, qrService: qrService}
}

func (s *vehicleService) CreateVehicle(ctx context.Context, userID int64, req dto.CreateVehicleRequest, image io.Reader, imageExt string) (dto.VehicleResponse, error) {
	var imageURL *string
	if image != nil {
		objectPath := "vehicles/" + helper.RandomHex(16) + imageExt
		url, err := s.storage.Save(ctx, image, objectPath)
		if err != nil {
			return dto.VehicleResponse{}, err
		}
		imageURL = &url
	}

	vehicle, err := s.repository.Create(ctx, model.Vehicle{
		UserID:            userID,
		BrandID:           req.BrandID,
		ModelID:           req.ModelID,
		LicensePlate:      strings.ToUpper(strings.TrimSpace(req.LicensePlate)),
		ManufacturingYear: req.ManufacturingYear,
		Color:             nullableString(req.Color),
		PurchaseDate:      parsePurchaseDate(req.PurchaseDate),
		EngineNumber:      strings.TrimSpace(req.EngineNumber),
		CurrentMileage:    req.CurrentMileage,
		Status:            defaultStatus(req.Status),
		ImageURL:          imageURL,
	})
	if err != nil {
		if imageURL != nil {
			_ = s.storage.Delete(ctx, *imageURL)
		}
		return dto.VehicleResponse{}, err
	}

	qr, qrErr := s.qrService.GenerateForVehicle(ctx, vehicle.ID)
	if qrErr != nil {
		_ = s.repository.Delete(ctx, vehicle.ID, userID)
		if imageURL != nil {
			_ = s.storage.Delete(ctx, *imageURL)
		}
		return dto.VehicleResponse{}, qrErr
	}
	vehicle.QRToken = qr.Token

	return toVehicleResponse(vehicle), nil
}

func (s *vehicleService) GetVehicles(ctx context.Context, userID int64) ([]dto.VehicleResponse, error) {
	vehicles, err := s.repository.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.VehicleResponse, 0, len(vehicles))
	for _, vehicle := range vehicles {
		responses = append(responses, toVehicleResponse(vehicle))
	}
	return responses, nil
}

func (s *vehicleService) GetVehicleByID(ctx context.Context, id int64, userID int64) (dto.VehicleResponse, error) {
	vehicle, err := s.repository.FindByID(ctx, id, userID)
	if err != nil {
		return dto.VehicleResponse{}, err
	}
	return toVehicleResponse(vehicle), nil
}

func (s *vehicleService) UpdateVehicle(ctx context.Context, id int64, userID int64, req dto.UpdateVehicleRequest, image io.Reader, imageExt string) (dto.VehicleResponse, error) {
	existing, err := s.repository.FindByID(ctx, id, userID)
	if err != nil {
		return dto.VehicleResponse{}, err
	}

	var imageURL *string
	if image != nil {
		objectPath := "vehicles/" + helper.RandomHex(16) + imageExt
		url, err := s.storage.Save(ctx, image, objectPath)
		if err != nil {
			return dto.VehicleResponse{}, err
		}
		imageURL = &url
	} else {
		imageURL = existing.ImageURL
	}

	vehicle, err := s.repository.Update(ctx, model.Vehicle{
		ID:                id,
		UserID:            userID,
		BrandID:           req.BrandID,
		ModelID:           req.ModelID,
		LicensePlate:      strings.ToUpper(strings.TrimSpace(req.LicensePlate)),
		ManufacturingYear: req.ManufacturingYear,
		Color:             nullableString(req.Color),
		PurchaseDate:      parsePurchaseDate(req.PurchaseDate),
		EngineNumber:      strings.TrimSpace(req.EngineNumber),
		CurrentMileage:    req.CurrentMileage,
		Status:            defaultStatus(req.Status),
		ImageURL:          imageURL,
	})
	if err != nil {
		if image != nil {
			_ = s.storage.Delete(ctx, *imageURL)
		}
		return dto.VehicleResponse{}, err
	}

	if image != nil && existing.ImageURL != nil {
		_ = s.storage.Delete(ctx, *existing.ImageURL)
	}

	return toVehicleResponse(vehicle), nil
}

func (s *vehicleService) DeleteVehicle(ctx context.Context, id int64, userID int64) error {
	existing, err := s.repository.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}

	if err := s.repository.Delete(ctx, id, userID); err != nil {
		return err
	}

	if existing.ImageURL != nil {
		_ = s.storage.Delete(ctx, *existing.ImageURL)
	}

	_ = s.qrService.CleanupForVehicle(ctx, id)

	return nil
}

func toVehicleResponse(vehicle model.Vehicle) dto.VehicleResponse {
	var qrImageURL *string
	var qrToken *string
	if vehicle.QRToken != "" {
		imageURL := "/uploads/qr/" + strconv.FormatInt(vehicle.ID, 10) + ".png"
		token := vehicle.QRToken
		qrImageURL = &imageURL
		qrToken = &token
	}
	return dto.VehicleResponse{
		ID:                vehicle.ID,
		UserID:            vehicle.UserID,
		BrandID:           vehicle.BrandID,
		BrandName:         vehicle.BrandName,
		ModelID:           vehicle.ModelID,
		ModelName:         vehicle.ModelName,
		LicensePlate:      vehicle.LicensePlate,
		ManufacturingYear: vehicle.ManufacturingYear,
		Color:             vehicle.Color,
		PurchaseDate:      vehicle.PurchaseDate.Format("2006-01-02"),
		EngineNumber:      vehicle.EngineNumber,
		CurrentMileage:    vehicle.CurrentMileage,
		Status:            vehicle.Status,
		ImageURL:          vehicle.ImageURL,
		QRImageURL:        qrImageURL,
		QRToken:           qrToken,
	}
}

func nullableString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func parsePurchaseDate(value string) time.Time {
	t, _ := time.Parse("2006-01-02", value)
	return t
}

func defaultStatus(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "active"
	}
	return value
}
