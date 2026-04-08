package app_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-zoox/connect/app"
	"github.com/go-zoox/connect/app/admin/api"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/errors"
)

func readBody(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	b, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(b)
}

// Integration: full Connect.Setup path (defaults, validation, bootstrap, router, oauth2 register)
// must preserve legacy behavior when admin is off.

func TestIntegration_ConnectSetup_DefaultMode_BackwardCompatible(t *testing.T) {
	c := app.New()
	cfg := &config.Config{
		Auth: config.Auth{Mode: "password"},
	}
	cfg.ApplyDefault()
	if cfg.Admin.Enabled {
		t.Fatal("expected default admin.enabled false")
	}

	if err := c.Setup(cfg); err != nil {
		t.Fatalf("Setup: %v", err)
	}
	h := c.TestHTTPHandler()
	if h == nil {
		t.Fatal("TestHTTPHandler: nil")
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/roles", nil)
	req.Header.Set("Accept", "application/json")
	h.ServeHTTP(rec, req)
	body := readBody(t, rec)
	if strings.Contains(body, api.MarkerRoles) {
		t.Fatalf("default mode must not serve admin roles stub; body %q", body)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/roles without session: want %d, got %d body %q", http.StatusUnauthorized, rec.Code, body)
	}

	loginPath := "/api" + cfg.BuiltInAPIs.Login
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, loginPath, strings.NewReader(`{"username":"x","password":"y"}`))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Accept", "application/json")
	h.ServeHTTP(rec2, req2)
	body2 := readBody(t, rec2)
	if strings.Contains(body2, errors.AdminLoginInvalid.Message) {
		t.Fatalf("default mode must use core login (captcha), not admin login; body %q", body2)
	}
	if !strings.Contains(body2, errors.InvalidCaptcha.Message) {
		t.Fatalf("core login should require captcha first; want %q in body, got %q", errors.InvalidCaptcha.Message, body2)
	}
}

func TestIntegration_ConnectSetup_AdminEnabled_LoginAndProtectedAPI(t *testing.T) {
	c := app.New()
	dsn := "file:" + strings.ReplaceAll(strings.ReplaceAll(t.Name(), "/", "_"), " ", "_") + "?mode=memory&cache=shared"
	cfg := &config.Config{
		Auth: config.Auth{Mode: "password"},
		Admin: config.Admin{
			Enabled: true,
			Auth: config.AdminAuthConfig{
				Admin: config.AdminRootCredentials{
					Username: "root",
					Password: "integration-secret",
				},
			},
			Database: config.AdminDatabaseConfig{
				Driver: "sqlite",
				DSN:    dsn,
			},
		},
	}
	cfg.ApplyDefault()
	// Allow unauthenticated access to roles/groups for this test (same as router tests).
	cfg.Auth.IgnorePaths = append(cfg.Auth.IgnorePaths,
		fmt.Sprintf("^/api%s$", cfg.BuiltInAPIs.Roles),
		fmt.Sprintf("^/api%s$", cfg.BuiltInAPIs.Groups),
	)

	if err := c.Setup(cfg); err != nil {
		t.Fatalf("Setup: %v", err)
	}
	h := c.TestHTTPHandler()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"username":"root","password":"integration-secret"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("admin login: want 200, got %d body %q", rec.Code, readBody(t, rec))
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/api/roles", nil)
	req2.Header.Set("Accept", "application/json")
	for _, c := range rec.Result().Cookies() {
		req2.AddCookie(c)
	}
	h.ServeHTTP(rec2, req2)
	if rec2.Code == http.StatusNotFound {
		t.Fatal("GET /api/roles: unexpected 404")
	}
	body := readBody(t, rec2)
	if !strings.Contains(body, api.MarkerRoles) {
		t.Fatalf("after login, GET /api/roles should hit admin handler; body %q", body)
	}
}
