package service

import (
	"context"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
)

type ServiceRecordService interface {
	CreateServiceRecord(ctx context.Context, userID int64, req dto.CreateServiceRecordRequest) (dto.ServiceRecordResponse, error)
	GetServiceRecords(ctx context.Context, userID int64) ([]dto.ServiceRecordResponse, error)
	GetServiceRecordByID(ctx context.Context, id int64, userID int64) (dto.ServiceRecordResponse, error)
	UpdateServiceRecord(ctx context.Context, id int64, userID int64, req dto.UpdateServiceRecordRequest) (dto.ServiceRecordResponse, error)
	DeleteServiceRecord(ctx context.Context, id int64, userID int64) error
}

type serviceRecordService struct {
	repository             repository.ServiceRecordRepository
	serviceItemRepository  repository.ServiceItemRepository
}

func NewServiceRecordService(repository repository.ServiceRecordRepository, serviceItemRepository repository.ServiceItemRepository) ServiceRecordService {
	return &serviceRecordService{repository: repository, serviceItemRepository: serviceItemRepository}
}

func (s *serviceRecordService) CreateServiceRecord(ctx context.Context, userID int64, req dto.CreateServiceRecordRequest) (dto.ServiceRecordResponse, error) {
	record, err := s.repository.Create(ctx, model.ServiceRecord{
		VehicleID:     req.VehicleID,
		ServiceTypeID: req.ServiceTypeID,
		ServiceDate:   parsePurchaseDate(req.ServiceDate),
		Odometer:      req.Odometer,
		WorkshopName:  req.WorkshopName,
		LaborCost:     req.LaborCost,
		PartsCost:     req.PartsCost,
		TotalCost:     req.LaborCost + req.PartsCost,
		Notes:         req.Notes,
	}, userID)
	if err != nil {
		return dto.ServiceRecordResponse{}, err
	}

	return toServiceRecordResponse(record), nil
}

func (s *serviceRecordService) GetServiceRecords(ctx context.Context, userID int64) ([]dto.ServiceRecordResponse, error) {
	records, err := s.repository.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ServiceRecordResponse, 0, len(records))
	recordIDs := make([]int64, 0, len(records))
	idToIndex := make(map[int64]int, len(records))
	for i, record := range records {
		responses = append(responses, toServiceRecordResponse(record))
		recordIDs = append(recordIDs, record.ID)
		idToIndex[record.ID] = i
	}

	if len(recordIDs) > 0 {
		items, err := s.serviceItemRepository.FindByServiceRecordIDs(ctx, recordIDs)
		if err != nil {
			return nil, err
		}
		itemsByRecord := make(map[int64][]dto.ServiceItemResponse)
		for _, item := range items {
			itemsByRecord[item.ServiceRecordID] = append(itemsByRecord[item.ServiceRecordID], toServiceItemResponse(item))
		}
		for recordID, index := range idToIndex {
			items := itemsByRecord[recordID]
			if items == nil {
				items = make([]dto.ServiceItemResponse, 0)
			}
			responses[index].ServiceItems = items
			for _, item := range items {
				responses[index].TotalCost += float64(item.Quantity) * item.Cost
			}
		}
	}

	return responses, nil
}

func (s *serviceRecordService) GetServiceRecordByID(ctx context.Context, id int64, userID int64) (dto.ServiceRecordResponse, error) {
	record, err := s.repository.FindByID(ctx, id, userID)
	if err != nil {
		return dto.ServiceRecordResponse{}, err
	}

	response := toServiceRecordResponse(record)
	items, err := s.serviceItemRepository.FindByServiceRecordIDs(ctx, []int64{record.ID})
	if err != nil {
		return dto.ServiceRecordResponse{}, err
	}
	response.ServiceItems = make([]dto.ServiceItemResponse, 0, len(items))
	for _, item := range items {
		response.ServiceItems = append(response.ServiceItems, toServiceItemResponse(item))
		response.TotalCost += float64(item.Quantity) * item.Cost
	}

	return response, nil
}

func (s *serviceRecordService) UpdateServiceRecord(ctx context.Context, id int64, userID int64, req dto.UpdateServiceRecordRequest) (dto.ServiceRecordResponse, error) {
	record, err := s.repository.Update(ctx, model.ServiceRecord{
		ID:            id,
		VehicleID:     req.VehicleID,
		ServiceTypeID: req.ServiceTypeID,
		ServiceDate:   parsePurchaseDate(req.ServiceDate),
		Odometer:      req.Odometer,
		WorkshopName:  req.WorkshopName,
		LaborCost:     req.LaborCost,
		PartsCost:     req.PartsCost,
		TotalCost:     req.LaborCost + req.PartsCost,
		Notes:         req.Notes,
	}, userID)
	if err != nil {
		return dto.ServiceRecordResponse{}, err
	}

	return toServiceRecordResponse(record), nil
}

func (s *serviceRecordService) DeleteServiceRecord(ctx context.Context, id int64, userID int64) error {
	return s.repository.Delete(ctx, id, userID)
}

func toServiceRecordResponse(record model.ServiceRecord) dto.ServiceRecordResponse {
	return dto.ServiceRecordResponse{
		ID:              record.ID,
		VehicleID:       record.VehicleID,
		LicensePlate:    record.LicensePlate,
		ServiceTypeID:   record.ServiceTypeID,
		ServiceTypeName: record.ServiceTypeName,
		ServiceDate:     record.ServiceDate.Format("2006-01-02"),
		Odometer:        record.Odometer,
		WorkshopName:    record.WorkshopName,
		LaborCost:       record.LaborCost,
		PartsCost:       record.PartsCost,
		TotalCost:       record.TotalCost,
		Notes:           record.Notes,
	}
}

func toServiceItemResponse(serviceItem model.ServiceItem) dto.ServiceItemResponse {
	return dto.ServiceItemResponse{
		ID:              serviceItem.ID,
		ServiceRecordID: serviceItem.ServiceRecordID,
		Name:            serviceItem.Name,
		Quantity:        serviceItem.Quantity,
		Cost:            serviceItem.Cost,
		Notes:           serviceItem.Notes,
	}
}
