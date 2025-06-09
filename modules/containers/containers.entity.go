package containers

import (
	"time"

	"deployer.com/modules/users"
	"gorm.io/gorm"
)

type Container struct {
	*gorm.Model
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Registry  string         `gorm:"not null" json:"registry"`
	Image     string         `gorm:"not null" json:"image"`
	Tag       string         `gorm:"not null" json:"tag"`
	Username  string         `gorm:"not null" json:"username"`
	Password  string         `json:"password"`
	SecretKey string         `json:"secret_key"`
	Params    string         `json:"params"`
	User      users.User     `gorm:"foreignKey:UserID" json:"user"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
