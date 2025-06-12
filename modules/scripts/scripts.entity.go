package scripts

import (
	"time"

	"deployer.com/modules/users"
	"gorm.io/gorm"
)

type Script struct {
	gorm.Model
	// ID        uint           `gorm:"primaryKey" json:"id"`
	Name   string     `gorm:"not null;index" json:"name"`
	Script string     `gorm:"not null" json:"script"`
	User   users.User `gorm:"foreignKey:UserID" json:"-"`
	UserID uint       `gorm:"not null" json:"user_id"`
	// CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	LastRunAt *time.Time `gorm:"index;default:null" json:"last_run_at"`
}
