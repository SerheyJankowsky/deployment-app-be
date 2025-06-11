package dto

import "github.com/go-playground/validator/v10"

type CreateSubDomainDto struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	DomainID uint   `json:"domain_id" validate:"required,min=1"`
}

func ValidateCreateSubDomainDto(dto CreateSubDomainDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
