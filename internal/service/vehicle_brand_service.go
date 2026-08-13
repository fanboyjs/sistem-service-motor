package service

import (
	"context"
	"io"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
	"example.com/my-api/internal/storage"
	"example.com/my-api/pkg/helper"
)

type VehicleBrandService interface {
	CreateVehicleBrand(ctx context.Context, req dto.CreateVehicleBrandRequest, logo io.Reader, logoExt string) (dto.VehicleBrandResponse, error)
	GetVehicleBrands(ctx context.Context) ([]model.VehicleBrand, error)
	GetVehicleBrandByID(ctx context.Context, id int64) (model.VehicleBrand, error)
	UpdateVehicleBrand(ctx context.Context, id int64, req dto.UpdateVehicleBrandRequest, logo io.Reader, logoExt string) (dto.VehicleBrandResponse, error)
}

type vehicleBrandService struct {
	repository repository.VehicleBrandRepository
	storage    storage.Storage
}

func NewVehicleBrandService(repository repository.VehicleBrandRepository, storage storage.Storage) VehicleBrandService {
	return &vehicleBrandService{repository: repository, storage: storage}
}

func (s *vehicleBrandService) CreateVehicleBrand(ctx context.Context, req dto.CreateVehicleBrandRequest, logo io.Reader, logoExt string) (dto.VehicleBrandResponse, error) {
	var logoURL *string
	if logo != nil {
		objectPath := "brands/" + helper.RandomHex(16) + logoExt
		url, err := s.storage.Save(ctx, logo, objectPath)
		if err != nil {
			return dto.VehicleBrandResponse{}, err
		}
		logoURL = &url
	}

	brand, err := s.repository.Create(ctx, model.VehicleBrand{
		Name:    req.Name,
		LogoURL: logoURL,
	})
	if err != nil {
		if logoURL != nil {
			_ = s.storage.Delete(ctx, *logoURL)
		}
		return dto.VehicleBrandResponse{}, err
	}

	return dto.VehicleBrandResponse{
		ID:      brand.ID,
		Name:    brand.Name,
		LogoURL: brand.LogoURL,
	}, nil
}

func (s *vehicleBrandService) GetVehicleBrands(ctx context.Context) ([]model.VehicleBrand, error) {
	return s.repository.FindAll(ctx)
}

func (s *vehicleBrandService) GetVehicleBrandByID(ctx context.Context, id int64) (model.VehicleBrand, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *vehicleBrandService) UpdateVehicleBrand(ctx context.Context, id int64, req dto.UpdateVehicleBrandRequest, logo io.Reader, logoExt string) (dto.VehicleBrandResponse, error) {
	existing, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return dto.VehicleBrandResponse{}, err
	}

	var logoURL *string
	if logo != nil {
		objectPath := "brands/" + helper.RandomHex(16) + logoExt
		url, err := s.storage.Save(ctx, logo, objectPath)
		if err != nil {
			return dto.VehicleBrandResponse{}, err
		}
		logoURL = &url
	} else {
		logoURL = existing.LogoURL
	}

	brand, err := s.repository.Update(ctx, model.VehicleBrand{
		ID:      id,
		Name:    req.Name,
		LogoURL: logoURL,
	})
	if err != nil {
		if logo != nil {
			_ = s.storage.Delete(ctx, *logoURL)
		}
		return dto.VehicleBrandResponse{}, err
	}

	if logo != nil && existing.LogoURL != nil {
		_ = s.storage.Delete(ctx, *existing.LogoURL)
	}

	return dto.VehicleBrandResponse{
		ID:      brand.ID,
		Name:    brand.Name,
		LogoURL: brand.LogoURL,
	}, nil
}
