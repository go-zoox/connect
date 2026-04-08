package router_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-zoox/connect/app/admin/api"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/router"
	"github.com/go-zoox/zoox/defaults"
)

func testBuiltInIgnorePaths() []string {
	return []string{
		`^/api/roles$`,
		`^/api/groups$`,
		`^/api/_/roles$`,
		`^/api/_/groups$`,
	}
}

func newTestApp(t *testing.T, adminEnabled bool, ignoreRolesAndGroups bool) (http.Handler, *config.Config) {
	t.Helper()

	app := defaults.Application()
	auth := config.Auth{Mode: "password"}
	if ignoreRolesAndGroups {
		auth.IgnorePaths = testBuiltInIgnorePaths()
	}
	cfg := &config.Config{
		Auth: auth,
		Admin: config.Admin{
			Enabled: adminEnabled,
		},
	}
	cfg.ApplyDefault()
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

func TestAdminEnabled_RegistersRolesAndGroupsUnderPublicPrefix(t *testing.T) {
	h, _ := newTestApp(t, true, true)
	public := "/api" + cfgDefaultPublic() + "/roles"

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, public, nil)
	req.Header.Set("Accept", "application/json")
	h.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Fatalf("%s: got 404, expected a registered handler", public)
	}
	body := getBody(t, rec)
	if !strings.Contains(body, api.MarkerRoles) {
		t.Fatalf("expected admin roles marker in body for public path, got %q", body)
	}
}

func cfgDefaultPublic() string {
	c := &config.Config{}
	c.ApplyDefault()
	return c.BuiltInAPIs.Public
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
	if rec.Code == http.StatusOK {
		t.Fatalf("expected non-200 without admin roles route (auth or proxy), got 200 body %q", body)
	}
}
