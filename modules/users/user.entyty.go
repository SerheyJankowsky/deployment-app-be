package users

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Username  string `gorm:"unique;not null" json:"username"`
	Email     string `gorm:"unique;not null" json:"email"`
	Phone     string `gorm:"not null" json:"phone"`
	Country   string `gorm:"not null" json:"country"`
	ApiKey    string `gorm:"default:null" json:"-"`
	// City         string    `gorm:"not null" json:"city"`
	// Address      string    `gorm:"not null" json:"address"`
	// ZipCode      string    `gorm:"not null" json:"zip_code"`
	// Company      string    `gorm:"not null" json:"company"`
	// JobTitle     string    `gorm:"not null" json:"job_title"`
	// Bio          string    `gorm:"not null" json:"bio"`
	PasswordHash string `gorm:"not null" json:"-"`
	IV           string `gorm:"not null" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
