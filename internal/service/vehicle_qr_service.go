package service

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	qrcode "github.com/skip2/go-qrcode"

	"example.com/my-api/config"
	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
	"example.com/my-api/internal/storage"
	"example.com/my-api/pkg/helper"
)

type VehicleQRService interface {
	GenerateForVehicle(ctx context.Context, vehicleID int64) (dto.VehicleQRResponse, error)
	GetVehicleByToken(ctx context.Context, token string) (dto.VehicleQRScanResponse, error)
	RefreshForVehicle(ctx context.Context, vehicleID int64, userID int64) (dto.VehicleQRResponse, error)
	CleanupForVehicle(ctx context.Context, vehicleID int64) error
}

type vehicleQRService struct {
	qrRepository      repository.VehicleQRRepository
	vehicleRepository repository.VehicleRepository
	storage           storage.Storage
	baseURL           string
}

func NewVehicleQRService(qrRepository repository.VehicleQRRepository, vehicleRepository repository.VehicleRepository, storage storage.Storage, cfg config.Config) VehicleQRService {
	return &vehicleQRService{
		qrRepository:      qrRepository,
		vehicleRepository: vehicleRepository,
		storage:           storage,
		baseURL:           strings.TrimRight(cfg.PublicBaseURL, "/"),
	}
}

func (s *vehicleQRService) GenerateForVehicle(ctx context.Context, vehicleID int64) (dto.VehicleQRResponse, error) {
	token := "qr_" + helper.RandomHex(16)
	if _, err := s.qrRepository.Upsert(ctx, vehicleID, token); err != nil {
		return dto.VehicleQRResponse{}, err
	}

	scanURL := s.baseURL + "/api/qr/vehicle/" + token
	imageURL, err := s.saveQRImage(ctx, vehicleID, scanURL)
	if err != nil {
		return dto.VehicleQRResponse{}, err
	}

	return dto.VehicleQRResponse{
		Token:    token,
		ImageURL: imageURL,
		ScanURL:  scanURL,
	}, nil
}

func (s *vehicleQRService) GetVehicleByToken(ctx context.Context, token string) (dto.VehicleQRScanResponse, error) {
	vehicle, err := s.qrRepository.FindVehicleByToken(ctx, token)
	if err != nil {
		return dto.VehicleQRScanResponse{}, err
	}
	return toVehicleQRScanResponse(vehicle), nil
}

func (s *vehicleQRService) RefreshForVehicle(ctx context.Context, vehicleID int64, userID int64) (dto.VehicleQRResponse, error) {
	if _, err := s.vehicleRepository.FindByID(ctx, vehicleID, userID); err != nil {
		return dto.VehicleQRResponse{}, err
	}
	return s.GenerateForVehicle(ctx, vehicleID)
}

func (s *vehicleQRService) CleanupForVehicle(ctx context.Context, vehicleID int64) error {
	return s.storage.Delete(ctx, qrImagePath(vehicleID))
}

func (s *vehicleQRService) saveQRImage(ctx context.Context, vehicleID int64, scanURL string) (string, error) {
	png, err := qrcode.Encode(scanURL, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	return s.storage.Save(ctx, bytes.NewReader(png), qrImagePath(vehicleID))
}

func qrImagePath(vehicleID int64) string {
	return "qr/" + strconv.FormatInt(vehicleID, 10) + ".png"
}

func toVehicleQRScanResponse(vehicle model.Vehicle) dto.VehicleQRScanResponse {
	return dto.VehicleQRScanResponse{
		ID:                vehicle.ID,
		BrandName:         vehicle.BrandName,
		ModelName:         vehicle.ModelName,
		LicensePlate:      vehicle.LicensePlate,
		ManufacturingYear: vehicle.ManufacturingYear,
		Color:             vehicle.Color,
		PurchaseDate:      vehicle.PurchaseDate.Format("2006-01-02"),
		CurrentMileage:    vehicle.CurrentMileage,
		Status:            vehicle.Status,
		ImageURL:          vehicle.ImageURL,
	}
}
