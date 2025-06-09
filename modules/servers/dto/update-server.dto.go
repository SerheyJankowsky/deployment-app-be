package dto

import "github.com/go-playground/validator/v10"

type UpdateServerDto struct {
	Name     *string `json:"name" validate:"omitempty,min=1,max=255"`
	Host     *string `json:"host" validate:"omitempty,min=1,max=255"`
	Port     *int    `json:"port" validate:"omitempty,min=1,max=65535"`
	Username *string `json:"username" validate:"omitempty,min=1,max=255"`
	Password *string `json:"password" validate:"omitempty,min=1,max=255"`
	SSHKey   *string `json:"ssh_key" validate:"omitempty,min=1,max=10000"`
}

func (dto *UpdateServerDto) GetUpdates() (map[string]interface{}, []string) {
	updates := make(map[string]interface{})
	fields := make([]string, 0)
	if dto.Name != nil {
		updates["name"] = *dto.Name
		fields = append(fields, "name")
	}
	if dto.Host != nil {
		updates["host"] = *dto.Host
		fields = append(fields, "host")
	}
	if dto.Port != nil {
		updates["port"] = *dto.Port
		fields = append(fields, "port")
	}
	if dto.Username != nil {
		updates["username"] = *dto.Username
		fields = append(fields, "username")
	}
	if dto.Password != nil {
		updates["password"] = *dto.Password
		fields = append(fields, "password")
	}
	if dto.SSHKey != nil {
		updates["ssh_key"] = *dto.SSHKey
		fields = append(fields, "ssh_key")
	}
	return updates, fields
}

func (dto UpdateServerDto) HasUpdates() bool {
	_, fields := dto.GetUpdates()
	return len(fields) > 0
}

func ValidateUpdateServerDto(dto UpdateServerDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
