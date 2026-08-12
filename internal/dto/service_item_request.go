package dto

type CreateServiceItemRequest struct {
	ServiceRecordID int64   `json:"service_record_id" validate:"required"`
	Name            string  `json:"name" validate:"required,min=1,max=150"`
	Quantity        int     `json:"quantity" validate:"required,min=1"`
	Cost            float64 `json:"cost" validate:"required,gt=0"`
	Notes           string  `json:"notes" validate:"required"`
}

type UpdateServiceItemRequest struct {
	ServiceRecordID int64   `json:"service_record_id" validate:"required"`
	Name            string  `json:"name" validate:"required,min=1,max=150"`
	Quantity        int     `json:"quantity" validate:"required,min=1"`
	Cost            float64 `json:"cost" validate:"required,gt=0"`
	Notes           string  `json:"notes" validate:"required"`
}
