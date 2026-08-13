package dto

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Email string `json:"email" validate:"required,email,max=255"`
	Phone string `json:"phone" validate:"omitempty,max=20"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Email string `json:"email" validate:"required,email,max=255"`
	Phone string `json:"phone" validate:"omitempty,max=20"`
}
