package dto

import "github.com/go-playground/validator/v10"

type CreateServerDto struct {
	Name     string  `json:"name" validate:"required,min=1,max=255"`
	Host     string  `json:"host" validate:"required,min=1,max=255"`
	Port     int     `json:"port" validate:"omitempty,min=1,max=65535"`
	Username string  `json:"username" validate:"required,min=1,max=255"`
	Password string  `json:"password" validate:"required,min=1,max=255"`
	SSHKey   *string `json:"ssh_key" validate:"omitempty,min=1,max=10000"`
}

func ValidateCreateServerDto(dto CreateServerDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
