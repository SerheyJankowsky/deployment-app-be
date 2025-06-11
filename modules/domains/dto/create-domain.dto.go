package dto

import "github.com/go-playground/validator/v10"

type CreateDomainDto struct {
	Name    string `json:"name" validate:"required,min=1,max=255"`
	SSLCert string `json:"ssl_cert" validate:"required,min=1,max=10000"`
	SSLKey  string `json:"ssl_key" validate:"required,min=1,max=10000"`
}

func ValidateCreateDomainDto(dto CreateDomainDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
