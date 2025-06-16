package dto

import "github.com/go-playground/validator/v10"

type UpdateScriptDto struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Script      *string `json:"script" validate:"omitempty,min=1,max=10000"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1000"`
}

func (dto *UpdateScriptDto) GetUpdates() (map[string]interface{}, []string) {
	updates := make(map[string]interface{})
	fields := make([]string, 0)
	if dto.Name != nil {
		updates["name"] = *dto.Name
		fields = append(fields, "name")
	}
	if dto.Script != nil {
		updates["script"] = *dto.Script
		fields = append(fields, "script")
	}
	if dto.Description != nil {
		updates["description"] = *dto.Description
		fields = append(fields, "description")
	}
	return updates, fields
}

func (dto UpdateScriptDto) HasUpdates() bool {
	_, fields := dto.GetUpdates()
	return len(fields) > 0
}

func ValidateUpdateScriptDto(dto UpdateScriptDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
