package router_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-zoox/connect/admin/api"
	"github.com/go-zoox/connect/admin/bootstrap"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/router"
	"github.com/go-zoox/zoox/defaults"
)

func testBuiltInIgnorePaths(cfg *config.Config) []string {
	return []string{
		fmt.Sprintf("^/api%s$", cfg.BuiltInAPIs.Roles),
		fmt.Sprintf("^/api%s$", cfg.BuiltInAPIs.Groups),
		fmt.Sprintf("^/api%s%s$", cfg.BuiltInAPIs.Public, cfg.BuiltInAPIs.Roles),
		fmt.Sprintf("^/api%s%s$", cfg.BuiltInAPIs.Public, cfg.BuiltInAPIs.Groups),
	}
}

func newTestApp(t *testing.T, adminEnabled bool, ignoreRolesAndGroups bool) (http.Handler, *config.Config) {
	t.Helper()

	cfg := &config.Config{
		Auth: config.Auth{Mode: "password"},
		Admin: config.Admin{
			Enabled: adminEnabled,
		},
	}
	cfg.ApplyDefault()
	if ignoreRolesAndGroups {
		cfg.Auth.IgnorePaths = append(cfg.Auth.IgnorePaths, testBuiltInIgnorePaths(cfg)...)
	}
	app := defaults.Application()
	router.New(app, cfg)
	return app, cfg
}

func newTestAppWithAdminBootstrap(t *testing.T) (http.Handler, *config.Config) {
	t.Helper()
	dsn := "file:" + strings.ReplaceAll(strings.ReplaceAll(t.Name(), "/", "_"), " ", "_") + "?mode=memory&cache=shared"
	cfg := &config.Config{
		Auth: config.Auth{Mode: "password"},
		Admin: config.Admin{
			Enabled: true,
			Auth: config.AdminAuthConfig{
				Admin: config.AdminRootCredentials{
					Username: "root",
					Password: "secret-root-pass",
				},
			},
			Database: config.AdminDatabaseConfig{
				Driver: "sqlite",
				DSN:    dsn,
			},
		},
	}
	cfg.ApplyDefault()
	cfg.Auth.IgnorePaths = append(cfg.Auth.IgnorePaths, testBuiltInIgnorePaths(cfg)...)

	if err := bootstrap.Bootstrap(cfg); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	app := defaults.Application()
	router.New(app, cfg)
	return app, cfg
}

func getBody(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	b, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(b)
}

func TestAdminEnabled_RegistersRolesAndGroupsUnderAPI(t *testing.T) {
	h, cfg := newTestApp(t, true, true)

	for _, path := range []string{"/api/roles", "/api/groups"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Accept", "application/json")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code == http.StatusNotFound {
			t.Fatalf("%s: got 404, expected a registered handler (not missing route)", path)
		}
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/roles", nil)
	req.Header.Set("Accept", "application/json")
	h.ServeHTTP(rec, req)
	body := getBody(t, rec)
	if !strings.Contains(body, api.MarkerRoles) {
		t.Fatalf("expected admin roles marker in body, got %q (cfg admin enabled=%v)", body, cfg.Admin.Enabled)
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/api/groups", nil)
	req2.Header.Set("Accept", "application/json")
	h.ServeHTTP(rec2, req2)
	body2 := getBody(t, rec2)
	if !strings.Contains(body2, api.MarkerGroups) {
		t.Fatalf("expected admin groups marker in body, got %q", body2)
	}
}

func TestAdminEnabled_DoesNotExposeAdminDataUnderPublicBuiltinPrefix(t *testing.T) {
	h, _ := newTestApp(t, true, true)
	pub := cfgDefaultPublic()
	for _, suffix := range []string{
		"/roles", "/groups", "/users", "/user", "/app", "/menus", "/permissions",
	} {
		path := "/api" + pub + suffix
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Accept", "application/json")
		h.ServeHTTP(rec, req)
		body := getBody(t, rec)
		if strings.Contains(body, api.MarkerRoles) {
			t.Fatalf("%s: must not expose admin roles handler on public prefix, body %q", path, body)
		}
		if strings.Contains(body, api.MarkerGroups) {
			t.Fatalf("%s: must not expose admin groups handler on public prefix, body %q", path, body)
		}
		if strings.Contains(body, api.MarkerUsers) {
			t.Fatalf("%s: must not expose admin users handler on public prefix, body %q", path, body)
		}
	}
}

func cfgDefaultPublic() string {
	c := &config.Config{}
	c.ApplyDefault()
	return c.BuiltInAPIs.Public
}

func TestAdminEnabled_LoginRejectsInvalidCredentials(t *testing.T) {
	h, _ := newTestAppWithAdminBootstrap(t)

	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"username":"root","password":"wrong-password"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for bad password, got %d body %q", rec.Code, getBody(t, rec))
	}
	body := getBody(t, rec)
	if !strings.Contains(body, "Invalid username or password") {
		t.Fatalf("expected invalid-credentials message in body, got %q", body)
	}
}

func TestAdminEnabled_LoginAcceptsSeededRoot(t *testing.T) {
	h, cfg := newTestAppWithAdminBootstrap(t)
	loginPath := "/api" + cfg.BuiltInAPIs.Login

	payload := fmt.Sprintf(`{"username":%q,"password":%q}`, "root", "secret-root-pass")
	req := httptest.NewRequest(http.MethodPost, loginPath, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid admin login, got %d body %q", rec.Code, getBody(t, rec))
	}
}

func TestAdminDisabled_DoesNotRegisterBuiltinRolesHandler(t *testing.T) {
	// Do not ignore /api/roles: auth rejects before the catch-all backend proxy (httptest recorder lacks CloseNotify).
	h, _ := newTestApp(t, false, false)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/roles", nil)
	req.Header.Set("Accept", "application/json")
	h.ServeHTTP(rec, req)

	// Without the admin-only route, the request should not be served by the stub that emits MarkerRoles.
	body := getBody(t, rec)
	if strings.Contains(body, api.MarkerRoles) {
		t.Fatalf("did not expect admin roles marker when admin is disabled, got body %q", body)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("without session, want %d for protected built-in API when admin disabled, got %d body %q", http.StatusUnauthorized, rec.Code, body)
	}
}

func TestAdminStatic_WithPasswordAuth_AdminShellWithoutSession_ProtectedAPIStillUnauthorized(t *testing.T) {
	h, _ := newTestAppWithAdminBootstrap(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /admin/ without session: want %d, got %d body %q", http.StatusOK, rec.Code, getBody(t, rec))
	}
	shell := getBody(t, rec)
	if !strings.Contains(shell, `data-connect-builtin-admin="1"`) {
		snippetLen := 200
		if len(shell) < snippetLen {
			snippetLen = len(shell)
		}
		t.Fatalf("expected built-in admin HTML marker in body, got snippet %q", shell[:snippetLen])
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req2.Header.Set("Accept", "application/json")
	h.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/users without session: want %d (must not be ignored like admin entry), got %d body %q", http.StatusUnauthorized, rec2.Code, getBody(t, rec2))
	}
}
