package dto

type VehicleResponse struct {
	ID                int64   `json:"id"`
	UserID            int64   `json:"user_id"`
	BrandID           int64   `json:"brand_id"`
	BrandName         string  `json:"brand_name"`
	ModelID           int64   `json:"model_id"`
	ModelName         string  `json:"model_name"`
	LicensePlate      string  `json:"license_plate"`
	ManufacturingYear int     `json:"manufacturing_year"`
	Color             *string `json:"color"`
	PurchaseDate      string  `json:"purchase_date"`
	EngineNumber      string  `json:"engine_number"`
	CurrentMileage    int     `json:"current_mileage"`
	Status            string  `json:"status"`
	ImageURL          *string `json:"image_url"`
}
