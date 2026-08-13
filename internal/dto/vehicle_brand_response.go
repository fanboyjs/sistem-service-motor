package dto

type VehicleBrandResponse struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	LogoURL *string `json:"logo_url"`
}
