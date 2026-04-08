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

// normalizedBackendPrefix returns backend.prefix as a non-empty URL path: leading slash, no trailing slash.
func normalizedBackendPrefix(c *Config) string {
	p := strings.TrimSpace(c.Backend.Prefix)
	if p == "" {
		return "/api"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	p = strings.TrimSuffix(p, "/")
	if p == "" {
		return "/api"
	}
	return p
}

func (c *Config) validateAdminEntryForAuthSafety() error {
	entry := NormalizedAdminEntry(c.Admin.Entry)
	if entry == "/" {
		return fmt.Errorf("admin: admin.entry %q resolves to %q, which is not allowed: the admin auth ignore rule would treat the site root as unauthenticated", strings.TrimSpace(c.Admin.Entry), entry)
	}
	if entry == "/api" {
		return fmt.Errorf("admin: admin.entry must not be %q (normalized: %q): it would skip authentication for the entire built-in /api route tree", strings.TrimSpace(c.Admin.Entry), entry)
	}
	backend := normalizedBackendPrefix(c)
	if entry == backend {
		return fmt.Errorf("admin: admin.entry must not equal backend.prefix %q: that path carries proxied API traffic and cannot be used as the admin UI entry", backend)
	}
	for _, reserved := range []string{"/login", "/logout", "/register"} {
		if entry == reserved {
			return fmt.Errorf("admin: admin.entry must not be %q: that path is reserved for the application authentication UI", entry)
		}
	}
	return nil
}

// ValidateAdmin returns an error when admin is enabled but required settings are missing.
func (c *Config) ValidateAdmin() error {
	if !c.Admin.Enabled {
		return nil
	}

	if err := c.validateAdminEntryForAuthSafety(); err != nil {
		return err
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
