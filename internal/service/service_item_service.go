package service

import (
	"context"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
)

type ServiceItemService interface {
	CreateServiceItem(ctx context.Context, req dto.CreateServiceItemRequest) (dto.ServiceItemResponse, error)
	GetServiceItems(ctx context.Context) ([]model.ServiceItem, error)
	GetServiceItemByID(ctx context.Context, id int64) (model.ServiceItem, error)
	UpdateServiceItem(ctx context.Context, id int64, req dto.UpdateServiceItemRequest) (dto.ServiceItemResponse, error)
}

type serviceItemService struct {
	repository repository.ServiceItemRepository
}

func NewServiceItemService(repository repository.ServiceItemRepository) ServiceItemService {
	return &serviceItemService{repository: repository}
}

func (s *serviceItemService) CreateServiceItem(ctx context.Context, req dto.CreateServiceItemRequest) (dto.ServiceItemResponse, error) {
	serviceItem, err := s.repository.Create(ctx, model.ServiceItem{
		ServiceRecordID: req.ServiceRecordID,
		Name:            req.Name,
		Quantity:        req.Quantity,
		Cost:            req.Cost,
		Notes:           req.Notes,
	})
	if err != nil {
		return dto.ServiceItemResponse{}, err
	}

	return dto.ServiceItemResponse{
		ID:              serviceItem.ID,
		ServiceRecordID: serviceItem.ServiceRecordID,
		Name:            serviceItem.Name,
		Quantity:        serviceItem.Quantity,
		Cost:            serviceItem.Cost,
		Notes:           serviceItem.Notes,
	}, nil
}

func (s *serviceItemService) GetServiceItems(ctx context.Context) ([]model.ServiceItem, error) {
	return s.repository.FindAll(ctx)
}

func (s *serviceItemService) GetServiceItemByID(ctx context.Context, id int64) (model.ServiceItem, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *serviceItemService) UpdateServiceItem(ctx context.Context, id int64, req dto.UpdateServiceItemRequest) (dto.ServiceItemResponse, error) {
	serviceItem, err := s.repository.Update(ctx, model.ServiceItem{
		ID:              id,
		ServiceRecordID: req.ServiceRecordID,
		Name:            req.Name,
		Quantity:        req.Quantity,
		Cost:            req.Cost,
		Notes:           req.Notes,
	})
	if err != nil {
		return dto.ServiceItemResponse{}, err
	}

	return dto.ServiceItemResponse{
		ID:              serviceItem.ID,
		ServiceRecordID: serviceItem.ServiceRecordID,
		Name:            serviceItem.Name,
		Quantity:        serviceItem.Quantity,
		Cost:            serviceItem.Cost,
		Notes:           serviceItem.Notes,
	}, nil
}
