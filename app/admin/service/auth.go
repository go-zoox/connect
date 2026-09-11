package service

import (
	"errors"
	"strings"

	"github.com/go-zoox/connect/app/admin/model"
	"github.com/go-zoox/connect/app/admin/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthLoginStatus is the outcome of a login attempt.
type AuthLoginStatus string

const (
	// AuthLoginOK means credentials matched an active admin user.
	AuthLoginOK AuthLoginStatus = "ok"
	// AuthLoginInvalidCredentials means the username or password did not match.
	AuthLoginInvalidCredentials AuthLoginStatus = "invalid_credentials"
)

// AuthLoginResult is returned by AuthService.Login.
type AuthLoginResult struct {
	Status AuthLoginStatus
	User   *model.User
}

// AuthService verifies admin credentials against the users table.
type AuthService struct {
	users *repository.UserRepo
}

// NewAuthService builds an AuthService.
func NewAuthService(users *repository.UserRepo) *AuthService {
	return &AuthService{users: users}
}

// Login checks username and password using the stored bcrypt hash.
func (s *AuthService) Login(username, password string) (*AuthLoginResult, error) {
	username = strings.TrimSpace(username)
	u, err := s.users.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &AuthLoginResult{Status: AuthLoginInvalidCredentials}, nil
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return &AuthLoginResult{Status: AuthLoginInvalidCredentials}, nil
	}

	out := *u
	out.PasswordHash = ""
	return &AuthLoginResult{Status: AuthLoginOK, User: &out}, nil
}
