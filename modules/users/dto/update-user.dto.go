package dto

import "github.com/go-playground/validator/v10"

type UpdateUserDto struct {
	ID        uint   `json:"id" validate:"required"`
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name" validate:"omitempty"`
	// Email     string `json:"email" validate:"omitempty,email"`
	Phone   string `json:"phone" validate:"omitempty"`
	Country string `json:"country" validate:"omitempty"`
	// City         string    `gorm:"not null" json:"city"`
	// Address      string    `gorm:"not null" json:"address"`
	// ZipCode      string    `gorm:"not null" json:"zip_code"`
}

func ValidateUpdateUser(dto UpdateUserDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
