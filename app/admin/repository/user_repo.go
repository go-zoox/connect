package repository

import (
	"github.com/go-zoox/connect/app/admin/model"
	"gorm.io/gorm"
)

// UserRepo loads admin users from the database.
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo returns a UserRepo backed by db.
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetByUsername returns the active admin user with the given username, or gorm.ErrRecordNotFound.
func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
