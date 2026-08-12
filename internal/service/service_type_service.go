package service

import (
	"context"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
)

type ServiceTypeService interface {
	CreateServiceType(ctx context.Context, req dto.CreateServiceTypeRequest) (dto.ServiceTypeResponse, error)
	GetServiceTypes(ctx context.Context) ([]model.ServiceType, error)
	GetServiceTypeByID(ctx context.Context, id int64) (model.ServiceType, error)
	UpdateServiceType(ctx context.Context, id int64, req dto.UpdateServiceTypeRequest) (dto.ServiceTypeResponse, error)
}

type serviceTypeService struct {
	repository repository.ServiceTypeRepository
}

func NewServiceTypeService(repository repository.ServiceTypeRepository) ServiceTypeService {
	return &serviceTypeService{repository: repository}
}

func (s *serviceTypeService) CreateServiceType(ctx context.Context, req dto.CreateServiceTypeRequest) (dto.ServiceTypeResponse, error) {
	serviceType, err := s.repository.Create(ctx, model.ServiceType{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return dto.ServiceTypeResponse{}, err
	}

	return dto.ServiceTypeResponse{
		ID:          serviceType.ID,
		Name:        serviceType.Name,
		Description: serviceType.Description,
	}, nil
}

func (s *serviceTypeService) GetServiceTypes(ctx context.Context) ([]model.ServiceType, error) {
	return s.repository.FindAll(ctx)
}

func (s *serviceTypeService) GetServiceTypeByID(ctx context.Context, id int64) (model.ServiceType, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *serviceTypeService) UpdateServiceType(ctx context.Context, id int64, req dto.UpdateServiceTypeRequest) (dto.ServiceTypeResponse, error) {
	serviceType, err := s.repository.Update(ctx, model.ServiceType{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return dto.ServiceTypeResponse{}, err
	}

	return dto.ServiceTypeResponse{
		ID:          serviceType.ID,
		Name:        serviceType.Name,
		Description: serviceType.Description,
	}, nil
}
