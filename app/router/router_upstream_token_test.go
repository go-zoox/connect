package router_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/router"
	"github.com/go-zoox/connect/user"
	"github.com/go-zoox/crypto/jwt"
	"github.com/go-zoox/zoox/defaults"
)

// testSecretKey is the secret key of the connect apps used by these tests.
const testSecretKey = "test-secret-key"

// startConnect starts an upstream which captures X-Connect-Token and a connect app proxying to it.
// A real server is required, the upstream proxy needs a http.CloseNotifier capable ResponseWriter.
func startConnect(t *testing.T, cfg *config.Config) (connectURL string, capturedToken func() string) {
	t.Helper()

	captured := ""
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Header.Get("X-Connect-Token")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(upstream.Close)

	target, err := url.Parse(upstream.URL)
	if err != nil {
		t.Fatalf("parse upstream url: %v", err)
	}
	port, err := strconv.ParseInt(target.Port(), 10, 64)
	if err != nil {
		t.Fatalf("parse upstream port: %v", err)
	}

	cfg.Upstream = config.UpstreamService{Protocol: target.Scheme, Host: target.Hostname(), Port: port}
	cfg.ApplyDefault()

	app := defaults.Application()
	router.New(app, cfg)

	server := httptest.NewServer(app)
	t.Cleanup(server.Close)

	return server.URL, func() string { return captured }
}

// connectUserOf proxies one request and decodes the X-Connect-Token received by the upstream.
func connectUserOf(t *testing.T, connectURL string, header map[string]string, capturedToken func() string) *user.User {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, connectURL+"/anything", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request connect: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("proxy request: want %d, got %d", http.StatusOK, res.StatusCode)
	}

	token := capturedToken()
	if token == "" {
		t.Fatal("upstream request must carry X-Connect-Token")
	}

	connectUser := &user.User{}
	if err := connectUser.Decode(jwt.New(testSecretKey), token); err != nil {
		t.Fatalf("decode X-Connect-Token: %v", err)
	}
	return connectUser
}

// Upstream mode must carry the current user permissions inside X-Connect-Token, so that upstream
// services resolve roles without calling back into the permissions service.
func TestUpstreamMode_InjectsPermissionsIntoConnectToken(t *testing.T) {
	cfg := &config.Config{SecretKey: testSecretKey, Auth: config.Auth{Mode: "none"}}
	cfg.Services.Permissions.Mode = "local"
	cfg.Services.Permissions.Local = []config.PermissionItem{"ADMIN"}

	connectURL, capturedToken := startConnect(t, cfg)

	connectUser := connectUserOf(t, connectURL, nil, capturedToken)
	if len(connectUser.Permissions) != 1 || connectUser.Permissions[0] != "ADMIN" {
		t.Fatalf("X-Connect-Token permissions want [ADMIN], got %v", connectUser.Permissions)
	}
}

// The permissions answered by the user service are the permissions of the connected application,
// they must be kept and must not be replaced by a second lookup.
func TestUpstreamMode_KeepsUserServicePermissions(t *testing.T) {
	cfg := &config.Config{SecretKey: testSecretKey, Auth: config.Auth{Mode: "password"}}
	cfg.Services.User.Mode = "local"
	cfg.Services.User.Local.Username = "zero"
	cfg.Services.User.Local.Permissions = []string{"USER_PERMISSION"}
	cfg.Services.App.Mode = "local"
	cfg.Services.App.Local.Name = "test-app"
	cfg.Services.Permissions.Mode = "local"
	cfg.Services.Permissions.Local = []config.PermissionItem{"PERMISSIONS_SERVICE_PERMISSION"}

	connectURL, capturedToken := startConnect(t, cfg)

	header := map[string]string{"authorization": "Bearer user-token"}
	connectUser := connectUserOf(t, connectURL, header, capturedToken)
	if len(connectUser.Permissions) != 1 || connectUser.Permissions[0] != "USER_PERMISSION" {
		t.Fatalf("X-Connect-Token permissions want [USER_PERMISSION], got %v", connectUser.Permissions)
	}
}
