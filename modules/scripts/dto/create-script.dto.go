package dto

import "github.com/go-playground/validator/v10"

type CreateScriptDto struct {
	Name   string `json:"name" validate:"required,min=1,max=255"`
	Script string `json:"script" validate:"required,min=1,max=10000"`
}

func ValidateCreateScriptDto(dto CreateScriptDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
