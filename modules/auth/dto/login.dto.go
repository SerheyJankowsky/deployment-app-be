package dto

import "github.com/go-playground/validator/v10"

type LoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func ValidateLogin(dto LoginDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
