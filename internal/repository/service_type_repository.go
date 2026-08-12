package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/my-api/internal/model"
)

var ErrServiceTypeNotFound = errors.New("tipe servis tidak ditemukan")
var ErrServiceTypeNameExists = errors.New("nama tipe servis sudah terdaftar")

const serviceTypeSelectColumns = `
	id, name, description, created_at, updated_at`

type ServiceTypeRepository interface {
	Create(ctx context.Context, serviceType model.ServiceType) (model.ServiceType, error)
	FindAll(ctx context.Context) ([]model.ServiceType, error)
	FindByID(ctx context.Context, id int64) (model.ServiceType, error)
	Update(ctx context.Context, serviceType model.ServiceType) (model.ServiceType, error)
}

type serviceTypeRepository struct {
	db *pgxpool.Pool
}

func NewServiceTypeRepository(db *pgxpool.Pool) ServiceTypeRepository {
	return &serviceTypeRepository{db: db}
}

func (r *serviceTypeRepository) Create(ctx context.Context, serviceType model.ServiceType) (model.ServiceType, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO service_types (name, description)
		VALUES ($1, $2)
		RETURNING `+serviceTypeSelectColumns+`
	`, serviceType.Name, serviceType.Description).Scan(
		&serviceType.ID,
		&serviceType.Name,
		&serviceType.Description,
		&serviceType.CreatedAt,
		&serviceType.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return model.ServiceType{}, ErrServiceTypeNameExists
		}
		return model.ServiceType{}, err
	}
	return serviceType, nil
}

func (r *serviceTypeRepository) FindAll(ctx context.Context) ([]model.ServiceType, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+serviceTypeSelectColumns+`
		FROM service_types
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	serviceTypes := make([]model.ServiceType, 0)
	for rows.Next() {
		var serviceType model.ServiceType
		if err := rows.Scan(
			&serviceType.ID,
			&serviceType.Name,
			&serviceType.Description,
			&serviceType.CreatedAt,
			&serviceType.UpdatedAt,
		); err != nil {
			return nil, err
		}
		serviceTypes = append(serviceTypes, serviceType)
	}

	return serviceTypes, rows.Err()
}

func (r *serviceTypeRepository) FindByID(ctx context.Context, id int64) (model.ServiceType, error) {
	var serviceType model.ServiceType
	err := r.db.QueryRow(ctx, `
		SELECT `+serviceTypeSelectColumns+`
		FROM service_types
		WHERE id = $1
	`, id).Scan(
		&serviceType.ID,
		&serviceType.Name,
		&serviceType.Description,
		&serviceType.CreatedAt,
		&serviceType.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ServiceType{}, ErrServiceTypeNotFound
	}
	if err != nil {
		return model.ServiceType{}, err
	}
	return serviceType, nil
}

func (r *serviceTypeRepository) Update(ctx context.Context, serviceType model.ServiceType) (model.ServiceType, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE service_types
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING `+serviceTypeSelectColumns+`
	`, serviceType.Name, serviceType.Description, serviceType.ID).Scan(
		&serviceType.ID,
		&serviceType.Name,
		&serviceType.Description,
		&serviceType.CreatedAt,
		&serviceType.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ServiceType{}, ErrServiceTypeNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return model.ServiceType{}, ErrServiceTypeNameExists
		}
		return model.ServiceType{}, err
	}
	return serviceType, nil
}
