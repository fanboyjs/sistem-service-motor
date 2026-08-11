package dto

type VehicleModelResponse struct {
	ID        int64  `json:"id"`
	BrandID   int64  `json:"brand_id"`
	BrandName string `json:"brand_name"`
	Name      string `json:"name"`
}
