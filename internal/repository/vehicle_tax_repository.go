package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example.com/my-api/internal/model"
)

var ErrVehicleTaxNotFound = errors.New("pajak kendaraan tidak ditemukan")
var ErrVehicleNotFoundFK = errors.New("kendaraan tidak ditemukan")

const vehicleTaxSelectColumns = `
	vt.id, vt.vehicle_id, v.license_plate, vt.tax_year, vt.pkb_amount, vt.swdkllj_amount,
	vt.other_amount, vt.total_amount, vt.payment_date, vt.due_date, vt.status, vt.notes,
	vt.created_at, vt.updated_at`

type VehicleTaxRepository interface {
	Create(ctx context.Context, tax model.VehicleTax, userID int64) (model.VehicleTax, error)
	FindAll(ctx context.Context, userID int64) ([]model.VehicleTax, error)
	FindByID(ctx context.Context, id int64, userID int64) (model.VehicleTax, error)
	Update(ctx context.Context, tax model.VehicleTax, userID int64) (model.VehicleTax, error)
	Delete(ctx context.Context, id int64, userID int64) error
}

type vehicleTaxRepository struct {
	db *pgxpool.Pool
}

func NewVehicleTaxRepository(db *pgxpool.Pool) VehicleTaxRepository {
	return &vehicleTaxRepository{db: db}
}

func (r *vehicleTaxRepository) Create(ctx context.Context, tax model.VehicleTax, userID int64) (model.VehicleTax, error) {
	err := r.db.QueryRow(ctx, `
		WITH ins AS (
			INSERT INTO vehicle_taxes (vehicle_id, tax_year, pkb_amount, swdkllj_amount, other_amount, total_amount, payment_date, due_date, status, notes)
			SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
			FROM vehicles
			WHERE id = $1 AND user_id = $11
			RETURNING id, vehicle_id, tax_year, pkb_amount, swdkllj_amount, other_amount, total_amount, payment_date, due_date, status, notes, created_at, updated_at
		)
		SELECT `+vehicleTaxSelectColumns+`
		FROM ins vt
		JOIN vehicles v ON v.id = vt.vehicle_id
	`, tax.VehicleID, tax.TaxYear, tax.PKBAmount, tax.SWDKLLJAmount, tax.OtherAmount, tax.TotalAmount, tax.PaymentDate, tax.DueDate, tax.Status, tax.Notes, userID).Scan(
		&tax.ID,
		&tax.VehicleID,
		&tax.LicensePlate,
		&tax.TaxYear,
		&tax.PKBAmount,
		&tax.SWDKLLJAmount,
		&tax.OtherAmount,
		&tax.TotalAmount,
		&tax.PaymentDate,
		&tax.DueDate,
		&tax.Status,
		&tax.Notes,
		&tax.CreatedAt,
		&tax.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.VehicleTax{}, ErrVehicleNotFoundFK
	}
	if err != nil {
		return model.VehicleTax{}, err
	}
	return tax, nil
}

func (r *vehicleTaxRepository) FindAll(ctx context.Context, userID int64) ([]model.VehicleTax, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+vehicleTaxSelectColumns+`
		FROM vehicle_taxes vt
		JOIN vehicles v ON v.id = vt.vehicle_id
		WHERE v.user_id = $1
		ORDER BY vt.id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taxes := make([]model.VehicleTax, 0)
	for rows.Next() {
		var tax model.VehicleTax
		if err := rows.Scan(
			&tax.ID,
			&tax.VehicleID,
			&tax.LicensePlate,
			&tax.TaxYear,
			&tax.PKBAmount,
			&tax.SWDKLLJAmount,
			&tax.OtherAmount,
			&tax.TotalAmount,
			&tax.PaymentDate,
			&tax.DueDate,
			&tax.Status,
			&tax.Notes,
			&tax.CreatedAt,
			&tax.UpdatedAt,
		); err != nil {
			return nil, err
		}
		taxes = append(taxes, tax)
	}

	return taxes, rows.Err()
}

func (r *vehicleTaxRepository) FindByID(ctx context.Context, id int64, userID int64) (model.VehicleTax, error) {
	var tax model.VehicleTax
	err := r.db.QueryRow(ctx, `
		SELECT `+vehicleTaxSelectColumns+`
		FROM vehicle_taxes vt
		JOIN vehicles v ON v.id = vt.vehicle_id
		WHERE vt.id = $1 AND v.user_id = $2
	`, id, userID).Scan(
		&tax.ID,
		&tax.VehicleID,
		&tax.LicensePlate,
		&tax.TaxYear,
		&tax.PKBAmount,
		&tax.SWDKLLJAmount,
		&tax.OtherAmount,
		&tax.TotalAmount,
		&tax.PaymentDate,
		&tax.DueDate,
		&tax.Status,
		&tax.Notes,
		&tax.CreatedAt,
		&tax.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.VehicleTax{}, ErrVehicleTaxNotFound
	}
	if err != nil {
		return model.VehicleTax{}, err
	}
	return tax, nil
}

func (r *vehicleTaxRepository) Update(ctx context.Context, tax model.VehicleTax, userID int64) (model.VehicleTax, error) {
	err := r.db.QueryRow(ctx, `
		WITH upd AS (
			UPDATE vehicle_taxes
			SET vehicle_id = $1, tax_year = $2, pkb_amount = $3, swdkllj_amount = $4, other_amount = $5,
				total_amount = $6, payment_date = $7, due_date = $8, status = $9, notes = $10,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $11 AND EXISTS (
				SELECT 1 FROM vehicles WHERE id = $1 AND user_id = $12
			)
			RETURNING id, vehicle_id, tax_year, pkb_amount, swdkllj_amount, other_amount, total_amount, payment_date, due_date, status, notes, created_at, updated_at
		)
		SELECT `+vehicleTaxSelectColumns+`
		FROM upd vt
		JOIN vehicles v ON v.id = vt.vehicle_id
	`, tax.VehicleID, tax.TaxYear, tax.PKBAmount, tax.SWDKLLJAmount, tax.OtherAmount, tax.TotalAmount, tax.PaymentDate, tax.DueDate, tax.Status, tax.Notes, tax.ID, userID).Scan(
		&tax.ID,
		&tax.VehicleID,
		&tax.LicensePlate,
		&tax.TaxYear,
		&tax.PKBAmount,
		&tax.SWDKLLJAmount,
		&tax.OtherAmount,
		&tax.TotalAmount,
		&tax.PaymentDate,
		&tax.DueDate,
		&tax.Status,
		&tax.Notes,
		&tax.CreatedAt,
		&tax.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		var exists bool
		if checkErr := r.db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM vehicles WHERE id = $1 AND user_id = $2
			)
		`, tax.VehicleID, userID).Scan(&exists); checkErr != nil {
			return model.VehicleTax{}, checkErr
		}
		if !exists {
			return model.VehicleTax{}, ErrVehicleNotFoundFK
		}
		return model.VehicleTax{}, ErrVehicleTaxNotFound
	}
	if err != nil {
		return model.VehicleTax{}, err
	}
	return tax, nil
}

func (r *vehicleTaxRepository) Delete(ctx context.Context, id int64, userID int64) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM vehicle_taxes vt
		USING vehicles v
		WHERE vt.id = $1 AND vt.vehicle_id = v.id AND v.user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrVehicleTaxNotFound
	}
	return nil
}
