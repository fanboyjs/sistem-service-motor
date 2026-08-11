package dto

type CreateVehicleTaxRequest struct {
	VehicleID     int64   `json:"vehicle_id" validate:"required,gt=0"`
	TaxYear       int     `json:"tax_year" validate:"required,gt=1900,lte=2100"`
	PKBAmount     float64 `json:"pkb_amount" validate:"required,gte=0"`
	SWDKLLJAmount float64 `json:"swdkllj_amount" validate:"required,gte=0"`
	OtherAmount   float64 `json:"other_amount" validate:"required,gte=0"`
	PaymentDate   string  `json:"payment_date" validate:"required,datetime=2006-01-02"`
	DueDate       string  `json:"due_date" validate:"required,datetime=2006-01-02"`
	Status        string  `json:"status" validate:"omitempty,oneof=paid unpaid"`
	Notes         *string `json:"notes" validate:"omitempty"`
}

type UpdateVehicleTaxRequest struct {
	VehicleID     int64   `json:"vehicle_id" validate:"required,gt=0"`
	TaxYear       int     `json:"tax_year" validate:"required,gt=1900,lte=2100"`
	PKBAmount     float64 `json:"pkb_amount" validate:"required,gte=0"`
	SWDKLLJAmount float64 `json:"swdkllj_amount" validate:"required,gte=0"`
	OtherAmount   float64 `json:"other_amount" validate:"required,gte=0"`
	PaymentDate   string  `json:"payment_date" validate:"required,datetime=2006-01-02"`
	DueDate       string  `json:"due_date" validate:"required,datetime=2006-01-02"`
	Status        string  `json:"status" validate:"omitempty,oneof=paid unpaid"`
	Notes         *string `json:"notes" validate:"omitempty"`
}
