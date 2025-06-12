package dto

import (
	"reflect"

	"deployer.com/modules/containers"
	"deployer.com/modules/domains"
	"deployer.com/modules/scripts"
	"deployer.com/modules/secrets"
	"deployer.com/modules/servers"

	// "deployer.com/modules/servers" // Commented out due to import issue
	"github.com/go-playground/validator/v10"
)

type UpdateDeploymentDto struct {
	Name                  *string `json:"name" validate:"omitempty,min=1,max=255" db:"name"`
	SetUpDomains          *bool   `json:"setup_domains" db:"setup_domains"`
	PoolContainers        *bool   `json:"pool_containers" db:"pool_containers"`
	RunContainers         *bool   `json:"run_containers" db:"run_containers"`
	SetUpServers          *bool   `json:"setup_servers" db:"setup_servers"`
	SetSecretsToServer    *bool   `json:"set_secrets_to_server" db:"set_secrets_to_server"`
	SetSecretsToContainer *bool   `json:"set_secrets_to_container" db:"set_secrets_to_container"`
	RunScripts            *bool   `json:"run_scripts" db:"run_scripts"`

	// Fixed validation tags - removed 'min=1' from complex structs
	Domains    []domains.Domain       `json:"domains" validate:"omitempty,dive" db:"domains"`
	SubDomains []domains.SubDomain    `json:"sub_domains" validate:"omitempty,dive" db:"sub_domains"`
	Containers []containers.Container `json:"containers" validate:"omitempty,dive" db:"containers"`
	Servers    []servers.Server       `json:"servers" validate:"omitempty,dive" db:"servers"` // Commented out
	Scripts    []scripts.Script       `json:"scripts" validate:"omitempty,dive" db:"scripts"`
	Secrets    []secrets.Secret       `json:"secrets" validate:"omitempty,dive" db:"secrets"`

	// Alternative: Use IDs instead of full objects for relationships
	DomainIDs    []uint `json:"domain_ids" validate:"omitempty,dive,min=1" db:"domain_ids"`
	SubDomainIDs []uint `json:"subdomain_ids" validate:"omitempty,dive,min=1" db:"subdomain_ids"`
	ContainerIDs []uint `json:"container_ids" validate:"omitempty,dive,min=1" db:"container_ids"`
	ServerIDs    []uint `json:"server_ids" validate:"omitempty,dive,min=1" db:"server_ids"` // Commented out
	ScriptIDs    []uint `json:"script_ids" validate:"omitempty,dive,min=1" db:"script_ids"`
	SecretIDs    []uint `json:"secret_ids" validate:"omitempty,dive,min=1" db:"secret_ids"`
}

func (dto *UpdateDeploymentDto) GetUpdates() (map[string]interface{}, []string) {
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
			if fieldName == "DomainIDs" || fieldName == "SubDomainIDs" ||
				fieldName == "ContainerIDs" || fieldName == "ScriptIDs" ||
				fieldName == "SecretIDs" {
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

func (dto UpdateDeploymentDto) HasUpdates() bool {
	_, fields := dto.GetUpdates()
	return len(fields) > 0
}

func ValidateUpdateDeploymentDto(dto UpdateDeploymentDto) error {
	validate := validator.New()

	// Register custom validation for complex types if needed
	validate.RegisterValidation("valid_domain", validateDomain)
	validate.RegisterValidation("valid_container", validateContainer)
	validate.RegisterValidation("valid_script", validateScript)
	validate.RegisterValidation("valid_secret", validateSecret)

	return validate.Struct(dto)
}

// Custom validation functions
func validateDomain(fl validator.FieldLevel) bool {
	domain, ok := fl.Field().Interface().(domains.Domain)
	if !ok {
		return false
	}
	// Add your domain validation logic here
	return domain.ID > 0 // Example: check if ID exists
}

func validateContainer(fl validator.FieldLevel) bool {
	container, ok := fl.Field().Interface().(containers.Container)
	if !ok {
		return false
	}
	// Add your container validation logic here
	return container.ID > 0 // Example: check if ID exists
}

func validateScript(fl validator.FieldLevel) bool {
	script, ok := fl.Field().Interface().(scripts.Script)
	if !ok {
		return false
	}
	// Add your script validation logic here
	return script.ID > 0 // Example: check if ID exists
}

func validateSecret(fl validator.FieldLevel) bool {
	secret, ok := fl.Field().Interface().(secrets.Secret)
	if !ok {
		return false
	}
	// Add your secret validation logic here
	return secret.ID > 0 // Example: check if ID exists
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
