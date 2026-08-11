package dto

type CreateVehicleRequest struct {
	BrandID           int64  `json:"brand_id" validate:"required,gt=0"`
	ModelID           int64  `json:"model_id" validate:"required,gt=0"`
	LicensePlate      string `json:"license_plate" validate:"required,min=1,max=30"`
	ManufacturingYear int    `json:"manufacturing_year" validate:"required,gt=1900,lte=2100"`
	Color             string `json:"color" validate:"omitempty,max=50"`
	PurchaseDate      string `json:"purchase_date" validate:"required,datetime=2006-01-02"`
	EngineNumber      string `json:"engine_number" validate:"required,min=1,max=30"`
	CurrentMileage    int    `json:"current_mileage" validate:"min=0"`
	Status            string `json:"status" validate:"omitempty,oneof=active inactive"`
	// image_url tidak dikirim dalam request json karena dikirim dalam bentuk file
}

type UpdateVehicleRequest struct {
	BrandID           int64  `json:"brand_id" validate:"required,gt=0"`
	ModelID           int64  `json:"model_id" validate:"required,gt=0"`
	LicensePlate      string `json:"license_plate" validate:"required,min=1,max=30"`
	ManufacturingYear int    `json:"manufacturing_year" validate:"required,gt=1900,lte=2100"`
	Color             string `json:"color" validate:"omitempty,max=50"`
	PurchaseDate      string `json:"purchase_date" validate:"required,datetime=2006-01-02"`
	EngineNumber      string `json:"engine_number" validate:"required,min=1,max=30"`
	CurrentMileage    int    `json:"current_mileage" validate:"min=0"`
	Status            string `json:"status" validate:"omitempty,oneof=active inactive"`
	// image_url tidak dikirim dalam request json karena dikirim dalam bentuk file
}
