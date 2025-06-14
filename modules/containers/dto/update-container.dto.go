package dto

import "github.com/go-playground/validator/v10"

type UpdateContainerDto struct {
	Name      *string `json:"name" validate:"omitempty,min=1,max=255"`
	Registry  *string `json:"registry" validate:"omitempty,min=0,max=255"`
	Image     *string `json:"image" validate:"omitempty,min=1,max=255"`
	Tag       *string `json:"tag" validate:"omitempty,min=1,max=255"`
	Username  *string `json:"username" validate:"omitempty,min=0,max=255"`
	Password  *string `json:"password" validate:"omitempty,min=0,max=255"`
	SecretKey *string `json:"secret_key" validate:"omitempty,min=0,max=255"`
	Params    *string `json:"params" validate:"omitempty,min=0,max=10000"`
}

func (dto *UpdateContainerDto) GetUpdates() (map[string]interface{}, []string) {
	updates := make(map[string]interface{})
	fields := make([]string, 0)
	if dto.Name != nil {
		updates["name"] = *dto.Name
		fields = append(fields, "name")
	}
	if dto.Registry != nil {
		updates["registry"] = *dto.Registry
		fields = append(fields, "registry")
	}
	if dto.Image != nil {
		updates["image"] = *dto.Image
		fields = append(fields, "image")
	}
	if dto.Tag != nil {
		updates["tag"] = *dto.Tag
		fields = append(fields, "tag")
	}
	if dto.Username != nil {
		updates["username"] = *dto.Username
		fields = append(fields, "username")
	}
	if dto.Password != nil {
		updates["password"] = *dto.Password
		fields = append(fields, "password")
	}
	if dto.SecretKey != nil {
		updates["secret_key"] = *dto.SecretKey
		fields = append(fields, "secret_key")
	}
	if dto.Params != nil {
		updates["params"] = *dto.Params
		fields = append(fields, "params")
	}

	return updates, fields
}

func (dto UpdateContainerDto) HasUpdates() bool {
	_, fields := dto.GetUpdates()
	return len(fields) > 0
}

func ValidateUpdateContainerDto(dto UpdateContainerDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
