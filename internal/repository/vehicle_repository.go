package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/my-api/internal/model"
)

var ErrVehicleNotFound = errors.New("kendaraan tidak ditemukan")
var ErrVehicleModelNotFoundFK = errors.New("model kendaraan tidak ditemukan")

const vehicleSelectColumns = `
	v.id, v.user_id, v.brand_id, b.name, v.model_id, m.name,
	v.license_plate, v.manufacturing_year, v.color, v.purchase_date,
	v.engine_number, v.current_mileage, v.status, v.image_url,
	v.created_at, v.updated_at`

type VehicleRepository interface {
	Create(ctx context.Context, vehicle model.Vehicle) (model.Vehicle, error)
	FindAll(ctx context.Context, userID int64) ([]model.Vehicle, error)
	FindByID(ctx context.Context, id int64, userID int64) (model.Vehicle, error)
	Update(ctx context.Context, vehicle model.Vehicle) (model.Vehicle, error)
	Delete(ctx context.Context, id int64, userID int64) error
}

type vehicleRepository struct {
	db *pgxpool.Pool
}

func NewVehicleRepository(db *pgxpool.Pool) VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(ctx context.Context, vehicle model.Vehicle) (model.Vehicle, error) {
	err := r.db.QueryRow(ctx, `
		WITH ins AS (
			INSERT INTO vehicles (user_id, brand_id, model_id, license_plate, manufacturing_year, color, purchase_date, engine_number, current_mileage, status, image_url)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id, user_id, brand_id, model_id, license_plate, manufacturing_year, color, purchase_date, engine_number, current_mileage, status, image_url, created_at, updated_at
		)
		SELECT `+vehicleSelectColumns+`
		FROM ins v
		JOIN vehicle_brands b ON b.id = v.brand_id
		JOIN vehicle_models m ON m.id = v.model_id
	`, vehicle.UserID, vehicle.BrandID, vehicle.ModelID, vehicle.LicensePlate, vehicle.ManufacturingYear, vehicle.Color, vehicle.PurchaseDate, vehicle.EngineNumber, vehicle.CurrentMileage, vehicle.Status, vehicle.ImageURL).Scan(
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
	if err != nil {
		if vehicleForeignKeyError(err) != nil {
			return model.Vehicle{}, vehicleForeignKeyError(err)
		}
		return model.Vehicle{}, err
	}
	return vehicle, nil
}

func (r *vehicleRepository) FindAll(ctx context.Context, userID int64) ([]model.Vehicle, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+vehicleSelectColumns+`
		FROM vehicles v
		JOIN vehicle_brands b ON b.id = v.brand_id
		JOIN vehicle_models m ON m.id = v.model_id
		WHERE v.user_id = $1
		ORDER BY v.id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := make([]model.Vehicle, 0)
	for rows.Next() {
		var vehicle model.Vehicle
		if err := rows.Scan(
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
	); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, rows.Err()
}

func (r *vehicleRepository) FindByID(ctx context.Context, id int64, userID int64) (model.Vehicle, error) {
	var vehicle model.Vehicle
	err := r.db.QueryRow(ctx, `
		SELECT `+vehicleSelectColumns+`
		FROM vehicles v
		JOIN vehicle_brands b ON b.id = v.brand_id
		JOIN vehicle_models m ON m.id = v.model_id
		WHERE v.id = $1 AND v.user_id = $2
	`, id, userID).Scan(
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

func (r *vehicleRepository) Update(ctx context.Context, vehicle model.Vehicle) (model.Vehicle, error) {
	err := r.db.QueryRow(ctx, `
		WITH upd AS (
			UPDATE vehicles
			SET brand_id = $1, model_id = $2, license_plate = $3, manufacturing_year = $4,
				color = $5, purchase_date = $6, engine_number = $7, current_mileage = $8,
				status = $9, image_url = $10, updated_at = CURRENT_TIMESTAMP
			WHERE id = $11 AND user_id = $12
			RETURNING id, user_id, brand_id, model_id, license_plate, manufacturing_year, color, purchase_date, engine_number, current_mileage, status, image_url, created_at, updated_at
		)
		SELECT `+vehicleSelectColumns+`
		FROM upd v
		JOIN vehicle_brands b ON b.id = v.brand_id
		JOIN vehicle_models m ON m.id = v.model_id
	`, vehicle.BrandID, vehicle.ModelID, vehicle.LicensePlate, vehicle.ManufacturingYear, vehicle.Color, vehicle.PurchaseDate, vehicle.EngineNumber, vehicle.CurrentMileage, vehicle.Status, vehicle.ImageURL, vehicle.ID, vehicle.UserID).Scan(
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
		if vehicleForeignKeyError(err) != nil {
			return model.Vehicle{}, vehicleForeignKeyError(err)
		}
		return model.Vehicle{}, err
	}
	return vehicle, nil
}

func (r *vehicleRepository) Delete(ctx context.Context, id int64, userID int64) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM vehicles
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrVehicleNotFound
	}
	return nil
}

func vehicleForeignKeyError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23503" {
		return nil
	}
	switch pgErr.ConstraintName {
	case "vehicles_brand_id_fkey":
		return ErrVehicleBrandNotFoundFK
	case "vehicles_model_id_fkey":
		return ErrVehicleModelNotFoundFK
	default:
		return nil
	}
}
