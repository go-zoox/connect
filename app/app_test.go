package app_test

import (
	"strings"
	"testing"

	"github.com/go-zoox/connect/app"
	"github.com/go-zoox/connect/app/config"
)

// Ensures Connect.Start runs the same validation path as production (after ApplyDefault),
// not only direct unit tests on Config.ValidateAdmin.
func TestStart_RejectsAdminEnabledWithoutCredentials(t *testing.T) {
	c := app.New()
	cfg := &config.Config{
		Admin: config.Admin{Enabled: true},
	}
	err := c.Start(cfg)
	if err == nil {
		t.Fatal("expected startup to fail when admin is enabled without credentials")
	}
	if !strings.Contains(err.Error(), "config validation failed") {
		t.Fatalf("expected wrapped validation error, got: %v", err)
	}
}
