package dto

type VehicleModelResponse struct {
	ID      int64  `json:"id"`
	BrandID int64  `json:"brand_id"`
	Name    string `json:"name"`
}
