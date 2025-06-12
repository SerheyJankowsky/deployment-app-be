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

type DeploymentResponse struct {
	ID                    uint                   `json:"id"`
	Name                  string                 `json:"name"`
	Domains               []domains.Domain       `json:"domains"`
	SubDomains            []domains.SubDomain    `json:"sub_domains"`
	Containers            []containers.Container `json:"containers"`
	Servers               []servers.Server       `json:"servers"`
	Scripts               []scripts.Script       `json:"scripts"`
	Secrets               []secrets.Secret       `json:"secrets"`
	Status                DeploymentStatus       `json:"status"`
	LastRunAt             *time.Time             `json:"last_run_at"`
	SetUpDomains          bool                   `json:"setup_domains"`
	PoolContainers        bool                   `json:"pool_containers"`
	RunContainers         bool                   `json:"run_containers"`
	SetUpServers          bool                   `json:"setup_servers"`
	SetSecretsToServer    bool                   `json:"set_secrets_to_server"`
	SetSecretsToContainer bool                   `json:"set_secrets_to_container"`
	RunScripts            bool                   `json:"run_scripts"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
}

func NewDeploymentsService(db *gorm.DB) *DeploymentsService {
	return &DeploymentsService{db: db}
}

func (s *DeploymentsService) GetDeployments(userId uint, iv string) ([]DeploymentResponse, error) {
	var deployments []Deployment
	if err := s.db.Where("user_id = ?", userId).Order("created_at DESC").Find(&deployments).Error; err != nil {
		return nil, err
	}
	result := make([]DeploymentResponse, len(deployments))
	for i, deployment := range deployments {
		result[i] = DeploymentResponse{
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
			Domains:               deployment.Domains,
			SubDomains:            deployment.SubDomains,
			Containers:            deployment.Containers,
			Servers:               deployment.Servers,
			Scripts:               deployment.Scripts,
			Secrets:               deployment.Secrets,
			Status:                deployment.Status,
		}
	}
	return result, nil
}

func (s *DeploymentsService) GetDeployment(id, userId uint, iv string) (DeploymentResponse, error) {
	var deployment Deployment
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&deployment).Error; err != nil {
		return DeploymentResponse{}, err
	}
	return DeploymentResponse{
		ID:                    deployment.ID,
		Name:                  deployment.Name,
		CreatedAt:             deployment.CreatedAt,
		UpdatedAt:             deployment.UpdatedAt,
		LastRunAt:             deployment.LastRunAt,
		SetUpDomains:          deployment.SetUpDomains,
		PoolContainers:        deployment.PoolContainers,
		RunContainers:         deployment.RunContainers,
		SetUpServers:          deployment.SetUpServers,
		SetSecretsToServer:    deployment.SetSecretsToServer,
		SetSecretsToContainer: deployment.SetSecretsToContainer,
		RunScripts:            deployment.RunScripts,
		Domains:               deployment.Domains,
		SubDomains:            deployment.SubDomains,
		Containers:            deployment.Containers,
		Servers:               deployment.Servers,
		Scripts:               deployment.Scripts,
		Secrets:               deployment.Secrets,
		Status:                deployment.Status,
	}, nil
}

func (s *DeploymentsService) CreateDeployment(userId uint, dto dto.CreateDeploymentDto) (DeploymentResponse, error) {
	deployment := Deployment{
		Name:                  dto.Name,
		UserID:                userId,
		SetUpDomains:          dto.SetUpDomains,
		PoolContainers:        dto.PoolContainers,
		RunContainers:         dto.RunContainers,
		SetUpServers:          dto.SetUpServers,
		SetSecretsToServer:    dto.SetSecretsToServer,
		SetSecretsToContainer: dto.SetSecretsToContainer,
		RunScripts:            dto.RunScripts,
		Domains:               dto.Domains,
		SubDomains:            dto.SubDomains,
		Containers:            dto.Containers,
		Servers:               dto.Servers,
		Scripts:               dto.Scripts,
		Secrets:               dto.Secrets,
	}
	if err := s.db.Create(&deployment).Error; err != nil {
		return DeploymentResponse{}, err
	}
	return DeploymentResponse{
		ID:                    deployment.ID,
		Name:                  deployment.Name,
		CreatedAt:             deployment.CreatedAt,
		UpdatedAt:             deployment.UpdatedAt,
		LastRunAt:             deployment.LastRunAt,
		SetUpDomains:          deployment.SetUpDomains,
		PoolContainers:        deployment.PoolContainers,
		RunContainers:         deployment.RunContainers,
		SetUpServers:          deployment.SetUpServers,
		SetSecretsToServer:    deployment.SetSecretsToServer,
		SetSecretsToContainer: deployment.SetSecretsToContainer,
		RunScripts:            deployment.RunScripts,
		Domains:               deployment.Domains,
		SubDomains:            deployment.SubDomains,
		Containers:            deployment.Containers,
		Servers:               deployment.Servers,
		Scripts:               deployment.Scripts,
		Secrets:               deployment.Secrets,
		Status:                deployment.Status,
	}, nil
}

func (s *DeploymentsService) UpdateDeployment(id, userId uint, updates map[string]interface{}) (DeploymentResponse, error) {
	var deployment Deployment
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&deployment).Error; err != nil {
		return DeploymentResponse{}, err
	}
	libs.SetStructFieldsFromMap(&deployment, updates)
	if err := s.db.Save(&deployment).Error; err != nil {
		return DeploymentResponse{}, err
	}
	return DeploymentResponse{
		ID:                    deployment.ID,
		Name:                  deployment.Name,
		CreatedAt:             deployment.CreatedAt,
		UpdatedAt:             deployment.UpdatedAt,
		LastRunAt:             deployment.LastRunAt,
		SetUpDomains:          deployment.SetUpDomains,
		PoolContainers:        deployment.PoolContainers,
		RunContainers:         deployment.RunContainers,
		SetUpServers:          deployment.SetUpServers,
		SetSecretsToServer:    deployment.SetSecretsToServer,
		SetSecretsToContainer: deployment.SetSecretsToContainer,
		RunScripts:            deployment.RunScripts,
		Domains:               deployment.Domains,
		SubDomains:            deployment.SubDomains,
		Containers:            deployment.Containers,
		Servers:               deployment.Servers,
		Scripts:               deployment.Scripts,
		Secrets:               deployment.Secrets,
		Status:                deployment.Status,
	}, nil
}

func (s *DeploymentsService) DeleteDeployment(id, userId uint) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Delete(&Deployment{}).Error; err != nil {
		return err
	}
	return nil
}
