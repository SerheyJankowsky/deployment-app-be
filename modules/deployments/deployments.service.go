package deployments

import (
	"time"

	"deployer.com/libs"
	"deployer.com/modules/containers"
	"deployer.com/modules/deployments/dto"
	"deployer.com/modules/domains"
	"deployer.com/modules/scripts"
	"deployer.com/modules/secrets"
	"deployer.com/modules/servers"
	"gorm.io/gorm"
)

type DeploymentsService struct {
	db *gorm.DB
}

// Simplified response structures to get only ID and Name from relations
type DomainSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SubDomainSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ContainerSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ServerSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ScriptSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SecretSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type DeploymentResponse struct {
	ID                    uint               `json:"id"`
	Name                  string             `json:"name"`
	Domains               []DomainSummary    `json:"domains"`
	SubDomains            []SubDomainSummary `json:"sub_domains"`
	Containers            []ContainerSummary `json:"containers"`
	Servers               []ServerSummary    `json:"servers"`
	Scripts               []ScriptSummary    `json:"scripts"`
	Secrets               []SecretSummary    `json:"secrets"`
	Status                DeploymentStatus   `json:"status"`
	LastRunAt             *time.Time         `json:"last_run_at"`
	SetUpDomains          bool               `json:"setup_domains"`
	PoolContainers        bool               `json:"pool_containers"`
	RunContainers         bool               `json:"run_containers"`
	SetUpServers          bool               `json:"setup_servers"`
	SetSecretsToServer    bool               `json:"set_secrets_to_server"`
	SetSecretsToContainer bool               `json:"set_secrets_to_container"`
	RunScripts            bool               `json:"run_scripts"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
}

func NewDeploymentsService(db *gorm.DB) *DeploymentsService {
	return &DeploymentsService{db: db}
}

// Helper function to convert domain slice to summary slice
func convertDomainsToSummary(domains []domains.Domain) []DomainSummary {
	result := make([]DomainSummary, len(domains))
	for i, domain := range domains {
		result[i] = DomainSummary{
			ID:   domain.ID,
			Name: domain.Name,
		}
	}
	return result
}

// Helper function to convert subdomain slice to summary slice
func convertSubDomainsToSummary(subdomains []domains.SubDomain) []SubDomainSummary {
	result := make([]SubDomainSummary, len(subdomains))
	for i, subdomain := range subdomains {
		result[i] = SubDomainSummary{
			ID:   subdomain.ID,
			Name: subdomain.Name,
		}
	}
	return result
}

// Helper function to convert container slice to summary slice
func convertContainersToSummary(containers []containers.Container) []ContainerSummary {
	result := make([]ContainerSummary, len(containers))
	for i, container := range containers {
		result[i] = ContainerSummary{
			ID:   container.ID,
			Name: container.Name,
		}
	}
	return result
}

// Helper function to convert server slice to summary slice
func convertServersToSummary(servers []servers.Server) []ServerSummary {
	result := make([]ServerSummary, len(servers))
	for i, server := range servers {
		result[i] = ServerSummary{
			ID:   server.ID,
			Name: server.Name,
		}
	}
	return result
}

// Helper function to convert script slice to summary slice
func convertScriptsToSummary(scripts []scripts.Script) []ScriptSummary {
	result := make([]ScriptSummary, len(scripts))
	for i, script := range scripts {
		result[i] = ScriptSummary{
			ID:   script.ID,
			Name: script.Name,
		}
	}
	return result
}

// Helper function to convert secret slice to summary slice
func convertSecretsToSummary(secrets []secrets.Secret) []SecretSummary {
	result := make([]SecretSummary, len(secrets))
	for i, secret := range secrets {
		result[i] = SecretSummary{
			ID:   secret.ID,
			Name: secret.Name,
		}
	}
	return result
}

// Helper function to convert deployment to response
func (s *DeploymentsService) convertToResponse(deployment Deployment) DeploymentResponse {
	return DeploymentResponse{
		ID:                    deployment.ID,
		Name:                  deployment.Name,
		LastRunAt:             deployment.LastRunAt,
		SetUpDomains:          deployment.SetUpDomains,
		PoolContainers:        deployment.PoolContainers,
		RunContainers:         deployment.RunContainers,
		SetUpServers:          deployment.SetUpServers,
		SetSecretsToServer:    deployment.SetSecretsToServer,
		SetSecretsToContainer: deployment.SetSecretsToContainer,
		RunScripts:            deployment.RunScripts,
		CreatedAt:             deployment.CreatedAt,
		UpdatedAt:             deployment.UpdatedAt,
		Status:                deployment.Status,
		Domains:               convertDomainsToSummary(deployment.Domains),
		SubDomains:            convertSubDomainsToSummary(deployment.SubDomains),
		Containers:            convertContainersToSummary(deployment.Containers),
		Servers:               convertServersToSummary(deployment.Servers),
		Scripts:               convertScriptsToSummary(deployment.Scripts),
		Secrets:               convertSecretsToSummary(deployment.Secrets),
	}
}

func (s *DeploymentsService) GetDeployments(userId uint, iv string) ([]DeploymentResponse, error) {
	var deployments []Deployment

	// Preload all many-to-many relationships
	if err := s.db.Where("user_id = ?", userId).
		Preload("Domains").
		Preload("SubDomains").
		Preload("Containers").
		Preload("Servers").
		Preload("Scripts").
		Preload("Secrets").
		Order("created_at DESC").
		Find(&deployments).Error; err != nil {
		return nil, err
	}

	result := make([]DeploymentResponse, len(deployments))
	for i, deployment := range deployments {
		result[i] = s.convertToResponse(deployment)
	}

	return result, nil
}

func (s *DeploymentsService) GetDeployment(id, userId uint, iv string) (DeploymentResponse, error) {
	var deployment Deployment

	// Preload all many-to-many relationships
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).
		Preload("Domains").
		Preload("SubDomains").
		Preload("Containers").
		Preload("Servers").
		Preload("Scripts").
		Preload("Secrets").
		First(&deployment).Error; err != nil {
		return DeploymentResponse{}, err
	}

	return s.convertToResponse(deployment), nil
}

func (s *DeploymentsService) CreateDeployment(userId uint, dto dto.CreateDeploymentDto) (DeploymentResponse, error) {
	deployment := Deployment{
		Name:                  dto.Name,
		UserID:                userId,
		Status:                DeploymentStatusPending, // Set default status
		SetUpDomains:          dto.SetUpDomains,
		PoolContainers:        dto.PoolContainers,
		RunContainers:         dto.RunContainers,
		SetUpServers:          dto.SetUpServers,
		SetSecretsToServer:    dto.SetSecretsToServer,
		SetSecretsToContainer: dto.SetSecretsToContainer,
		RunScripts:            dto.RunScripts,
	}

	// Create the deployment first
	if err := s.db.Create(&deployment).Error; err != nil {
		return DeploymentResponse{}, err
	}

	// Handle many-to-many associations
	if len(dto.Domains) > 0 {
		if err := s.db.Model(&deployment).Association("Domains").Replace(dto.Domains); err != nil {
			return DeploymentResponse{}, err
		}
	}

	if len(dto.SubDomains) > 0 {
		if err := s.db.Model(&deployment).Association("SubDomains").Replace(dto.SubDomains); err != nil {
			return DeploymentResponse{}, err
		}
	}

	if len(dto.Containers) > 0 {
		if err := s.db.Model(&deployment).Association("Containers").Replace(dto.Containers); err != nil {
			return DeploymentResponse{}, err
		}
	}

	if len(dto.Servers) > 0 {
		if err := s.db.Model(&deployment).Association("Servers").Replace(dto.Servers); err != nil {
			return DeploymentResponse{}, err
		}
	}

	if len(dto.Scripts) > 0 {
		if err := s.db.Model(&deployment).Association("Scripts").Replace(dto.Scripts); err != nil {
			return DeploymentResponse{}, err
		}
	}

	if len(dto.Secrets) > 0 {
		if err := s.db.Model(&deployment).Association("Secrets").Replace(dto.Secrets); err != nil {
			return DeploymentResponse{}, err
		}
	}

	// Reload with associations
	if err := s.db.Preload("Domains").
		Preload("SubDomains").
		Preload("Containers").
		Preload("Servers").
		Preload("Scripts").
		Preload("Secrets").
		First(&deployment, deployment.ID).Error; err != nil {
		return DeploymentResponse{}, err
	}

	return s.convertToResponse(deployment), nil
}

func (s *DeploymentsService) UpdateDeployment(id, userId uint, updates map[string]interface{}) (DeploymentResponse, error) {
	var deployment Deployment

	// Find the deployment
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&deployment).Error; err != nil {
		return DeploymentResponse{}, err
	}

	// Handle many-to-many associations separately
	if d, ok := updates["domains"]; ok {
		if domainList, ok := d.([]domains.Domain); ok {
			if err := s.db.Model(&deployment).Association("Domains").Replace(domainList); err != nil {
				return DeploymentResponse{}, err
			}
		}
		delete(updates, "domains")
	}

	if subDomains, ok := updates["sub_domains"]; ok {
		if subDomainList, ok := subDomains.([]domains.SubDomain); ok {
			if err := s.db.Model(&deployment).Association("SubDomains").Replace(subDomainList); err != nil {
				return DeploymentResponse{}, err
			}
		}
		delete(updates, "sub_domains")
	}

	if c, ok := updates["containers"]; ok {
		if containerList, ok := c.([]containers.Container); ok {
			if err := s.db.Model(&deployment).Association("Containers").Replace(containerList); err != nil {
				return DeploymentResponse{}, err
			}
		}
		delete(updates, "containers")
	}

	if sr, ok := updates["servers"]; ok {
		if serverList, ok := sr.([]servers.Server); ok {
			if err := s.db.Model(&deployment).Association("Servers").Replace(serverList); err != nil {
				return DeploymentResponse{}, err
			}
		}
		delete(updates, "servers")
	}

	if sc, ok := updates["scripts"]; ok {
		if scriptList, ok := sc.([]scripts.Script); ok {
			if err := s.db.Model(&deployment).Association("Scripts").Replace(scriptList); err != nil {
				return DeploymentResponse{}, err
			}
		}
		delete(updates, "scripts")
	}

	if se, ok := updates["secrets"]; ok {
		if secretList, ok := se.([]secrets.Secret); ok {
			if err := s.db.Model(&deployment).Association("Secrets").Replace(secretList); err != nil {
				return DeploymentResponse{}, err
			}
		}
		delete(updates, "secrets")
	}

	// Update other fields
	if len(updates) > 0 {
		libs.SetStructFieldsFromMap(&deployment, updates)
		if err := s.db.Save(&deployment).Error; err != nil {
			return DeploymentResponse{}, err
		}
	}

	// Reload with associations
	if err := s.db.Preload("Domains").
		Preload("SubDomains").
		Preload("Containers").
		Preload("Servers").
		Preload("Scripts").
		Preload("Secrets").
		First(&deployment, deployment.ID).Error; err != nil {
		return DeploymentResponse{}, err
	}

	return s.convertToResponse(deployment), nil
}

func (s *DeploymentsService) DeleteDeployment(id, userId uint) error {
	// GORM will automatically handle the many-to-many relationship cleanup
	// when using CASCADE constraints
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Delete(&Deployment{}).Error; err != nil {
		return err
	}
	return nil
}
