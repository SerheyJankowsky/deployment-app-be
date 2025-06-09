package dto

import "github.com/go-playground/validator/v10"

type UpdateSecretDto struct {
	Name    *string `json:"name" validate:"omitempty,min=1,max=255"`
	Content *string `json:"content" validate:"omitempty,min=1,max=10000"`
}

func (dto *UpdateSecretDto) GetUpdates() (map[string]interface{}, []string) {
	updates := make(map[string]interface{})
	fields := make([]string, 0)
	if dto.Name != nil {
		updates["name"] = *dto.Name
		fields = append(fields, "name")
	}
	if dto.Content != nil {
		updates["content"] = *dto.Content
		fields = append(fields, "content")
	}
	return updates, fields
}

func (dto UpdateSecretDto) HasUpdates() bool {
	_, fields := dto.GetUpdates()
	return len(fields) > 0
}

func ValidateUpdateSecretDto(dto UpdateSecretDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
