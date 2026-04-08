package config

import (
	"strings"
	"testing"
)

func TestAdminDefaultsDisabledAndEntry(t *testing.T) {
	var c Config
	c.ApplyDefault()

	if c.Admin.Enabled {
		t.Fatalf("expected admin.enabled false by default, got true")
	}
	if c.Admin.Entry != "/admin" {
		t.Fatalf("expected admin.entry %q, got %q", "/admin", c.Admin.Entry)
	}
	if err := c.ValidateAdmin(); err != nil {
		t.Fatalf("ValidateAdmin when disabled: %v", err)
	}
}

func TestAdminEnabledRequiresRootCredentials(t *testing.T) {
	var c Config
	c.ApplyDefault()
	c.Admin.Enabled = true

	err := c.ValidateAdmin()
	if err == nil {
		t.Fatal("expected error when admin.enabled without credentials")
	}
	if !strings.Contains(err.Error(), "admin.auth.admin.username") {
		t.Fatalf("expected error to mention username/password requirement, got: %v", err)
	}

	c.Admin.Auth.Admin.Username = "root"
	c.Admin.Auth.Admin.Password = "secret"
	err = c.ValidateAdmin()
	if err == nil {
		t.Fatal("expected error when admin.enabled with credentials but no database.driver")
	}
	if !strings.Contains(err.Error(), "admin.database.driver") {
		t.Fatalf("expected database.driver requirement, got: %v", err)
	}

	c.Admin.Database.Driver = "sqlite"
	err = c.ValidateAdmin()
	if err == nil {
		t.Fatal("expected error when admin.enabled with driver but no database.dsn")
	}
	if !strings.Contains(err.Error(), "admin.database.dsn") {
		t.Fatalf("expected database.dsn requirement, got: %v", err)
	}

	c.Admin.Database.DSN = "file:validate_admin?mode=memory&cache=shared"
	if err := c.ValidateAdmin(); err != nil {
		t.Fatalf("ValidateAdmin with credentials and database: %v", err)
	}
}

func validAdminEnabledConfig(t *testing.T) Config {
	t.Helper()
	var c Config
	c.ApplyDefault()
	c.Admin.Enabled = true
	c.Admin.Auth.Admin.Username = "root"
	c.Admin.Auth.Admin.Password = "secret"
	c.Admin.Database.Driver = "sqlite"
	c.Admin.Database.DSN = "file:" + strings.ReplaceAll(strings.ReplaceAll(t.Name(), "/", "_"), " ", "_") + "?mode=memory&cache=shared"
	return c
}

func TestValidateAdmin_EntryRejectsAuthBypassOverlaps(t *testing.T) {
	tests := []struct {
		name       string
		mutate     func(*Config)
		wantSubstr string
	}{
		{
			name: "api_prefix",
			mutate: func(c *Config) {
				c.Admin.Entry = "/api"
			},
			wantSubstr: "built-in /api",
		},
		{
			name: "api_prefix_trailing_slash",
			mutate: func(c *Config) {
				c.Admin.Entry = "/api/"
			},
			wantSubstr: "built-in /api",
		},
		{
			name: "double_slash_normalizes_to_root",
			mutate: func(c *Config) {
				c.Admin.Entry = "//"
			},
			wantSubstr: "site root",
		},
		{
			name: "login_reserved",
			mutate: func(c *Config) {
				c.Admin.Entry = "/login"
			},
			wantSubstr: "reserved",
		},
		{
			name: "logout_reserved",
			mutate: func(c *Config) {
				c.Admin.Entry = "/logout"
			},
			wantSubstr: "reserved",
		},
		{
			name: "register_reserved",
			mutate: func(c *Config) {
				c.Admin.Entry = "/register"
			},
			wantSubstr: "reserved",
		},
		{
			name: "equals_backend_prefix",
			mutate: func(c *Config) {
				c.Backend.Prefix = "/internal-api"
				c.Admin.Entry = "/internal-api"
			},
			wantSubstr: "backend.prefix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := validAdminEnabledConfig(t)
			tt.mutate(&c)
			err := c.ValidateAdmin()
			if err == nil {
				t.Fatal("expected ValidateAdmin error for unsafe admin.entry")
			}
			if !strings.Contains(err.Error(), tt.wantSubstr) {
				t.Fatalf("error %q should contain %q", err.Error(), tt.wantSubstr)
			}
		})
	}
}

func TestValidateAdmin_EntryAllowsSafeCustomPaths(t *testing.T) {
	c := validAdminEnabledConfig(t)
	c.Admin.Entry = "/console"
	if err := c.ValidateAdmin(); err != nil {
		t.Fatalf("ValidateAdmin: %v", err)
	}
}
