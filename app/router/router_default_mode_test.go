package router_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-zoox/connect/app/admin/api"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/errors"
	"github.com/go-zoox/connect/app/router"
	"github.com/go-zoox/zoox/defaults"
)

// Default mode: admin is off unless explicitly enabled. These tests lock in backward-compatible
// routing when the admin block is omitted or admin.enabled is false.
// Absence of the built-in admin shell at /admin/ is covered by TestAdminStatic_Disabled_NoBuiltinUIAtDefaultEntry.

func TestDefaultMode_AdminDisabledWithoutExplicitAdminConfig(t *testing.T) {
	cfg := &config.Config{
		Auth: config.Auth{Mode: "password"},
	}
	cfg.ApplyDefault()
	if cfg.Admin.Enabled {
		t.Fatal("expected default admin.enabled false when admin config is omitted")
	}
	if cfg.Admin.Entry != "/admin" {
		t.Fatalf("ApplyDefault should set admin.entry to %q, got %q", "/admin", cfg.Admin.Entry)
	}

	app := defaults.Application()
	router.New(app, cfg)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/roles", nil)
	req.Header.Set("Accept", "application/json")
	app.ServeHTTP(rec, req)
	body := getBody(t, rec)
	if strings.Contains(body, api.MarkerRoles) {
		t.Fatalf("default mode must not register admin roles handler; got body containing %q", api.MarkerRoles)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated /api/roles: want %d, got %d body %q", http.StatusUnauthorized, rec.Code, body)
	}
}

func TestDefaultMode_LoginUsesCoreHandlerNotAdminLogin(t *testing.T) {
	cfg := &config.Config{
		Auth: config.Auth{Mode: "password"},
	}
	cfg.ApplyDefault()
	if cfg.Admin.Enabled {
		t.Fatal("expected admin.enabled false")
	}

	app := defaults.Application()
	router.New(app, cfg)

	loginPath := "/api" + cfg.BuiltInAPIs.Login
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, loginPath, strings.NewReader(`{"username":"nouser","password":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	app.ServeHTTP(rec, req)
	body := getBody(t, rec)
	if strings.Contains(body, errors.AdminLoginInvalid.Message) {
		t.Fatalf("default mode must register core POST /api login (captcha gate), not admin login; body %q", body)
	}
	if !strings.Contains(body, errors.InvalidCaptcha.Message) {
		t.Fatalf("core login should reject missing captcha before credentials; want %q in body, got %q", errors.InvalidCaptcha.Message, body)
	}
}
