package dto

import "github.com/go-playground/validator/v10"

type UpdateSubDomainDto struct {
	Name *string `json:"name" validate:"omitempty,min=1,max=255"`
}

func (dto *UpdateSubDomainDto) GetUpdatesSubDomain() (map[string]interface{}, []string) {
	updates := make(map[string]interface{})
	fields := make([]string, 0)
	if dto.Name != nil {
		updates["name"] = *dto.Name
		fields = append(fields, "name")
	}
	return updates, fields
}

func (dto UpdateSubDomainDto) HasUpdatesSubDomain() bool {
	_, fields := dto.GetUpdatesSubDomain()
	return len(fields) > 0
}

func ValidateUpdateSubDomainDto(dto UpdateSubDomainDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
