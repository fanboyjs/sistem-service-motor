package model

import "time"

type Vehicle struct {
	ID                int64     `json:"id"`
	UserID            int64     `json:"user_id"`
	BrandID           int64     `json:"brand_id"`
	ModelID           int64     `json:"model_id"`
	LicensePlate      string    `json:"license_plate"`
	ManufacturingYear int       `json:"manufacturing_year"`
	Color             *string   `json:"color"`
	PurchaseDate      time.Time `json:"purchase_date"`
	EngineNumber      string    `json:"engine_number"`
	CurrentMileage    int       `json:"current_mileage"`
	Status            string    `json:"status"`
	ImageURL          *string   `json:"image_url"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
}

type VehicleBrand struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	LogoURL   *string   `json:"logo_url"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type VehicleModel struct {
	ID        int64     `json:"id"`
	BrandID   int64     `json:"brand_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
