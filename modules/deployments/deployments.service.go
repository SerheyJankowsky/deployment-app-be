package deployments

import (
	"fmt"
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

// Helper function to safely preload relations, ignoring errors if tables don't exist
func (s *DeploymentsService) safePreloadRelations(deployment *Deployment) {
	// Try to preload each relationship, ignore errors if junction tables don't exist

	// Try Domains
	if err := s.db.Model(deployment).Association("Domains").Find(&deployment.Domains); err != nil {
		deployment.Domains = []domains.Domain{} // Set empty slice if error
	}

	// Try SubDomains
	if err := s.db.Model(deployment).Association("SubDomains").Find(&deployment.SubDomains); err != nil {
		deployment.SubDomains = []domains.SubDomain{} // Set empty slice if error
	}

	// Try Containers
	if err := s.db.Model(deployment).Association("Containers").Find(&deployment.Containers); err != nil {
		deployment.Containers = []containers.Container{} // Set empty slice if error
	}

	// Try Servers
	if err := s.db.Model(deployment).Association("Servers").Find(&deployment.Servers); err != nil {
		deployment.Servers = []servers.Server{} // Set empty slice if error
	}

	// Try Scripts
	if err := s.db.Model(deployment).Association("Scripts").Find(&deployment.Scripts); err != nil {
		deployment.Scripts = []scripts.Script{} // Set empty slice if error
	}

	// Try Secrets
	if err := s.db.Model(deployment).Association("Secrets").Find(&deployment.Secrets); err != nil {
		deployment.Secrets = []secrets.Secret{} // Set empty slice if error
	}
}

// Helper function to check if a table exists
func (s *DeploymentsService) tableExists(tableName string) bool {
	var count int64
	err := s.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?", tableName).Scan(&count).Error
	return err == nil && count > 0
}

// Helper function to safely create associations with user validation
func (s *DeploymentsService) safeCreateAssociations(deployment *Deployment, dto dto.CreateDeploymentDto, userId uint) error {
	// Only create associations if the data exists and tables exist, with user validation

	if len(dto.Domains) > 0 && s.tableExists("deployment_domains") {
		// Extract IDs and validate ownership in one query
		domainIDs := make([]uint, len(dto.Domains))
		for i, domain := range dto.Domains {
			domainIDs[i] = domain.ID
		}
		var count int64
		s.db.Model(&domains.Domain{}).Where("id IN ? AND user_id = ?", domainIDs, userId).Count(&count)
		if count != int64(len(domainIDs)) {
			return fmt.Errorf("some domains do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("Domains").Replace(dto.Domains)
	}

	if len(dto.SubDomains) > 0 && s.tableExists("deployment_subdomains") {
		// Extract IDs and validate ownership in one query
		subDomainIDs := make([]uint, len(dto.SubDomains))
		for i, subdomain := range dto.SubDomains {
			subDomainIDs[i] = subdomain.ID
		}
		var count int64
		s.db.Model(&domains.SubDomain{}).Where("id IN ? AND user_id = ?", subDomainIDs, userId).Count(&count)
		if count != int64(len(subDomainIDs)) {
			return fmt.Errorf("some subdomains do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("SubDomains").Replace(dto.SubDomains)
	}

	if len(dto.Containers) > 0 && s.tableExists("deployment_containers") {
		// Extract IDs and validate ownership in one query
		containerIDs := make([]uint, len(dto.Containers))
		for i, container := range dto.Containers {
			containerIDs[i] = container.ID
		}
		var count int64
		s.db.Model(&containers.Container{}).Where("id IN ? AND user_id = ?", containerIDs, userId).Count(&count)
		if count != int64(len(containerIDs)) {
			return fmt.Errorf("some containers do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("Containers").Replace(dto.Containers)
	}

	if len(dto.Servers) > 0 && s.tableExists("deployment_servers") {
		// Extract IDs and validate ownership in one query
		serverIDs := make([]uint, len(dto.Servers))
		for i, server := range dto.Servers {
			serverIDs[i] = server.ID
		}
		var count int64
		s.db.Model(&servers.Server{}).Where("id IN ? AND user_id = ?", serverIDs, userId).Count(&count)
		if count != int64(len(serverIDs)) {
			return fmt.Errorf("some servers do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("Servers").Replace(dto.Servers)
	}

	if len(dto.Scripts) > 0 && s.tableExists("deployment_scripts") {
		// Extract IDs and validate ownership in one query
		scriptIDs := make([]uint, len(dto.Scripts))
		for i, script := range dto.Scripts {
			scriptIDs[i] = script.ID
		}
		var count int64
		s.db.Model(&scripts.Script{}).Where("id IN ? AND user_id = ?", scriptIDs, userId).Count(&count)
		if count != int64(len(scriptIDs)) {
			return fmt.Errorf("some scripts do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("Scripts").Replace(dto.Scripts)
	}

	if len(dto.Secrets) > 0 && s.tableExists("deployment_secrets") {
		// Extract IDs and validate ownership in one query
		secretIDs := make([]uint, len(dto.Secrets))
		for i, secret := range dto.Secrets {
			secretIDs[i] = secret.ID
		}
		var count int64
		s.db.Model(&secrets.Secret{}).Where("id IN ? AND user_id = ?", secretIDs, userId).Count(&count)
		if count != int64(len(secretIDs)) {
			return fmt.Errorf("some secrets do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("Secrets").Replace(dto.Secrets)
	}

	return nil
}

// Helper function to safely create associations from IDs with user validation
func (s *DeploymentsService) safeCreateAssociationsFromIDs(deployment *Deployment, dto dto.CreateDeploymentDto, userId uint) error {
	// Convert IDs to objects and create associations with user validation

	if len(dto.DomainIDs) > 0 && s.tableExists("deployment_domains") {
		var domains []domains.Domain
		if err := s.db.Where("id IN ? AND user_id = ?", dto.DomainIDs, userId).Find(&domains).Error; err != nil {
			return err
		}
		// Check if all requested domains were found (belong to user)
		if len(domains) != len(dto.DomainIDs) {
			return fmt.Errorf("some domains do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("Domains").Replace(domains)
	}

	if len(dto.SubDomainIDs) > 0 && s.tableExists("deployment_subdomains") {
		var subdomains []domains.SubDomain
		if err := s.db.Where("id IN ? AND user_id = ?", dto.SubDomainIDs, userId).Find(&subdomains).Error; err != nil {
			return err
		}
		// Check if all requested subdomains were found (belong to user)
		if len(subdomains) != len(dto.SubDomainIDs) {
			return fmt.Errorf("some subdomains do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("SubDomains").Replace(subdomains)
	}

	if len(dto.ContainerIDs) > 0 && s.tableExists("deployment_containers") {
		var containers []containers.Container
		if err := s.db.Where("id IN ? AND user_id = ?", dto.ContainerIDs, userId).Find(&containers).Error; err != nil {
			return err
		}
		// Check if all requested containers were found (belong to user)
		if len(containers) != len(dto.ContainerIDs) {
			return fmt.Errorf("some containers do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("Containers").Replace(containers)
	}

	if len(dto.ServerIDs) > 0 && s.tableExists("deployment_servers") {
		var servers []servers.Server
		if err := s.db.Where("id IN ? AND user_id = ?", dto.ServerIDs, userId).Find(&servers).Error; err != nil {
			return err
		}
		// Check if all requested servers were found (belong to user)
		if len(servers) != len(dto.ServerIDs) {
			return fmt.Errorf("some servers do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("Servers").Replace(servers)
	}

	if len(dto.ScriptIDs) > 0 && s.tableExists("deployment_scripts") {
		var scripts []scripts.Script
		if err := s.db.Where("id IN ? AND user_id = ?", dto.ScriptIDs, userId).Find(&scripts).Error; err != nil {
			return err
		}
		// Check if all requested scripts were found (belong to user)
		if len(scripts) != len(dto.ScriptIDs) {
			return fmt.Errorf("some scripts do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("Scripts").Replace(scripts)
	}

	if len(dto.SecretIDs) > 0 && s.tableExists("deployment_secrets") {
		var secrets []secrets.Secret
		if err := s.db.Where("id IN ? AND user_id = ?", dto.SecretIDs, userId).Find(&secrets).Error; err != nil {
			return err
		}
		// Check if all requested secrets were found (belong to user)
		if len(secrets) != len(dto.SecretIDs) {
			return fmt.Errorf("some secrets do not belong to this user or do not exist")
		}
		s.db.Model(deployment).Association("Secrets").Replace(secrets)
	}

	return nil
}

func (s *DeploymentsService) safeUpdateAssociations(deployment *Deployment, updates map[string]interface{}, userId uint) error {
	// Handle many-to-many associations separately (only if tables exist) with user validation

	if d, ok := updates["domains"]; ok {
		if domainList, ok := d.([]domains.Domain); ok && s.tableExists("deployment_domains") {
			// Extract IDs and validate ownership in one query
			domainIDs := make([]uint, len(domainList))
			for i, domain := range domainList {
				domainIDs[i] = domain.ID
			}
			var count int64
			s.db.Model(&domains.Domain{}).Where("id IN ? AND user_id = ?", domainIDs, userId).Count(&count)
			if count != int64(len(domainIDs)) {
				return fmt.Errorf("some domains do not belong to this user or do not exist")
			}
			s.db.Model(deployment).Association("Domains").Replace(domainList)
		}
		delete(updates, "domains")
	}

	if subDomains, ok := updates["sub_domains"]; ok {
		if subDomainList, ok := subDomains.([]domains.SubDomain); ok && s.tableExists("deployment_subdomains") {
			// Extract IDs and validate ownership in one query
			subDomainIDs := make([]uint, len(subDomainList))
			for i, subdomain := range subDomainList {
				subDomainIDs[i] = subdomain.ID
			}
			var count int64
			s.db.Model(&domains.SubDomain{}).Where("id IN ? AND user_id = ?", subDomainIDs, userId).Count(&count)
			if count != int64(len(subDomainIDs)) {
				return fmt.Errorf("some subdomains do not belong to this user or do not exist")
			}
			s.db.Model(deployment).Association("SubDomains").Replace(subDomainList)
		}
		delete(updates, "sub_domains")
	}

	if c, ok := updates["containers"]; ok {
		if containerList, ok := c.([]containers.Container); ok && s.tableExists("deployment_containers") {
			// Extract IDs and validate ownership in one query
			containerIDs := make([]uint, len(containerList))
			for i, container := range containerList {
				containerIDs[i] = container.ID
			}
			var count int64
			s.db.Model(&containers.Container{}).Where("id IN ? AND user_id = ?", containerIDs, userId).Count(&count)
			if count != int64(len(containerIDs)) {
				return fmt.Errorf("some containers do not belong to this user or do not exist")
			}
			s.db.Model(deployment).Association("Containers").Replace(containerList)
		}
		delete(updates, "containers")
	}

	if srv, ok := updates["servers"]; ok {
		if serverList, ok := srv.([]servers.Server); ok && s.tableExists("deployment_servers") {
			// Extract IDs and validate ownership in one query
			serverIDs := make([]uint, len(serverList))
			for i, server := range serverList {
				serverIDs[i] = server.ID
			}
			var count int64
			s.db.Model(&servers.Server{}).Where("id IN ? AND user_id = ?", serverIDs, userId).Count(&count)
			if count != int64(len(serverIDs)) {
				return fmt.Errorf("some servers do not belong to this user or do not exist")
			}
			s.db.Model(deployment).Association("Servers").Replace(serverList)
		}
		delete(updates, "servers")
	}

	if sc, ok := updates["scripts"]; ok {
		if scriptList, ok := sc.([]scripts.Script); ok && s.tableExists("deployment_scripts") {
			// Extract IDs and validate ownership in one query
			scriptIDs := make([]uint, len(scriptList))
			for i, script := range scriptList {
				scriptIDs[i] = script.ID
			}
			var count int64
			s.db.Model(&scripts.Script{}).Where("id IN ? AND user_id = ?", scriptIDs, userId).Count(&count)
			if count != int64(len(scriptIDs)) {
				return fmt.Errorf("some scripts do not belong to this user or do not exist")
			}
			s.db.Model(deployment).Association("Scripts").Replace(scriptList)
		}
		delete(updates, "scripts")
	}

	if se, ok := updates["secrets"]; ok {
		if secretList, ok := se.([]secrets.Secret); ok && s.tableExists("deployment_secrets") {
			// Extract IDs and validate ownership in one query
			secretIDs := make([]uint, len(secretList))
			for i, secret := range secretList {
				secretIDs[i] = secret.ID
			}
			var count int64
			s.db.Model(&secrets.Secret{}).Where("id IN ? AND user_id = ?", secretIDs, userId).Count(&count)
			if count != int64(len(secretIDs)) {
				return fmt.Errorf("some secrets do not belong to this user or do not exist")
			}
			s.db.Model(deployment).Association("Secrets").Replace(secretList)
		}
		delete(updates, "secrets")
	}

	return nil
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

	// First, get deployments without preloading to avoid junction table errors
	if err := s.db.Where("user_id = ?", userId).
		Order("created_at DESC").
		Find(&deployments).Error; err != nil {
		return nil, err
	}

	result := make([]DeploymentResponse, len(deployments))
	for i, deployment := range deployments {
		// Try to preload each relationship individually and handle errors gracefully
		s.safePreloadRelations(&deployment)
		result[i] = s.convertToResponse(deployment)
	}

	return result, nil
}

func (s *DeploymentsService) GetDeployment(id, userId uint, iv string) (DeploymentResponse, error) {
	var deployment Deployment

	// First, get deployment without preloading to avoid junction table errors
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).
		First(&deployment).Error; err != nil {
		return DeploymentResponse{}, err
	}

	// Try to preload each relationship individually and handle errors gracefully
	s.safePreloadRelations(&deployment)

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

	// Handle relationships based on which approach is used
	if dto.UseIDsOnly() {
		// Convert IDs to objects and create associations with user validation
		if err := s.safeCreateAssociationsFromIDs(&deployment, dto, userId); err != nil {
			// Delete the created deployment if association fails
			s.db.Delete(&deployment)
			return DeploymentResponse{}, err
		}
	} else if dto.UseFullObjects() {
		// Use full objects directly with user validation
		if err := s.safeCreateAssociations(&deployment, dto, userId); err != nil {
			// Delete the created deployment if association fails
			s.db.Delete(&deployment)
			return DeploymentResponse{}, err
		}
	}

	// Try to preload relations safely
	s.safePreloadRelations(&deployment)

	return s.convertToResponse(deployment), nil
}

func (s *DeploymentsService) UpdateDeployment(id, userId uint, updates map[string]interface{}) (DeploymentResponse, error) {
	var deployment Deployment

	// Find the deployment
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&deployment).Error; err != nil {
		return DeploymentResponse{}, err
	}

	// Handle associations safely with user validation
	if err := s.safeUpdateAssociations(&deployment, updates, userId); err != nil {
		return DeploymentResponse{}, err
	}

	// Update other fields
	if len(updates) > 0 {
		libs.SetStructFieldsFromMap(&deployment, updates)
		if err := s.db.Save(&deployment).Error; err != nil {
			return DeploymentResponse{}, err
		}
	}

	// Try to preload relations safely
	s.safePreloadRelations(&deployment)

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
