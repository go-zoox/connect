package doreamon

import (
	"testing"

	"github.com/go-zoox/core-utils/fmt"

	"github.com/go-zoox/testify"
)

func TestCreate(t *testing.T) {
	cfg, err := Create(&Config{
		Port:          8080,
		SecretKey:     "7777",
		SessionMaxAge: 3600,
		ClientID:      "client_id",
		ClientSecret:  "client_secret",
		RedirectURI:   "https://doreamon.example.com/login/doreamon/callback",
		Upstream:      "http://upstream:8080",
	})

	fmt.PrintJSON(map[string]any{
		"cfg": cfg,
		"err": err,
	})

	testify.Assert(t, err == nil, "failed to create doreamon")
	testify.Equal(t, 8080, cfg.Port)
	testify.Equal(t, "7777", cfg.SecretKey)
	testify.Equal(t, 3600, cfg.SessionMaxAge)
	testify.Equal(t, "doreamon", cfg.Auth.Provider)
	testify.Equal(t, "service", cfg.Services.App.Mode)
	testify.Equal(t, "https://api.zcorky.com/oauth/app", cfg.Services.App.Service)
	testify.Equal(t, "service", cfg.Services.User.Mode)
	testify.Equal(t, "https://api.zcorky.com/user", cfg.Services.User.Service)
	testify.Equal(t, "service", cfg.Services.Menus.Mode)
	testify.Equal(t, "https://api.zcorky.com/menus", cfg.Services.Menus.Service)
	testify.Equal(t, "service", cfg.Services.Users.Mode)
	testify.Equal(t, "https://api.zcorky.com/users", cfg.Services.Users.Service)
	testify.Equal(t, "service", cfg.Services.OpenID.Mode)
	testify.Equal(t, "https://api.zcorky.com/oauth/app/user/open_id", cfg.Services.OpenID.Service)
	testify.Equal(t, "doreamon", cfg.OAuth2[0].Name)
	testify.Equal(t, "client_id", cfg.OAuth2[0].ClientID)
	testify.Equal(t, "client_secret", cfg.OAuth2[0].ClientSecret)
	testify.Equal(t, "https://doreamon.example.com/login/doreamon/callback", cfg.OAuth2[0].RedirectURI)
	// testify.Equal(t, "http://upstream:8080", fmt.Sprintf(""))
}
