package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/my-api/internal/model"
)

var ErrVehicleModelNotFound = errors.New("model kendaraan tidak ditemukan")
var ErrModelNameExists = errors.New("nama model kendaraan sudah terdaftar")
var ErrVehicleBrandNotFoundFK = errors.New("brand kendaraan tidak ditemukan")

const vehicleModelSelectColumns = `
	vm.id, vm.brand_id, b.name, vm.name,
	vm.created_at, vm.updated_at`

type VehicleModelRepository interface {
	Create(ctx context.Context, modelData model.VehicleModel) (model.VehicleModel, error)
	FindAll(ctx context.Context) ([]model.VehicleModel, error)
	FindByID(ctx context.Context, id int64) (model.VehicleModel, error)
	Update(ctx context.Context, modelData model.VehicleModel) (model.VehicleModel, error)
}

type vehicleModelRepository struct {
	db *pgxpool.Pool
}

func NewVehicleModelRepository(db *pgxpool.Pool) VehicleModelRepository {
	return &vehicleModelRepository{db: db}
}

func (r *vehicleModelRepository) Create(ctx context.Context, modelData model.VehicleModel) (model.VehicleModel, error) {
	err := r.db.QueryRow(ctx, `
		WITH ins AS (
			INSERT INTO vehicle_models (brand_id, name)
			VALUES ($1, $2)
			RETURNING id, brand_id, name, created_at, updated_at
		)
		SELECT `+vehicleModelSelectColumns+`
		FROM ins vm
		JOIN vehicle_brands b ON b.id = vm.brand_id
	`, modelData.BrandID, modelData.Name).Scan(
		&modelData.ID,
		&modelData.BrandID,
		&modelData.BrandName,
		&modelData.Name,
		&modelData.CreatedAt,
		&modelData.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return model.VehicleModel{}, ErrModelNameExists
		}
		if isForeignKeyViolation(err) {
			return model.VehicleModel{}, ErrVehicleBrandNotFoundFK
		}
		return model.VehicleModel{}, err
	}
	return modelData, nil
}

func (r *vehicleModelRepository) FindAll(ctx context.Context) ([]model.VehicleModel, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+vehicleModelSelectColumns+`
		FROM vehicle_models vm
		JOIN vehicle_brands b ON b.id = vm.brand_id
		ORDER BY vm.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models := make([]model.VehicleModel, 0)
	for rows.Next() {
		var modelData model.VehicleModel
		if err := rows.Scan(
			&modelData.ID,
			&modelData.BrandID,
			&modelData.BrandName,
			&modelData.Name,
			&modelData.CreatedAt,
			&modelData.UpdatedAt,
		); err != nil {
			return nil, err
		}
		models = append(models, modelData)
	}

	return models, rows.Err()
}

func (r *vehicleModelRepository) FindByID(ctx context.Context, id int64) (model.VehicleModel, error) {
	var modelData model.VehicleModel
	err := r.db.QueryRow(ctx, `
		SELECT `+vehicleModelSelectColumns+`
		FROM vehicle_models vm
		JOIN vehicle_brands b ON b.id = vm.brand_id
		WHERE vm.id = $1
	`, id).Scan(
		&modelData.ID,
		&modelData.BrandID,
		&modelData.BrandName,
		&modelData.Name,
		&modelData.CreatedAt,
		&modelData.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.VehicleModel{}, ErrVehicleModelNotFound
	}
	if err != nil {
		return model.VehicleModel{}, err
	}
	return modelData, nil
}

func (r *vehicleModelRepository) Update(ctx context.Context, modelData model.VehicleModel) (model.VehicleModel, error) {
	err := r.db.QueryRow(ctx, `
		WITH upd AS (
			UPDATE vehicle_models
			SET brand_id = $1, name = $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $3
			RETURNING id, brand_id, name, created_at, updated_at
		)
		SELECT `+vehicleModelSelectColumns+`
		FROM upd vm
		JOIN vehicle_brands b ON b.id = vm.brand_id
	`, modelData.BrandID, modelData.Name, modelData.ID).Scan(
		&modelData.ID,
		&modelData.BrandID,
		&modelData.BrandName,
		&modelData.Name,
		&modelData.CreatedAt,
		&modelData.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.VehicleModel{}, ErrVehicleModelNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return model.VehicleModel{}, ErrModelNameExists
		}
		if isForeignKeyViolation(err) {
			return model.VehicleModel{}, ErrVehicleBrandNotFoundFK
		}
		return model.VehicleModel{}, err
	}
	return modelData, nil
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
