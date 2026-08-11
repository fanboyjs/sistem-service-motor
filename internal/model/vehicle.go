package model

import "time"

type Vehicle struct {
	ID                int64     `json:"id"`
	UserID            int64     `json:"user_id"`
	BrandID           int64     `json:"brand_id"`
	BrandName         string    `json:"brand_name"`
	ModelID           int64     `json:"model_id"`
	ModelName         string    `json:"model_name"`
	LicensePlate      string    `json:"license_plate"`
	ManufacturingYear int       `json:"manufacturing_year"`
	Color             *string   `json:"color"`
	PurchaseDate      time.Time `json:"purchase_date"`
	EngineNumber      string    `json:"engine_number"`
	CurrentMileage    int       `json:"current_mileage"`
	Status            string    `json:"status"`
	ImageURL          *string   `json:"image_url"`
	QRToken           string    `json:"-"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
}

type VehicleQRCode struct {
	ID        int64     `json:"id"`
	VehicleID int64     `json:"vehicle_id"`
	Token     string    `json:"token"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
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
	BrandName string    `json:"brand_name"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type VehicleTax struct {
	ID            int64     `json:"id"`
	VehicleID     int64     `json:"vehicle_id"`
	LicensePlate  string    `json:"license_plate"`
	TaxYear       int       `json:"tax_year"`
	PKBAmount     float64   `json:"pkb_amount"`
	SWDKLLJAmount float64   `json:"swdkllj_amount"`
	OtherAmount   float64   `json:"other_amount"`
	TotalAmount   float64   `json:"total_amount"`
	PaymentDate   time.Time `json:"payment_date"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"`
	Notes         *string   `json:"notes"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}
