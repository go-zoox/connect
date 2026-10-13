package service_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

const testPermissionProvider = "connect-test-permission-provider"

// probeResult collects what the probe handler observed.
type probeResult struct {
	statuses    []int
	messages    []string
	permissions []service.PermissionItem
}

// registerOAuth2Provider registers provider once, GetPermission needs its client credentials.
func registerOAuth2Provider(t *testing.T, provider string) {
	t.Helper()

	if _, err := oauth2.Get(provider); err == nil {
		return
	}

	if err := oauth2.Register(provider, &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}); err != nil {
		t.Fatalf("register oauth2 provider(%s): %v", provider, err)
	}
}

// startProbeServer serves /probe, which resolves the permissions of one token three times.
func startProbeServer(t *testing.T, cfg *config.Config, out *probeResult) string {
	t.Helper()

	app := defaults.Application()
	app.Config.SecretKey = "connect-test-secret-key"

	app.Get("/probe", func(ctx *zoox.Context) {
		for i := 0; i < 3; i++ {
			permissions, status, err := service.GetPermission(ctx, cfg, testPermissionProvider, "token-1")

			out.statuses = append(out.statuses, status)
			if err != nil {
				out.messages = append(out.messages, err.Error())
			}
			out.permissions = permissions
		}

		ctx.Status(http.StatusOK)
	})

	server := httptest.NewServer(app)
	t.Cleanup(server.Close)

	return server.URL
}

// A failed permissions lookup is an external round trip on the request path, it must not be
// retried on every request while the permissions service keeps failing.
func TestGetPermission_CachesFailure(t *testing.T) {
	var calls int64
	permissionsService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&calls, 1)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"token rejected"}`))
	}))
	t.Cleanup(permissionsService.Close)

	registerOAuth2Provider(t, testPermissionProvider)

	cfg := &config.Config{}
	cfg.Services.Permissions.Mode = "service"
	cfg.Services.Permissions.Service = permissionsService.URL

	out := &probeResult{}
	url := startProbeServer(t, cfg, out)

	res, err := http.Get(url + "/probe")
	if err != nil {
		t.Fatalf("probe request: %v", err)
	}
	defer res.Body.Close()

	if got := atomic.LoadInt64(&calls); got != 1 {
		t.Fatalf("failed lookup must be cached: want 1 permissions service call, got %d", got)
	}

	if len(out.statuses) != 3 {
		t.Fatalf("want 3 lookups, got %d", len(out.statuses))
	}
	for i, status := range out.statuses {
		if status != http.StatusUnauthorized {
			t.Fatalf("lookup#%d status: want %d, got %d", i, http.StatusUnauthorized, status)
		}
	}
	if len(out.messages) != 3 {
		t.Fatalf("cached failure must keep returning an error, got %d errors", len(out.messages))
	}
}

// The successful path keeps its cache behavior: one lookup, then the cached answer.
func TestGetPermission_CachesSuccess(t *testing.T) {
	var calls int64
	permissionsService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&calls, 1)
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"result":["USER_PERMISSION"]}`))
	}))
	t.Cleanup(permissionsService.Close)

	registerOAuth2Provider(t, testPermissionProvider)

	cfg := &config.Config{}
	cfg.Services.Permissions.Mode = "service"
	cfg.Services.Permissions.Service = permissionsService.URL

	out := &probeResult{}
	url := startProbeServer(t, cfg, out)

	res, err := http.Get(url + "/probe")
	if err != nil {
		t.Fatalf("probe request: %v", err)
	}
	defer res.Body.Close()

	if got := atomic.LoadInt64(&calls); got != 1 {
		t.Fatalf("successful lookup must be cached: want 1 permissions service call, got %d", got)
	}

	if len(out.permissions) != 1 || out.permissions[0] != "USER_PERMISSION" {
		t.Fatalf("permissions want [USER_PERMISSION], got %v", out.permissions)
	}
}
