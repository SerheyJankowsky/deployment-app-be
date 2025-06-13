package dto

import (
	"reflect"

	// "deployer.com/modules/servers" // Commented out due to import issue

	"github.com/go-playground/validator/v10"
)

type UpdateProjectDto struct {
	Name               *string              `json:"name" validate:"omitempty,min=1,max=255" db:"name"`
	ProjectDeployments []ProjectDeployments `json:"project_deployments" validate:"omitempty,dive" db:"project_deployments"`
}

func (dto *UpdateProjectDto) GetUpdates() (map[string]interface{}, []string) {
	updates := make(map[string]interface{})
	fields := make([]string, 0)

	v := reflect.ValueOf(dto).Elem()
	t := reflect.TypeOf(dto).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Получаем имя поля из db тега
		dbTag := fieldType.Tag.Get("db")
		if dbTag == "" {
			continue
		}

		// Обрабатываем указатели
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			updates[dbTag] = field.Elem().Interface()
			fields = append(fields, dbTag)
		}

		// Обрабатываем слайсы (nil означает "не обновлять")
		if field.Kind() == reflect.Slice && !field.IsNil() {
			// Handle ID slices vs object slices differently
			fieldName := fieldType.Name
			if fieldName == "ProjectDeployments" {
				// Convert IDs to actual objects (you'll need to implement this)
				updates[dbTag] = field.Interface()
			} else {
				// Handle full object slices
				updates[dbTag] = field.Interface()
			}
			fields = append(fields, dbTag)
		}
	}

	return updates, fields
}

func (dto UpdateProjectDto) HasUpdates() bool {
	_, fields := dto.GetUpdates()
	return len(fields) > 0
}

func ValidateUpdateProjectDto(dto UpdateProjectDto) error {
	validate := validator.New()

	// Register custom validation for complex types if needed
	validate.RegisterValidation("valid_project_deployment", validateProjectDeployment)

	return validate.Struct(dto)
}

// Вспомогательные функции для создания указателей
func StringPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}

func UintSlicePtr(slice []uint) []uint {
	if slice == nil {
		return []uint{}
	}
	return slice
}
