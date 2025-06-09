package dto

import "github.com/go-playground/validator/v10"

type CreateUserDto struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=32"`
	Phone     string `json:"phone" validate:"required"`
	Country   string `json:"country" validate:"required"`
}

func ValidateCreateUser(dto CreateUserDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
