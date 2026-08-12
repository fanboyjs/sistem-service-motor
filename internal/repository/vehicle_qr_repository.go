package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/my-api/internal/model"
)

type VehicleQRRepository interface {
	FindVehicleByID(ctx context.Context, id int64) (model.Vehicle, error)
}

type vehicleQRRepository struct {
	db *pgxpool.Pool
}

func NewVehicleQRRepository(db *pgxpool.Pool) VehicleQRRepository {
	return &vehicleQRRepository{db: db}
}

func (r *vehicleQRRepository) FindVehicleByID(ctx context.Context, id int64) (model.Vehicle, error) {
	var vehicle model.Vehicle
	err := r.db.QueryRow(ctx, `
		SELECT `+vehicleSelectColumns+`
		FROM vehicles v
		JOIN vehicle_brands b ON b.id = v.brand_id
		JOIN vehicle_models m ON m.id = v.model_id
		WHERE v.id = $1
	`, id).Scan(
		&vehicle.ID,
		&vehicle.UserID,
		&vehicle.BrandID,
		&vehicle.BrandName,
		&vehicle.ModelID,
		&vehicle.ModelName,
		&vehicle.LicensePlate,
		&vehicle.ManufacturingYear,
		&vehicle.Color,
		&vehicle.PurchaseDate,
		&vehicle.EngineNumber,
		&vehicle.CurrentMileage,
		&vehicle.Status,
		&vehicle.ImageURL,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Vehicle{}, ErrVehicleNotFound
	}
	if err != nil {
		return model.Vehicle{}, err
	}
	return vehicle, nil
}
