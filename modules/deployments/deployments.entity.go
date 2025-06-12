package deployments

import (
	"time"

	"deployer.com/modules/containers"
	"deployer.com/modules/domains"
	"deployer.com/modules/scripts"
	"deployer.com/modules/secrets"
	"deployer.com/modules/servers"
	"deployer.com/modules/users"
	"gorm.io/gorm"
)

type DeploymentStatus string

const (
	DeploymentStatusPending  DeploymentStatus = "pending"
	DeploymentStatusRunning  DeploymentStatus = "running"
	DeploymentStatusSuccess  DeploymentStatus = "success"
	DeploymentStatusFailed   DeploymentStatus = "failed"
	DeploymentStatusCanceled DeploymentStatus = "canceled"
	DeploymentStatusSkipped  DeploymentStatus = "skipped"
)

type Deployment struct {
	gorm.Model
	Name   string     `gorm:"not null;index" json:"name"`
	User   users.User `gorm:"foreignKey:UserID" json:"-"`
	UserID uint       `gorm:"not null" json:"user_id"`

	// Many-to-many relationships with junction tables
	Domains    []domains.Domain       `gorm:"many2many:deployment_domains;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"domains"`
	SubDomains []domains.SubDomain    `gorm:"many2many:deployment_subdomains;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"sub_domains"`
	Containers []containers.Container `gorm:"many2many:deployment_containers;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"containers"`
	Servers    []servers.Server       `gorm:"many2many:deployment_servers;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"servers"`
	Scripts    []scripts.Script       `gorm:"many2many:deployment_scripts;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"scripts"`
	Secrets    []secrets.Secret       `gorm:"many2many:deployment_secrets;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"secrets"`

	Status                DeploymentStatus `gorm:"not null" json:"status"`
	LastRunAt             *time.Time       `gorm:"index;default:null" json:"last_run_at"`
	SetUpDomains          bool             `gorm:"not null;default:false" json:"setup_domains"`
	PoolContainers        bool             `gorm:"not null;default:false" json:"pool_containers"`
	RunContainers         bool             `gorm:"not null;default:false" json:"run_containers"`
	SetUpServers          bool             `gorm:"not null;default:false" json:"setup_servers"`
	SetSecretsToServer    bool             `gorm:"not null;default:false" json:"set_secrets_to_server"`
	SetSecretsToContainer bool             `gorm:"not null;default:false" json:"set_secrets_to_container"`
	RunScripts            bool             `gorm:"not null;default:false" json:"run_script"`
}
