package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/my-api/internal/model"
)

var ErrQRCodeNotFound = errors.New("qr code tidak ditemukan")

type VehicleQRRepository interface {
	Upsert(ctx context.Context, vehicleID int64, token string) (model.VehicleQRCode, error)
	FindVehicleByToken(ctx context.Context, token string) (model.Vehicle, error)
}

type vehicleQRRepository struct {
	db *pgxpool.Pool
}

func NewVehicleQRRepository(db *pgxpool.Pool) VehicleQRRepository {
	return &vehicleQRRepository{db: db}
}

func (r *vehicleQRRepository) Upsert(ctx context.Context, vehicleID int64, token string) (model.VehicleQRCode, error) {
	var qr model.VehicleQRCode
	err := r.db.QueryRow(ctx, `
		INSERT INTO vehicle_qr_codes (vehicle_id, token)
		VALUES ($1, $2)
		ON CONFLICT (vehicle_id) DO UPDATE
			SET token = EXCLUDED.token, status = 'active', updated_at = CURRENT_TIMESTAMP
		RETURNING id, vehicle_id, token, status, created_at, updated_at
	`, vehicleID, token).Scan(
		&qr.ID,
		&qr.VehicleID,
		&qr.Token,
		&qr.Status,
		&qr.CreatedAt,
		&qr.UpdatedAt,
	)
	if err != nil {
		return model.VehicleQRCode{}, err
	}
	return qr, nil
}

func (r *vehicleQRRepository) FindVehicleByToken(ctx context.Context, token string) (model.Vehicle, error) {
	var vehicle model.Vehicle
	err := r.db.QueryRow(ctx, `
		SELECT `+vehicleSelectColumns+`
		FROM vehicle_qr_codes q
		JOIN vehicles v ON v.id = q.vehicle_id
		JOIN vehicle_brands b ON b.id = v.brand_id
		JOIN vehicle_models m ON m.id = v.model_id
		WHERE q.token = $1 AND q.status = 'active'
	`, token).Scan(
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
		&vehicle.QRToken,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Vehicle{}, ErrQRCodeNotFound
	}
	if err != nil {
		return model.Vehicle{}, err
	}
	return vehicle, nil
}
