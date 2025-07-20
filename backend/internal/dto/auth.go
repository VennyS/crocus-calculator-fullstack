package dto

type RegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Permit   string `json:"permit"`
	FullName string `json:"full_name" validate:"required"`
}
