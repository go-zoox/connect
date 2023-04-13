package doreamon

import (
	"net/url"

	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/debug"
	"github.com/go-zoox/random"
	"github.com/spf13/cast"
)

type Config struct {
	Port          int64
	SecretKey     string
	SessionMaxAge int64
	ClientID      string
	ClientSecret  string
	RedirectURI   string
	Frontend      string
	Backend       string
	Upstream      string
}

func Create(cfg *Config) (*config.Config, error) {
	if cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.RedirectURI == "" {
		return nil, fmt.Errorf("client_id, client_secret, redirect_uri are required")
	}

	if cfg.Upstream != "" {
		// ok
	} else if cfg.Frontend != "" && cfg.Backend != "" {
		// ok
	} else {
		return nil, fmt.Errorf("upstream or frontend and backend are required")
	}

	cfgX := &config.Config{
		Port:          cfg.Port,
		SecretKey:     cfg.SecretKey,
		SessionMaxAge: cfg.SessionMaxAge,
	}
	cfgX.Auth.Provider = "doreamon"
	cfgX.Services.App.Mode = "service"
	cfgX.Services.App.Service = "https://api.zcorky.com/oauth/app"
	cfgX.Services.User.Mode = "service"
	cfgX.Services.User.Service = "https://api.zcorky.com/user"
	cfgX.Services.Menus.Mode = "service"
	cfgX.Services.Menus.Service = "https://api.zcorky.com/menus"
	cfgX.Services.Users.Mode = "service"
	cfgX.Services.Users.Service = "https://api.zcorky.com/users"
	cfgX.Services.OpenID.Mode = "service"
	cfgX.Services.OpenID.Service = "https://api.zcorky.com/oauth/app/user/open_id"
	//
	cfgX.OAuth2 = []config.ConfigPartAuthOAuth2{
		{
			Name:         "doreamon",
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURI:  cfg.RedirectURI,
		},
	}

	if cfg.Upstream != "" {
		u, err := url.Parse(config.FixUpstream(cfg.Upstream))
		if err != nil {
			return nil, fmt.Errorf("upstream format error, protocol://host:port")
		}

		cfgX.Upstream = config.ConfigUpstreamService{
			Protocol: u.Scheme,
			Host:     u.Hostname(),
			Port:     cast.ToInt64(u.Port()),
		}
	} else {
		if cfg.Frontend == "" || cfg.Backend == "" {
			return nil, fmt.Errorf("frontend and backend are required")
		}

		{
			u, err := url.Parse(config.FixUpstream((cfg.Frontend)))
			if err != nil {
				return nil, fmt.Errorf("frontend format error, protocol://host:port")
			}

			cfgX.Frontend = config.ConfigFrontendService{
				Protocol: u.Scheme,
				Host:     u.Hostname(),
				Port:     cast.ToInt64(u.Port()),
			}
		}

		{
			u, err := url.Parse(config.FixUpstream((cfg.Frontend)))
			if err != nil {
				return nil, fmt.Errorf("frontend format error, protocol://host:port")
			}

			cfgX.Frontend = config.ConfigFrontendService{
				Protocol: u.Scheme,
				Host:     u.Hostname(),
				Port:     cast.ToInt64(u.Port()),
			}
		}
	}

	if cfgX.SecretKey == "" {
		cfgX.SecretKey = random.String(10)
	}

	if debug.IsDebugMode() {
		fmt.PrintJSON("config:", cfgX)
	}

	return cfgX, nil
}
