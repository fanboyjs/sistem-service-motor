package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/my-api/internal/model"
)

var ErrServiceRecordNotFound = errors.New("record servis tidak ditemukan")
var ErrServiceTypeNotFoundFK = errors.New("jenis servis tidak ditemukan")

const serviceRecordSelectColumns = `
	sr.id, sr.vehicle_id, v.license_plate, sr.service_type_id, st.name, sr.service_date,
	sr.odometer, sr.workshop_name, sr.labor_cost, sr.parts_cost, sr.total_cost, sr.notes,
	sr.created_at, sr.updated_at`

type ServiceRecordRepository interface {
	Create(ctx context.Context, record model.ServiceRecord, userID int64) (model.ServiceRecord, error)
	FindAll(ctx context.Context, userID int64) ([]model.ServiceRecord, error)
	FindByID(ctx context.Context, id int64, userID int64) (model.ServiceRecord, error)
	Update(ctx context.Context, record model.ServiceRecord, userID int64) (model.ServiceRecord, error)
	Delete(ctx context.Context, id int64, userID int64) error
}

type serviceRecordRepository struct {
	db *pgxpool.Pool
}

func NewServiceRecordRepository(db *pgxpool.Pool) ServiceRecordRepository {
	return &serviceRecordRepository{db: db}
}

func (r *serviceRecordRepository) Create(ctx context.Context, record model.ServiceRecord, userID int64) (model.ServiceRecord, error) {
	err := r.db.QueryRow(ctx, `
		WITH ins AS (
			INSERT INTO service_records (vehicle_id, service_type_id, service_date, odometer, workshop_name, labor_cost, parts_cost, total_cost, notes)
			SELECT v.id, $2, $3, $4, $5, $6, $7, $8, $9
			FROM vehicles v
			WHERE v.id = $1 AND v.user_id = $10
				AND EXISTS (SELECT 1 FROM service_types WHERE id = $2)
			RETURNING id, vehicle_id, service_type_id, service_date, odometer, workshop_name, labor_cost, parts_cost, total_cost, notes, created_at, updated_at
		)
		SELECT `+serviceRecordSelectColumns+`
		FROM ins sr
		JOIN vehicles v ON v.id = sr.vehicle_id
		JOIN service_types st ON st.id = sr.service_type_id
	`, record.VehicleID, record.ServiceTypeID, record.ServiceDate, record.Odometer, record.WorkshopName, record.LaborCost, record.PartsCost, record.TotalCost, record.Notes, userID).Scan(
		&record.ID,
		&record.VehicleID,
		&record.LicensePlate,
		&record.ServiceTypeID,
		&record.ServiceTypeName,
		&record.ServiceDate,
		&record.Odometer,
		&record.WorkshopName,
		&record.LaborCost,
		&record.PartsCost,
		&record.TotalCost,
		&record.Notes,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ServiceRecord{}, r.checkCreateFK(ctx, record.VehicleID, record.ServiceTypeID, userID)
	}
	if err != nil {
		return model.ServiceRecord{}, err
	}
	return record, nil
}

func (r *serviceRecordRepository) FindAll(ctx context.Context, userID int64) ([]model.ServiceRecord, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+serviceRecordSelectColumns+`
		FROM service_records sr
		JOIN vehicles v ON v.id = sr.vehicle_id
		JOIN service_types st ON st.id = sr.service_type_id
		WHERE v.user_id = $1
		ORDER BY sr.id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]model.ServiceRecord, 0)
	for rows.Next() {
		var record model.ServiceRecord
		if err := rows.Scan(
			&record.ID,
			&record.VehicleID,
			&record.LicensePlate,
			&record.ServiceTypeID,
			&record.ServiceTypeName,
			&record.ServiceDate,
			&record.Odometer,
			&record.WorkshopName,
			&record.LaborCost,
			&record.PartsCost,
			&record.TotalCost,
			&record.Notes,
			&record.CreatedAt,
			&record.UpdatedAt,
		); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

func (r *serviceRecordRepository) FindByID(ctx context.Context, id int64, userID int64) (model.ServiceRecord, error) {
	var record model.ServiceRecord
	err := r.db.QueryRow(ctx, `
		SELECT `+serviceRecordSelectColumns+`
		FROM service_records sr
		JOIN vehicles v ON v.id = sr.vehicle_id
		JOIN service_types st ON st.id = sr.service_type_id
		WHERE sr.id = $1 AND v.user_id = $2
	`, id, userID).Scan(
		&record.ID,
		&record.VehicleID,
		&record.LicensePlate,
		&record.ServiceTypeID,
		&record.ServiceTypeName,
		&record.ServiceDate,
		&record.Odometer,
		&record.WorkshopName,
		&record.LaborCost,
		&record.PartsCost,
		&record.TotalCost,
		&record.Notes,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ServiceRecord{}, ErrServiceRecordNotFound
	}
	if err != nil {
		return model.ServiceRecord{}, err
	}
	return record, nil
}

func (r *serviceRecordRepository) Update(ctx context.Context, record model.ServiceRecord, userID int64) (model.ServiceRecord, error) {
	err := r.db.QueryRow(ctx, `
		WITH upd AS (
			UPDATE service_records
			SET vehicle_id = $1, service_type_id = $2, service_date = $3, odometer = $4, workshop_name = $5,
				labor_cost = $6, parts_cost = $7, total_cost = $8, notes = $9, updated_at = CURRENT_TIMESTAMP
			WHERE id = $10 AND EXISTS (
				SELECT 1 FROM vehicles WHERE id = $1 AND user_id = $11
			) AND EXISTS (
				SELECT 1 FROM service_types WHERE id = $2
			)
			RETURNING id, vehicle_id, service_type_id, service_date, odometer, workshop_name, labor_cost, parts_cost, total_cost, notes, created_at, updated_at
		)
		SELECT `+serviceRecordSelectColumns+`
		FROM upd sr
		JOIN vehicles v ON v.id = sr.vehicle_id
		JOIN service_types st ON st.id = sr.service_type_id
	`, record.VehicleID, record.ServiceTypeID, record.ServiceDate, record.Odometer, record.WorkshopName, record.LaborCost, record.PartsCost, record.TotalCost, record.Notes, record.ID, userID).Scan(
		&record.ID,
		&record.VehicleID,
		&record.LicensePlate,
		&record.ServiceTypeID,
		&record.ServiceTypeName,
		&record.ServiceDate,
		&record.Odometer,
		&record.WorkshopName,
		&record.LaborCost,
		&record.PartsCost,
		&record.TotalCost,
		&record.Notes,
		&record.CreatedAt,
		&record.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ServiceRecord{}, r.checkUpdateFK(ctx, record.VehicleID, record.ServiceTypeID, record.ID, userID)
	}
	if err != nil {
		return model.ServiceRecord{}, err
	}
	return record, nil
}

func (r *serviceRecordRepository) Delete(ctx context.Context, id int64, userID int64) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM service_records sr
		USING vehicles v
		WHERE sr.id = $1 AND sr.vehicle_id = v.id AND v.user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrServiceRecordNotFound
	}
	return nil
}

func (r *serviceRecordRepository) checkCreateFK(ctx context.Context, vehicleID, serviceTypeID, userID int64) error {
	var vehicleExists bool
	if err := r.db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM vehicles WHERE id = $1 AND user_id = $2)
	`, vehicleID, userID).Scan(&vehicleExists); err != nil {
		return err
	}
	if !vehicleExists {
		return ErrVehicleNotFoundFK
	}
	return r.checkServiceTypeExists(ctx, serviceTypeID)
}

func (r *serviceRecordRepository) checkUpdateFK(ctx context.Context, vehicleID, serviceTypeID, recordID, userID int64) error {
	var vehicleExists bool
	if err := r.db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM vehicles WHERE id = $1 AND user_id = $2)
	`, vehicleID, userID).Scan(&vehicleExists); err != nil {
		return err
	}
	if !vehicleExists {
		return ErrVehicleNotFoundFK
	}
	if err := r.checkServiceTypeExists(ctx, serviceTypeID); err != nil {
		return err
	}
	return ErrServiceRecordNotFound
}

func (r *serviceRecordRepository) checkServiceTypeExists(ctx context.Context, serviceTypeID int64) error {
	var serviceTypeExists bool
	if err := r.db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM service_types WHERE id = $1)
	`, serviceTypeID).Scan(&serviceTypeExists); err != nil {
		return err
	}
	if !serviceTypeExists {
		return ErrServiceTypeNotFoundFK
	}
	return nil
}
