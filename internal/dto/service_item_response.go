package dto

type ServiceItemResponse struct {
	ID              int64   `json:"id"`
	ServiceRecordID int64   `json:"service_record_id"`
	Name            string  `json:"name"`
	Quantity        int     `json:"quantity"`
	Cost            float64 `json:"cost"`
	Notes           string  `json:"notes"`
}
