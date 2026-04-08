package config

import (
	"fmt"
	"regexp"
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

// NormalizedAdminEntry returns the URL path prefix for the built-in admin UI: leading slash, no trailing slash.
func NormalizedAdminEntry(entry string) string {
	e := strings.TrimSpace(entry)
	if e == "" {
		return "/admin"
	}
	if !strings.HasPrefix(e, "/") {
		e = "/" + e
	}
	e = strings.TrimSuffix(e, "/")
	if e == "" {
		return "/admin"
	}
	return e
}

// AdminStaticAuthIgnorePattern is a regex matching the admin UI entry path so auth middleware allows the static shell without a session.
func AdminStaticAuthIgnorePattern(entry string) string {
	p := NormalizedAdminEntry(entry)
	return "^" + regexp.QuoteMeta(p) + "(?:/|$)"
}

// ValidateAdmin returns an error when admin is enabled but required settings are missing.
func (c *Config) ValidateAdmin() error {
	if !c.Admin.Enabled {
		return nil
	}

	u := strings.TrimSpace(c.Admin.Auth.Admin.Username)
	p := strings.TrimSpace(c.Admin.Auth.Admin.Password)
	if u == "" || p == "" {
		return fmt.Errorf("admin: when admin.enabled is true, admin.auth.admin.username and admin.auth.admin.password must both be non-empty")
	}

	driver := strings.TrimSpace(string(c.Admin.Database.Driver))
	if driver == "" {
		return fmt.Errorf("admin: when admin.enabled is true, admin.database.driver must be non-empty")
	}

	dsn := strings.TrimSpace(c.Admin.Database.DSN)
	if dsn == "" {
		return fmt.Errorf("admin: when admin.enabled is true, admin.database.dsn must be non-empty")
	}

	return nil
}
