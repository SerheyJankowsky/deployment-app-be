package servers

import (
	"deployer.com/modules/users"
	"gorm.io/gorm"
)

type Server struct {
	*gorm.Model
	// ID        uint           `gorm:"primaryKey" json:"id"`
	Name     string     `gorm:"not null" json:"name"`
	Host     string     `gorm:"not null" json:"host"`
	Port     int        `gorm:"not null;default:22" json:"port"`
	Password string     `gorm:"not null" json:"-"`
	SSHKey   *string    `json:"ssh_key"`
	Username string     `gorm:"not null" json:"username"`
	User     users.User `gorm:"foreignKey:UserID" json:"user"`
	UserID   uint       `gorm:"not null" json:"user_id"`
	// CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
