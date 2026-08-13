package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/my-api/internal/model"
)

var ErrVehicleBrandNotFound = errors.New("brand kendaraan tidak ditemukan")
var ErrBrandNameExists = errors.New("nama brand sudah terdaftar")

type VehicleBrandRepository interface {
	Create(ctx context.Context, brand model.VehicleBrand) (model.VehicleBrand, error)
	FindAll(ctx context.Context) ([]model.VehicleBrand, error)
	FindByID(ctx context.Context, id int64) (model.VehicleBrand, error)
	Update(ctx context.Context, brand model.VehicleBrand) (model.VehicleBrand, error)
}

type vehicleBrandRepository struct {
	db *pgxpool.Pool
}

func NewVehicleBrandRepository(db *pgxpool.Pool) VehicleBrandRepository {
	return &vehicleBrandRepository{db: db}
}

func (r *vehicleBrandRepository) Create(ctx context.Context, brand model.VehicleBrand) (model.VehicleBrand, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO vehicle_brands (name, logo_url)
		VALUES ($1, $2)
		RETURNING id, name, logo_url, created_at, updated_at
	`, brand.Name, brand.LogoURL).Scan(
		&brand.ID,
		&brand.Name,
		&brand.LogoURL,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return model.VehicleBrand{}, ErrBrandNameExists
		}
		return model.VehicleBrand{}, err
	}
	return brand, nil
}

func (r *vehicleBrandRepository) FindAll(ctx context.Context) ([]model.VehicleBrand, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, logo_url, created_at, updated_at
		FROM vehicle_brands
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	brands := make([]model.VehicleBrand, 0)
	for rows.Next() {
		var brand model.VehicleBrand
		if err := rows.Scan(
			&brand.ID,
			&brand.Name,
			&brand.LogoURL,
			&brand.CreatedAt,
			&brand.UpdatedAt,
		); err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}

	return brands, rows.Err()
}

func (r *vehicleBrandRepository) FindByID(ctx context.Context, id int64) (model.VehicleBrand, error) {
	var brand model.VehicleBrand
	err := r.db.QueryRow(ctx, `
		SELECT id, name, logo_url, created_at, updated_at
		FROM vehicle_brands
		WHERE id = $1
	`, id).Scan(
		&brand.ID,
		&brand.Name,
		&brand.LogoURL,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.VehicleBrand{}, ErrVehicleBrandNotFound
	}
	if err != nil {
		return model.VehicleBrand{}, err
	}
	return brand, nil
}

func (r *vehicleBrandRepository) Update(ctx context.Context, brand model.VehicleBrand) (model.VehicleBrand, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE vehicle_brands
		SET name = $1, logo_url = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, name, logo_url, created_at, updated_at
	`, brand.Name, brand.LogoURL, brand.ID).Scan(
		&brand.ID,
		&brand.Name,
		&brand.LogoURL,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.VehicleBrand{}, ErrVehicleBrandNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return model.VehicleBrand{}, ErrBrandNameExists
		}
		return model.VehicleBrand{}, err
	}
	return brand, nil
}
