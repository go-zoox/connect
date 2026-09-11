package model

import (
	"time"

	"gorm.io/gorm"
)

// Permission is an admin RBAC permission key.
type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Key         string `gorm:"size:191;uniqueIndex;not null"`
	Description string `gorm:"size:512"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Permission) TableName() string {
	return "admin_permissions"
}
