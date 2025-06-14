package dto

import "github.com/go-playground/validator/v10"

type UpdateDomainDto struct {
	Name    *string `json:"name" validate:"omitempty,min=1,max=255"`
	SSLCert *string `json:"ssl_cert" validate:"omitempty,min=0,max=10000"`
	SSLKey  *string `json:"ssl_key" validate:"omitempty,min=0,max=10000"`
}

func (dto *UpdateDomainDto) GetUpdatesDomain() (map[string]interface{}, []string) {
	updates := make(map[string]interface{})
	fields := make([]string, 0)
	if dto.Name != nil {
		updates["name"] = *dto.Name
		fields = append(fields, "name")
	}
	if dto.SSLCert != nil {
		updates["ssl_cert"] = *dto.SSLCert
		fields = append(fields, "ssl_cert")
	}
	if dto.SSLKey != nil {
		updates["ssl_key"] = *dto.SSLKey
		fields = append(fields, "ssl_key")
	}
	return updates, fields
}

func (dto UpdateDomainDto) HasUpdatesDomain() bool {
	_, fields := dto.GetUpdatesDomain()
	return len(fields) > 0
}

func ValidateUpdateDomainDto(dto UpdateDomainDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
