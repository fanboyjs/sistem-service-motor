package dto

type CreateVehicleModelRequest struct {
	BrandID int64  `json:"brand_id" validate:"required,gt=0"`
	Name    string `json:"name" validate:"required,min=1,max=255"`
}

type UpdateVehicleModelRequest struct {
	BrandID int64  `json:"brand_id" validate:"required,gt=0"`
	Name    string `json:"name" validate:"required,min=1,max=255"`
}
