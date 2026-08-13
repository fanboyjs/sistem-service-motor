package dto

type CreateServiceRecordRequest struct {
	VehicleID     int64   `json:"vehicle_id" validate:"required,gt=0"`
	ServiceTypeID int64   `json:"service_type_id" validate:"required,gt=0"`
	ServiceDate   string  `json:"service_date" validate:"required,datetime=2006-01-02"`
	Odometer      float64 `json:"odometer" validate:"required,gte=0"`
	WorkshopName  string  `json:"workshop_name" validate:"required,min=1,max=150"`
	LaborCost     float64 `json:"labor_cost" validate:"required,gte=0"`
	PartsCost     float64 `json:"parts_cost" validate:"required,gte=0"`
	Notes         string  `json:"notes" validate:"required"`
}

type UpdateServiceRecordRequest struct {
	VehicleID     int64   `json:"vehicle_id" validate:"required,gt=0"`
	ServiceTypeID int64   `json:"service_type_id" validate:"required,gt=0"`
	ServiceDate   string  `json:"service_date" validate:"required,datetime=2006-01-02"`
	Odometer      float64 `json:"odometer" validate:"required,gte=0"`
	WorkshopName  string  `json:"workshop_name" validate:"required,min=1,max=150"`
	LaborCost     float64 `json:"labor_cost" validate:"required,gte=0"`
	PartsCost     float64 `json:"parts_cost" validate:"required,gte=0"`
	Notes         string  `json:"notes" validate:"required"`
}
