package model

import (
	"time"

	"gorm.io/gorm"
)

// User is an admin console account.
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"size:191;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName overrides singular-table naming to avoid reserved SQL keywords.
func (User) TableName() string {
	return "admin_users"
}
