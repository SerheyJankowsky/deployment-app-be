package domains

import (
	"deployer.com/modules/users"
	"gorm.io/gorm"
)

type Domain struct {
	*gorm.Model
	// ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string      `gorm:"not null;index" json:"name"`
	SSLCert    string      `gorm:"not null" json:"ssl_cert"`
	SSLKey     string      `gorm:"not null" json:"ssl_key"`
	SubDomains []SubDomain `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	User       users.User  `gorm:"foreignKey:UserID" json:"-"`
	UserID     uint        `gorm:"not null" json:"user_id"`
	// CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type SubDomain struct {
	*gorm.Model
	// ID        uint           `gorm:"primaryKey" json:"id"`
	Name     string     `gorm:"not null;uniqueIndex:idx_domain_subdomain" json:"name"`
	Domain   Domain     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	DomainID uint       `gorm:"not null;uniqueIndex:idx_domain_subdomain;index" json:"domain_id"`
	User     users.User `gorm:"foreignKey:UserID" json:"-"`
	UserID   uint       `gorm:"not null" json:"user_id"`
	// CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
