package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/my-api/internal/model"
)

var ErrServiceItemNotFound = errors.New("item servis tidak ditemukan")
var ErrServiceRecordNotFoundFK = errors.New("record servis tidak ditemukan")

const serviceItemSelectColumns = `
	id, service_record_id, name, quantity, cost, notes, created_at, updated_at`

type ServiceItemRepository interface {
	Create(ctx context.Context, serviceItem model.ServiceItem) (model.ServiceItem, error)
	FindAll(ctx context.Context) ([]model.ServiceItem, error)
	FindByID(ctx context.Context, id int64) (model.ServiceItem, error)
	FindByServiceRecordIDs(ctx context.Context, recordIDs []int64) ([]model.ServiceItem, error)
	Update(ctx context.Context, serviceItem model.ServiceItem) (model.ServiceItem, error)
}

type serviceItemRepository struct {
	db *pgxpool.Pool
}

func NewServiceItemRepository(db *pgxpool.Pool) ServiceItemRepository {
	return &serviceItemRepository{db: db}
}

func (r *serviceItemRepository) Create(ctx context.Context, serviceItem model.ServiceItem) (model.ServiceItem, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO service_items (service_record_id, name, quantity, cost, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+serviceItemSelectColumns+`
	`, serviceItem.ServiceRecordID, serviceItem.Name, serviceItem.Quantity, serviceItem.Cost, serviceItem.Notes).Scan(
		&serviceItem.ID,
		&serviceItem.ServiceRecordID,
		&serviceItem.Name,
		&serviceItem.Quantity,
		&serviceItem.Cost,
		&serviceItem.Notes,
		&serviceItem.CreatedAt,
		&serviceItem.UpdatedAt,
	)
	if err != nil {
		if isForeignKeyViolation(err) {
			return model.ServiceItem{}, ErrServiceRecordNotFoundFK
		}
		return model.ServiceItem{}, err
	}
	return serviceItem, nil
}

func (r *serviceItemRepository) FindAll(ctx context.Context) ([]model.ServiceItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+serviceItemSelectColumns+`
		FROM service_items
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	serviceItems := make([]model.ServiceItem, 0)
	for rows.Next() {
		var serviceItem model.ServiceItem
		if err := rows.Scan(
			&serviceItem.ID,
			&serviceItem.ServiceRecordID,
			&serviceItem.Name,
			&serviceItem.Quantity,
			&serviceItem.Cost,
			&serviceItem.Notes,
			&serviceItem.CreatedAt,
			&serviceItem.UpdatedAt,
		); err != nil {
			return nil, err
		}
		serviceItems = append(serviceItems, serviceItem)
	}

	return serviceItems, rows.Err()
}

func (r *serviceItemRepository) FindByID(ctx context.Context, id int64) (model.ServiceItem, error) {
	var serviceItem model.ServiceItem
	err := r.db.QueryRow(ctx, `
		SELECT `+serviceItemSelectColumns+`
		FROM service_items
		WHERE id = $1
	`, id).Scan(
		&serviceItem.ID,
		&serviceItem.ServiceRecordID,
		&serviceItem.Name,
		&serviceItem.Quantity,
		&serviceItem.Cost,
		&serviceItem.Notes,
		&serviceItem.CreatedAt,
		&serviceItem.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ServiceItem{}, ErrServiceItemNotFound
	}
	if err != nil {
		return model.ServiceItem{}, err
	}
	return serviceItem, nil
}

func (r *serviceItemRepository) FindByServiceRecordIDs(ctx context.Context, recordIDs []int64) ([]model.ServiceItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+serviceItemSelectColumns+`
		FROM service_items
		WHERE service_record_id = ANY($1::bigint[])
		ORDER BY id ASC
	`, recordIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	serviceItems := make([]model.ServiceItem, 0)
	for rows.Next() {
		var serviceItem model.ServiceItem
		if err := rows.Scan(
			&serviceItem.ID,
			&serviceItem.ServiceRecordID,
			&serviceItem.Name,
			&serviceItem.Quantity,
			&serviceItem.Cost,
			&serviceItem.Notes,
			&serviceItem.CreatedAt,
			&serviceItem.UpdatedAt,
		); err != nil {
			return nil, err
		}
		serviceItems = append(serviceItems, serviceItem)
	}

	return serviceItems, rows.Err()
}

func (r *serviceItemRepository) Update(ctx context.Context, serviceItem model.ServiceItem) (model.ServiceItem, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE service_items
		SET service_record_id = $1, name = $2, quantity = $3, cost = $4, notes = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING `+serviceItemSelectColumns+`
	`, serviceItem.ServiceRecordID, serviceItem.Name, serviceItem.Quantity, serviceItem.Cost, serviceItem.Notes, serviceItem.ID).Scan(
		&serviceItem.ID,
		&serviceItem.ServiceRecordID,
		&serviceItem.Name,
		&serviceItem.Quantity,
		&serviceItem.Cost,
		&serviceItem.Notes,
		&serviceItem.CreatedAt,
		&serviceItem.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ServiceItem{}, ErrServiceItemNotFound
	}
	if err != nil {
		if isForeignKeyViolation(err) {
			return model.ServiceItem{}, ErrServiceRecordNotFoundFK
		}
		return model.ServiceItem{}, err
	}
	return serviceItem, nil
}
