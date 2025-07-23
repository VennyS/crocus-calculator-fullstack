package dto

type RegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Permit   string `json:"permit" validate:"required,oneof=manager all"`
	FullName string `json:"full_name" validate:"required"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
