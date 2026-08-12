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
	repository repository.ServiceRecordRepository
}

func NewServiceRecordService(repository repository.ServiceRecordRepository) ServiceRecordService {
	return &serviceRecordService{repository: repository}
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
	for _, record := range records {
		responses = append(responses, toServiceRecordResponse(record))
	}
	return responses, nil
}

func (s *serviceRecordService) GetServiceRecordByID(ctx context.Context, id int64, userID int64) (dto.ServiceRecordResponse, error) {
	record, err := s.repository.FindByID(ctx, id, userID)
	if err != nil {
		return dto.ServiceRecordResponse{}, err
	}

	return toServiceRecordResponse(record), nil
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
