package dto

type CreateVehicleBrandRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	// logo_url tidak dikirim dalam request json karena dikirim dalam bentuk file
}

type UpdateVehicleBrandRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	// logo_url tidak dikirim dalam request json karena dikirim dalam bentuk file
}
