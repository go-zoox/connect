package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

const testFeishuProvider = "connect-test-feishu-provider"

// startUserService serves the user of a connect app without inline permissions, which is the
// case where the permissions service is asked for every proxied request.
func startUserService(t *testing.T) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{
				"id":       "1",
				"username": "zero",
				"nickname": "Zero",
				"email":    "zero@example.com",
			},
		})
	}))
	t.Cleanup(server.Close)

	return server
}

// Without an oauth2 provider in the session the feishu open id lookup can never succeed, it must
// be skipped: it used to run on every uncached user lookup and to sleep 3 seconds on failure.
func TestGetUser_SkipsFeishuOpenIDWithoutOAuth2Provider(t *testing.T) {
	var openIDCalls int64
	openIDService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&openIDCalls, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(openIDService.Close)

	userService := startUserService(t)

	cfg := &config.Config{}
	cfg.Services.User.Mode = "service"
	cfg.Services.User.Service = userService.URL
	cfg.Services.OpenID.Mode = "service"
	cfg.Services.OpenID.Service = openIDService.URL

	var (
		elapsed      time.Duration
		username     string
		feishuOpenID string
		userErr      error
	)

	app := defaults.Application()
	app.Config.SecretKey = "connect-test-secret-key"
	app.Get("/probe", func(ctx *zoox.Context) {
		start := time.Now()
		user, _, err := service.GetUser(ctx, cfg, "token-1")
		elapsed = time.Since(start)
		userErr = err

		if user != nil {
			username = user.Username
			feishuOpenID = user.FeishuOpenID
		}

		ctx.Status(http.StatusOK)
	})

	server := httptest.NewServer(app)
	t.Cleanup(server.Close)

	res, err := http.Get(server.URL + "/probe")
	if err != nil {
		t.Fatalf("probe request: %v", err)
	}
	defer res.Body.Close()

	if userErr != nil {
		t.Fatalf("GetUser: %v", userErr)
	}
	if username != "zero" {
		t.Fatalf("username: want zero, got %s", username)
	}
	if feishuOpenID != "" {
		t.Fatalf("feishu open id: want empty, got %s", feishuOpenID)
	}
	if got := atomic.LoadInt64(&openIDCalls); got != 0 {
		t.Fatalf("without an oauth2 provider the feishu open id lookup must be skipped, got %d call(s)", got)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("GetUser must not block on the feishu open id lookup, took %s", elapsed)
	}
}

// The feishu open id is still resolved for an oauth2 session, the skip must not remove the feature.
func TestGetUser_ResolvesFeishuOpenIDWithOAuth2Provider(t *testing.T) {
	var openIDCalls int64
	openIDService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&openIDCalls, 1)
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"result":{"openID":"ou_123"}}`))
	}))
	t.Cleanup(openIDService.Close)

	registerOAuth2Provider(t, testFeishuProvider)

	userService := startUserService(t)

	cfg := &config.Config{}
	cfg.Services.User.Mode = "service"
	cfg.Services.User.Service = userService.URL
	cfg.Services.OpenID.Mode = "service"
	cfg.Services.OpenID.Service = openIDService.URL

	var (
		username     string
		feishuOpenID string
		userErr      error
	)

	app := defaults.Application()
	app.Config.SecretKey = "connect-test-secret-key"
	// an oauth2 login sets the provider of the session
	app.Get("/session/provider", func(ctx *zoox.Context) {
		service.SetProvider(ctx, cfg, testFeishuProvider)
		ctx.Status(http.StatusOK)
	})
	app.Get("/probe", func(ctx *zoox.Context) {
		user, _, err := service.GetUser(ctx, cfg, "token-2")
		userErr = err

		if user != nil {
			username = user.Username
			feishuOpenID = user.FeishuOpenID
		}

		ctx.Status(http.StatusOK)
	})

	server := httptest.NewServer(app)
	t.Cleanup(server.Close)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookie jar: %v", err)
	}
	client := &http.Client{Jar: jar}

	loginRes, err := client.Get(server.URL + "/session/provider")
	if err != nil {
		t.Fatalf("session provider request: %v", err)
	}
	loginRes.Body.Close()

	res, err := client.Get(server.URL + "/probe")
	if err != nil {
		t.Fatalf("probe request: %v", err)
	}
	defer res.Body.Close()

	if userErr != nil {
		t.Fatalf("GetUser: %v", userErr)
	}
	if username != "zero" {
		t.Fatalf("username: want zero, got %s", username)
	}
	if feishuOpenID != "ou_123" {
		t.Fatalf("feishu open id: want ou_123, got %q", feishuOpenID)
	}
	if got := atomic.LoadInt64(&openIDCalls); got != 1 {
		t.Fatalf("feishu open id lookup: want 1 call, got %d", got)
	}
}
