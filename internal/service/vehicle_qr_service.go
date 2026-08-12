package service

import (
	"context"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
)

type VehicleQRService interface {
	GetVehicleByID(ctx context.Context, id int64) (dto.VehicleQRScanResponse, error)
}

type vehicleQRService struct {
	qrRepository repository.VehicleQRRepository
}

func NewVehicleQRService(qrRepository repository.VehicleQRRepository) VehicleQRService {
	return &vehicleQRService{qrRepository: qrRepository}
}

func (s *vehicleQRService) GetVehicleByID(ctx context.Context, id int64) (dto.VehicleQRScanResponse, error) {
	vehicle, err := s.qrRepository.FindVehicleByID(ctx, id)
	if err != nil {
		return dto.VehicleQRScanResponse{}, err
	}
	return toVehicleQRScanResponse(vehicle), nil
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
