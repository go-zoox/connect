package model

import (
	"time"

	"gorm.io/gorm"
)

// Group is an admin user group.
type Group struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:191;uniqueIndex;not null"`
	Description string `gorm:"size:512"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Group) TableName() string {
	return "admin_groups"
}
