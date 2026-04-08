package model

import (
	"time"

	"gorm.io/gorm"
)

// Role is an admin RBAC role.
type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:191;uniqueIndex;not null"`
	Description string `gorm:"size:512"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Role) TableName() string {
	return "admin_roles"
}
