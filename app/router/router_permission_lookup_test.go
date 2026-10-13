package router_test

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/zoox"
)

const testRouterProvider = "connect-test-router-provider"

// A user whose permissions are not answered by the user service makes every proxied request
// resolve permissions against the permissions service (router.fillPermissions). When that
// service fails, the failure must be cached: otherwise every keystroke of that account pays an
// extra external round trip (measured as ~0.6s per request in production).
func TestUpstreamMode_CachesFailingPermissionLookup(t *testing.T) {
	var permissionCalls int64
	permissionsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&permissionCalls, 1)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"token rejected"}`))
	}))
	t.Cleanup(permissionsServer.Close)

	userServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"result":{"id":"1","username":"zero","email":"zero@example.com"}}`))
	}))
	t.Cleanup(userServer.Close)

	if _, err := oauth2.Get(testRouterProvider); err != nil {
		if err := oauth2.Register(testRouterProvider, &oauth2.Config{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		}); err != nil {
			t.Fatalf("register oauth2 provider(%s): %v", testRouterProvider, err)
		}
	}

	cfg := &config.Config{SecretKey: testSecretKey, Auth: config.Auth{Mode: "oauth2"}}
	cfg.Services.User.Mode = "service"
	cfg.Services.User.Service = userServer.URL
	cfg.Services.App.Mode = "local"
	cfg.Services.App.Local.Name = "test-app"
	cfg.Services.Permissions.Mode = "service"
	cfg.Services.Permissions.Service = permissionsServer.URL
	cfg.Services.OpenID.Mode = "local"

	connectURL, _ := startConnect(t, cfg, func(app *zoox.Application) {
		// the session must survive requests: it carries the oauth2 provider of the account
		app.Config.SecretKey = testSecretKey
		app.Get("/__test/session_provider", func(ctx *zoox.Context) {
			service.SetProvider(ctx, cfg, testRouterProvider)
			ctx.Status(http.StatusOK)
		})
	})

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookie jar: %v", err)
	}
	client := &http.Client{Jar: jar}

	// act as an oauth2 login of that account
	loginReq, err := http.NewRequest(http.MethodGet, connectURL+"/__test/session_provider", nil)
	if err != nil {
		t.Fatalf("build session request: %v", err)
	}
	loginReq.Header.Set("authorization", "Bearer user-token")

	loginRes, err := client.Do(loginReq)
	if err != nil {
		t.Fatalf("session request: %v", err)
	}
	loginRes.Body.Close()
	if loginRes.StatusCode != http.StatusOK {
		t.Fatalf("session request: want %d, got %d", http.StatusOK, loginRes.StatusCode)
	}

	// three proxied requests of that account
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest(http.MethodGet, connectURL+"/anything", nil)
		if err != nil {
			t.Fatalf("build request#%d: %v", i, err)
		}
		req.Header.Set("authorization", "Bearer user-token")

		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("proxy request#%d: %v", i, err)
		}
		res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("proxy request#%d: want %d, got %d", i, http.StatusOK, res.StatusCode)
		}
	}

	if got := atomic.LoadInt64(&permissionCalls); got != 1 {
		t.Fatalf("a failing permission lookup must be cached: want 1 permissions service call for 3 requests, got %d", got)
	}
}
