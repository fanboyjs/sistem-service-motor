package service

import (
	"context"
	"strings"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
)

type VehicleTaxService interface {
	CreateVehicleTax(ctx context.Context, userID int64, req dto.CreateVehicleTaxRequest) (dto.VehicleTaxResponse, error)
	GetVehicleTaxes(ctx context.Context, userID int64) ([]dto.VehicleTaxResponse, error)
	GetVehicleTaxByID(ctx context.Context, id int64, userID int64) (dto.VehicleTaxResponse, error)
	UpdateVehicleTax(ctx context.Context, id int64, userID int64, req dto.UpdateVehicleTaxRequest) (dto.VehicleTaxResponse, error)
	DeleteVehicleTax(ctx context.Context, id int64, userID int64) error
}

type vehicleTaxService struct {
	repository repository.VehicleTaxRepository
}

func NewVehicleTaxService(repository repository.VehicleTaxRepository) VehicleTaxService {
	return &vehicleTaxService{repository: repository}
}

func (s *vehicleTaxService) CreateVehicleTax(ctx context.Context, userID int64, req dto.CreateVehicleTaxRequest) (dto.VehicleTaxResponse, error) {
	tax, err := s.repository.Create(ctx, model.VehicleTax{
		VehicleID:     req.VehicleID,
		TaxYear:       req.TaxYear,
		PKBAmount:     req.PKBAmount,
		SWDKLLJAmount: req.SWDKLLJAmount,
		OtherAmount:   req.OtherAmount,
		TotalAmount:   req.PKBAmount + req.SWDKLLJAmount + req.OtherAmount,
		PaymentDate:   parsePurchaseDate(req.PaymentDate),
		DueDate:       parsePurchaseDate(req.DueDate),
		Status:        defaultTaxStatus(req.Status),
		Notes:         req.Notes,
	}, userID)
	if err != nil {
		return dto.VehicleTaxResponse{}, err
	}

	return toVehicleTaxResponse(tax), nil
}

func (s *vehicleTaxService) GetVehicleTaxes(ctx context.Context, userID int64) ([]dto.VehicleTaxResponse, error) {
	taxes, err := s.repository.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.VehicleTaxResponse, 0, len(taxes))
	for _, tax := range taxes {
		responses = append(responses, toVehicleTaxResponse(tax))
	}
	return responses, nil
}

func (s *vehicleTaxService) GetVehicleTaxByID(ctx context.Context, id int64, userID int64) (dto.VehicleTaxResponse, error) {
	tax, err := s.repository.FindByID(ctx, id, userID)
	if err != nil {
		return dto.VehicleTaxResponse{}, err
	}

	return toVehicleTaxResponse(tax), nil
}

func (s *vehicleTaxService) UpdateVehicleTax(ctx context.Context, id int64, userID int64, req dto.UpdateVehicleTaxRequest) (dto.VehicleTaxResponse, error) {
	tax, err := s.repository.Update(ctx, model.VehicleTax{
		ID:            id,
		VehicleID:     req.VehicleID,
		TaxYear:       req.TaxYear,
		PKBAmount:     req.PKBAmount,
		SWDKLLJAmount: req.SWDKLLJAmount,
		OtherAmount:   req.OtherAmount,
		TotalAmount:   req.PKBAmount + req.SWDKLLJAmount + req.OtherAmount,
		PaymentDate:   parsePurchaseDate(req.PaymentDate),
		DueDate:       parsePurchaseDate(req.DueDate),
		Status:        defaultTaxStatus(req.Status),
		Notes:         req.Notes,
	}, userID)
	if err != nil {
		return dto.VehicleTaxResponse{}, err
	}

	return toVehicleTaxResponse(tax), nil
}

func (s *vehicleTaxService) DeleteVehicleTax(ctx context.Context, id int64, userID int64) error {
	return s.repository.Delete(ctx, id, userID)
}

func toVehicleTaxResponse(tax model.VehicleTax) dto.VehicleTaxResponse {
	return dto.VehicleTaxResponse{
		ID:            tax.ID,
		VehicleID:     tax.VehicleID,
		LicensePlate:  tax.LicensePlate,
		TaxYear:       tax.TaxYear,
		PKBAmount:     tax.PKBAmount,
		SWDKLLJAmount: tax.SWDKLLJAmount,
		OtherAmount:   tax.OtherAmount,
		TotalAmount:   tax.TotalAmount,
		PaymentDate:   tax.PaymentDate.Format("2006-01-02"),
		DueDate:       tax.DueDate.Format("2006-01-02"),
		Status:        tax.Status,
		Notes:         tax.Notes,
	}
}

func defaultTaxStatus(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unpaid"
	}
	return value
}
