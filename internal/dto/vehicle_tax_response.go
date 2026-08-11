package dto

type VehicleTaxResponse struct {
	ID            int64   `json:"id"`
	VehicleID     int64   `json:"vehicle_id"`
	LicensePlate  string  `json:"license_plate"`
	TaxYear       int     `json:"tax_year"`
	PKBAmount     float64 `json:"pkb_amount"`
	SWDKLLJAmount float64 `json:"swdkllj_amount"`
	OtherAmount   float64 `json:"other_amount"`
	TotalAmount   float64 `json:"total_amount"`
	PaymentDate   string  `json:"payment_date"`
	DueDate       string  `json:"due_date"`
	Status        string  `json:"status"`
	Notes         *string `json:"notes"`
}
