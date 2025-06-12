package dto

import (
	"deployer.com/modules/containers"
	"deployer.com/modules/domains"
	"deployer.com/modules/scripts"
	"deployer.com/modules/secrets"
	"deployer.com/modules/servers"
	"github.com/go-playground/validator/v10"
)

type CreateDeploymentDto struct {
	Name                  string                 `json:"name" validate:"required,min=1,max=255"`
	SetUpDomains          bool                   `json:"setup_domains"`
	PoolContainers        bool                   `json:"pool_containers"`
	RunContainers         bool                   `json:"run_containers"`
	SetUpServers          bool                   `json:"setup_servers"`
	SetSecretsToServer    bool                   `json:"set_secrets_to_server"`
	SetSecretsToContainer bool                   `json:"set_secrets_to_container"`
	RunScripts            bool                   `json:"run_scripts"`
	Domains               []domains.Domain       `json:"domains" validate:"omitempty,dive,min=1"`
	SubDomains            []domains.SubDomain    `json:"sub_domains" validate:"omitempty,dive,min=1"`
	Containers            []containers.Container `json:"containers" validate:"omitempty,dive,min=1"`
	Servers               []servers.Server       `json:"servers" validate:"omitempty,dive,min=1"`
	Scripts               []scripts.Script       `json:"scripts" validate:"omitempty,dive,min=1"`
	Secrets               []secrets.Secret       `json:"secrets" validate:"omitempty,dive,min=1"`
}

func ValidateCreateDeploymentDto(dto CreateDeploymentDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
