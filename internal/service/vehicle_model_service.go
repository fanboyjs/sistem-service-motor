package service

import (
	"context"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
)

type VehicleModelService interface {
	CreateVehicleModel(ctx context.Context, req dto.CreateVehicleModelRequest) (dto.VehicleModelResponse, error)
	GetVehicleModels(ctx context.Context) ([]model.VehicleModel, error)
	GetVehicleModelByID(ctx context.Context, id int64) (model.VehicleModel, error)
	UpdateVehicleModel(ctx context.Context, id int64, req dto.UpdateVehicleModelRequest) (dto.VehicleModelResponse, error)
}

type vehicleModelService struct {
	repository repository.VehicleModelRepository
}

func NewVehicleModelService(repository repository.VehicleModelRepository) VehicleModelService {
	return &vehicleModelService{repository: repository}
}

func (s *vehicleModelService) CreateVehicleModel(ctx context.Context, req dto.CreateVehicleModelRequest) (dto.VehicleModelResponse, error) {
	modelData, err := s.repository.Create(ctx, model.VehicleModel{
		BrandID: req.BrandID,
		Name:    req.Name,
	})
	if err != nil {
		return dto.VehicleModelResponse{}, err
	}

	return dto.VehicleModelResponse{
		ID:        modelData.ID,
		BrandID:   modelData.BrandID,
		BrandName: modelData.BrandName,
		Name:      modelData.Name,
	}, nil
}

func (s *vehicleModelService) GetVehicleModels(ctx context.Context) ([]model.VehicleModel, error) {
	return s.repository.FindAll(ctx)
}

func (s *vehicleModelService) GetVehicleModelByID(ctx context.Context, id int64) (model.VehicleModel, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *vehicleModelService) UpdateVehicleModel(ctx context.Context, id int64, req dto.UpdateVehicleModelRequest) (dto.VehicleModelResponse, error) {
	modelData, err := s.repository.Update(ctx, model.VehicleModel{
		ID:      id,
		BrandID: req.BrandID,
		Name:    req.Name,
	})
	if err != nil {
		return dto.VehicleModelResponse{}, err
	}

	return dto.VehicleModelResponse{
		ID:        modelData.ID,
		BrandID:   modelData.BrandID,
		BrandName: modelData.BrandName,
		Name:      modelData.Name,
	}, nil
}
