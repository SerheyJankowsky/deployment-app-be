package dto

import (
	"deployer.com/modules/users/dto"
	"github.com/go-playground/validator/v10"
)

type RegisterDto struct {
	IsRememberMe bool              `json:"is_remember_me"`
	User         dto.CreateUserDto `json:"user"`
}

func ValidateRegister(dto RegisterDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
