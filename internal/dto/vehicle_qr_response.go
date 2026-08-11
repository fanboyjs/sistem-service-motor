package dto

type VehicleQRResponse struct {
	Token    string `json:"token"`
	ImageURL string `json:"qr_image_url"`
	ScanURL  string `json:"scan_url"`
}

type VehicleQRScanResponse struct {
	ID                int64   `json:"id"`
	BrandName         string  `json:"brand_name"`
	ModelName         string  `json:"model_name"`
	LicensePlate      string  `json:"license_plate"`
	ManufacturingYear int     `json:"manufacturing_year"`
	Color             *string `json:"color"`
	PurchaseDate      string  `json:"purchase_date"`
	CurrentMileage    int     `json:"current_mileage"`
	Status            string  `json:"status"`
	ImageURL          *string `json:"image_url"`
}
