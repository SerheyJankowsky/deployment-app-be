package dto

import (
	"deployer.com/modules/containers"
	"deployer.com/modules/domains"
	"deployer.com/modules/scripts"
	"deployer.com/modules/secrets"
	"deployer.com/modules/servers"

	// "deployer.com/modules/servers" // Commented out due to import issue
	"github.com/go-playground/validator/v10"
)

type CreateDeploymentDto struct {
	Name                  string `json:"name" validate:"required,min=1,max=255"`
	SetUpDomains          bool   `json:"setup_domains"`
	PoolContainers        bool   `json:"pool_containers"`
	RunContainers         bool   `json:"run_containers"`
	SetUpServers          bool   `json:"setup_servers"`
	SetSecretsToServer    bool   `json:"set_secrets_to_server"`
	SetSecretsToContainer bool   `json:"set_secrets_to_container"`
	RunScripts            bool   `json:"run_scripts"`

	// Fixed validation tags - removed 'min=1' from complex structs
	Domains    []domains.Domain       `json:"domains" validate:"omitempty,dive"`
	SubDomains []domains.SubDomain    `json:"sub_domains" validate:"omitempty,dive"`
	Containers []containers.Container `json:"containers" validate:"omitempty,dive"`
	Servers    []servers.Server       `json:"servers" validate:"omitempty,dive"` // Commented out
	Scripts    []scripts.Script       `json:"scripts" validate:"omitempty,dive"`
	Secrets    []secrets.Secret       `json:"secrets" validate:"omitempty,dive"`

	// Alternative: Use IDs instead of full objects for relationships (recommended)
	DomainIDs    []uint `json:"domain_ids" validate:"omitempty,dive,min=1"`
	SubDomainIDs []uint `json:"subdomain_ids" validate:"omitempty,dive,min=1"`
	ContainerIDs []uint `json:"container_ids" validate:"omitempty,dive,min=1"`
	ServerIDs    []uint `json:"server_ids" validate:"omitempty,dive,min=1"` // Commented out
	ScriptIDs    []uint `json:"script_ids" validate:"omitempty,dive,min=1"`
	SecretIDs    []uint `json:"secret_ids" validate:"omitempty,dive,min=1"`
}

func ValidateCreateDeploymentDto(dto CreateDeploymentDto) error {
	validate := validator.New()

	// Register custom validation for complex types if needed
	validate.RegisterValidation("valid_domain", validateDomain)
	validate.RegisterValidation("valid_container", validateContainer)
	validate.RegisterValidation("valid_script", validateScript)
	validate.RegisterValidation("valid_secret", validateSecret)

	return validate.Struct(dto)
}

// Helper method to convert IDs to actual objects (to be used in service layer)
func (dto *CreateDeploymentDto) UseIDsOnly() bool {
	// Return true if any ID fields are populated, indicating client prefers ID-based approach
	return len(dto.DomainIDs) > 0 || len(dto.SubDomainIDs) > 0 ||
		len(dto.ContainerIDs) > 0 || len(dto.ScriptIDs) > 0 ||
		len(dto.SecretIDs) > 0
}

// Helper method to check if using full objects
func (dto *CreateDeploymentDto) UseFullObjects() bool {
	// Return true if any object fields are populated
	return len(dto.Domains) > 0 || len(dto.SubDomains) > 0 ||
		len(dto.Containers) > 0 || len(dto.Scripts) > 0 ||
		len(dto.Secrets) > 0
}
