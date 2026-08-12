package model

import "time"

type ServiceRecord struct {
	ID              int64     `json:"id"`
	VehicleID       int64     `json:"vehicle_id"`
	LicensePlate    string    `json:"license_plate"`
	ServiceTypeID   int64     `json:"service_type_id"`
	ServiceTypeName string    `json:"service_type_name"`
	ServiceDate     time.Time `json:"service_date"`
	Odometer        float64   `json:"odometer"`
	WorkshopName    string    `json:"workshop_name"`
	LaborCost       float64   `json:"labor_cost"`
	PartsCost       float64   `json:"parts_cost"`
	TotalCost       float64   `json:"total_cost"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}
