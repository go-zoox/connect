package config

import (
	"fmt"
	"strings"
)

// Admin configures the built-in admin UI and its backing store.
type Admin struct {
	Enabled  bool                `config:"enabled"`
	Entry    string              `config:"entry"`
	Auth     AdminAuthConfig     `config:"auth"`
	Database AdminDatabaseConfig `config:"database"`
}

// AdminAuthConfig holds authentication settings for the admin surface.
type AdminAuthConfig struct {
	Admin AdminRootCredentials `config:"admin"`
}

// AdminRootCredentials is the root account for built-in admin when enabled.
type AdminRootCredentials struct {
	Username string `config:"username"`
	Password string `config:"password"`
}

// AdminDatabaseDriver names the SQL driver used by the admin database.
type AdminDatabaseDriver string

// AdminDatabaseConfig holds DSN settings for the built-in admin database.
type AdminDatabaseConfig struct {
	Driver AdminDatabaseDriver `config:"driver"`
	DSN    string              `config:"dsn"`
}

// ValidateAdmin returns an error when admin is enabled but root credentials are missing.
func (c *Config) ValidateAdmin() error {
	if !c.Admin.Enabled {
		return nil
	}

	u := strings.TrimSpace(c.Admin.Auth.Admin.Username)
	p := strings.TrimSpace(c.Admin.Auth.Admin.Password)
	if u == "" || p == "" {
		return fmt.Errorf("admin: when admin.enabled is true, admin.auth.admin.username and admin.auth.admin.password must both be non-empty")
	}

	return nil
}
