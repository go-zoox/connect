package doreamon

import (
	"net/url"
	"strings"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/core-utils/cast"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/core-utils/regexp"
	"github.com/go-zoox/random"
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
		if regexp.Match("://", cfg.Upstream) {
			u, err := url.Parse(cfg.Upstream)
			if err != nil {
				return nil, fmt.Errorf("upstream format error, protocol://host:port")
			}

			cfgX.Upstream = config.ConfigUpstreamService{
				Protocol: u.Scheme,
				Host:     u.Hostname(),
				Port:     cast.ToInt64(u.Port()),
			}
		} else {
			parts := strings.Split(cfg.Upstream, ":")
			if len(parts) != 2 {
				return nil, fmt.Errorf("upstream format error, host:port")
			}

			cfgX.Upstream = config.ConfigUpstreamService{
				Protocol: "http",
				Host:     parts[0],
				Port:     cast.ToInt64(parts[1]),
			}
		}
	} else {
		if cfg.Frontend == "" || cfg.Backend == "" {
			return nil, fmt.Errorf("frontend and backend are required")
		}

		{
			if regexp.Match("://", cfg.Frontend) {
				u, err := url.Parse(cfg.Frontend)
				if err != nil {
					return nil, fmt.Errorf("frontend format error, protocol://host:port")
				}

				cfgX.Frontend = config.ConfigFrontendService{
					Protocol: u.Scheme,
					Host:     u.Hostname(),
					Port:     cast.ToInt64(u.Port()),
				}
			} else {
				parts := strings.Split(cfg.Frontend, ":")
				if len(parts) != 2 {
					return nil, fmt.Errorf("frontend format error, host:port")
				}

				cfgX.Frontend = config.ConfigFrontendService{
					Host: parts[0],
					Port: cast.ToInt64(parts[1]),
				}
			}
		}

		{
			if regexp.Match("://", cfg.Backend) {
				u, err := url.Parse(cfg.Backend)
				if err != nil {
					return nil, fmt.Errorf("backend format error, protocol://host:port")
				}

				cfgX.Backend = config.ConfigBackendService{
					Protocol: u.Scheme,
					Host:     u.Hostname(),
					Port:     cast.ToInt64(u.Port()),
				}
			} else {
				parts := strings.Split(cfg.Backend, ":")
				if len(parts) != 2 {
					return nil, fmt.Errorf("backend format error, host:port")
				}

				cfgX.Backend = config.ConfigBackendService{
					Host: parts[0],
					Port: cast.ToInt64(parts[1]),
				}
			}
		}
	}

	if cfgX.SecretKey == "" {
		cfgX.SecretKey = random.String(10)
	}

	if cfgX.SessionMaxAge == 0 {
		cfgX.SessionMaxAge = 1 * 24 * 60 * 60
	}

	return cfgX, nil
}
