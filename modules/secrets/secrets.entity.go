package secrets

import (
	"time"

	"deployer.com/modules/users"
	"gorm.io/gorm"
)

type Secret struct {
	*gorm.Model
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Content   string         `gorm:"not null" json:"content"`
	User      users.User     `gorm:"foreignKey:UserID" json:"user"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
