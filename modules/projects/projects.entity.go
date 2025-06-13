package projects

import (
	"deployer.com/modules/deployments"
	"deployer.com/modules/users"
	"gorm.io/gorm"
)

type Project struct {
	*gorm.Model
	// ID        uint           `gorm:"primaryKey" json:"id"`
	Name               string               `gorm:"not null" json:"name"`
	ProjectDeployments []ProjectDeployments `gorm:"foreignKey:ProjectID" json:"-"`
	User               users.User           `gorm:"foreignKey:UserID" json:"user"`
	UserID             uint                 `gorm:"not null" json:"user_id"`
	// CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProjectDeployments struct {
	gorm.Model
	Deployment   deployments.Deployment `gorm:"foreignKey:DeploymentID" json:"deployment"`
	DeploymentID uint                   `gorm:"not null" json:"deployment_id"`
	Project      Project                `gorm:"foreignKey:ProjectID" json:"project"`
	ProjectID    uint                   `gorm:"not null" json:"project_id"`
	Order        int                    `gorm:"not null" json:"order"`
	Status       string                 `gorm:"not null" json:"status"`
	Logs         string                 `json:"logs"`
}
