package dto

import "github.com/go-playground/validator/v10"

type CreateContainerDto struct {
	Name      string `json:"name" validate:"required,min=1,max=255"`
	Registry  string `json:"registry" validate:"omitempty,min=1,max=255"`
	Image     string `json:"image" validate:"required,min=1,max=255"`
	Tag       string `json:"tag" validate:"omitempty,min=1,max=255"`
	Username  string `json:"username" validate:"omitempty,min=0,max=255"`
	Password  string `json:"password" validate:"omitempty,min=0,max=255"`
	SecretKey string `json:"secret_key" validate:"omitempty,min=0,max=255"`
	Params    string `json:"params" validate:"omitempty,min=0,max=10000"`
}

func ValidateCreateContainerDto(dto CreateContainerDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
