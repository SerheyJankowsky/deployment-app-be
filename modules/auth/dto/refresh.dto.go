package dto

import "github.com/go-playground/validator/v10"

type RefreshTokenDto struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func ValidateRefreshToken(dto RefreshTokenDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
