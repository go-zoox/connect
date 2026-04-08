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
	if err := c.ValidateAdmin(); err != nil {
		t.Fatalf("ValidateAdmin with credentials: %v", err)
	}
}
