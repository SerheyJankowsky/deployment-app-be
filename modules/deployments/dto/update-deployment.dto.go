package dto

import (
	"reflect"

	"deployer.com/modules/containers"
	"deployer.com/modules/domains"
	"deployer.com/modules/scripts"
	"deployer.com/modules/secrets"
	"deployer.com/modules/servers"
	"github.com/go-playground/validator/v10"
)

type UpdateDeploymentDto struct {
	Name                  *string                `json:"name" validate:"omitempty,min=1,max=255" db:"name"`
	SetUpDomains          *bool                  `json:"setup_domains" db:"setup_domains"`
	PoolContainers        *bool                  `json:"pool_containers" db:"pool_containers"`
	RunContainers         *bool                  `json:"run_containers" db:"run_containers"`
	SetUpServers          *bool                  `json:"setup_servers" db:"setup_servers"`
	SetSecretsToServer    *bool                  `json:"set_secrets_to_server" db:"set_secrets_to_server"`
	SetSecretsToContainer *bool                  `json:"set_secrets_to_container" db:"set_secrets_to_container"`
	RunScripts            *bool                  `json:"run_scripts" db:"run_scripts"`
	Domains               []domains.Domain       `json:"domains" validate:"omitempty,dive,min=1" db:"domains"`
	SubDomains            []domains.SubDomain    `json:"sub_domains" validate:"omitempty,dive,min=1" db:"sub_domains"`
	Containers            []containers.Container `json:"containers" validate:"omitempty,dive,min=1" db:"containers"`
	Servers               []servers.Server       `json:"servers" validate:"omitempty,dive,min=1" db:"servers"`
	Scripts               []scripts.Script       `json:"scripts" validate:"omitempty,dive,min=1" db:"scripts"`
	Secrets               []secrets.Secret       `json:"secrets" validate:"omitempty,dive,min=1" db:"secrets"`
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
			updates[dbTag] = field.Interface()
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
