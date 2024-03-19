package github

import (
	"net/url"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/core-utils/cast"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/core-utils/regexp"
	"github.com/go-zoox/random"
)

// Config ...
type Config struct {
	Port          int64
	SecretKey     string
	SessionMaxAge int64
	//
	ClientID     string
	ClientSecret string
	RedirectURI  string
	//
	Frontend string
	Backend  string
	Upstream string
	//
	BackendPrefix                 string
	BackendIsDisablePrefixRewrite bool
	//
	AllowUsernames []string
}

// Create ...
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
	cfgX.Auth.Mode = "oauth2"
	cfgX.Auth.Provider = "github"
	cfgX.Auth.AllowUsernames = cfg.AllowUsernames
	//
	// cfgX.Services.App.Mode = "service"
	// cfgX.Services.App.Service = "https://api.github.com/oauth/app"
	// cfgX.Services.User.Mode = "service"
	// cfgX.Services.User.Service = "https://api.github.com/user"
	// cfgX.Services.Menus.Mode = "service"
	// cfgX.Services.Menus.Service = "https://api.github.com/menus"
	// cfgX.Services.Users.Mode = "service"
	// cfgX.Services.Users.Service = "https://api.github.com/users"
	// cfgX.Services.OpenID.Mode = "service"
	// cfgX.Services.OpenID.Service = "https://api.github.com/oauth/app/user/open_id"
	//
	cfgX.OAuth2 = []config.AuthOAuth2{
		{
			Name:         "github",
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURI:  cfg.RedirectURI,
		},
	}

	if cfg.Upstream != "" {
		if !regexp.Match("^https?://", cfg.Upstream) {
			cfg.Upstream = fmt.Sprintf("http://%s", cfg.Upstream)
		}

		u, err := url.Parse(cfg.Upstream)
		if err != nil {
			return nil, fmt.Errorf("upstream format error, protocol://host:port")
		}

		cfgX.Upstream = config.UpstreamService{
			Protocol: u.Scheme,
			Host:     u.Hostname(),
			Port:     cast.ToInt64(u.Port()),
		}

		if cfgX.Upstream.Port == 0 {
			switch cfgX.Upstream.Protocol {
			case "http":
				cfgX.Upstream.Port = 80
			case "https":
				cfgX.Upstream.Port = 443
			}
		}
	} else {
		if cfg.Frontend == "" || cfg.Backend == "" {
			return nil, fmt.Errorf("frontend and backend are required")
		}

		{
			if !regexp.Match("^https?://", cfg.Frontend) {
				cfg.Frontend = fmt.Sprintf("http://%s", cfg.Frontend)
			}

			u, err := url.Parse(cfg.Frontend)
			if err != nil {
				return nil, fmt.Errorf("frontend format error, protocol://host:port")
			}

			cfgX.Frontend = config.FrontendService{
				Protocol: u.Scheme,
				Host:     u.Hostname(),
				Port:     cast.ToInt64(u.Port()),
			}

			if cfgX.Frontend.Port == 0 {
				switch cfgX.Frontend.Protocol {
				case "http":
					cfgX.Frontend.Port = 80
				case "https":
					cfgX.Frontend.Port = 443
				}
			}
		}

		{
			if !regexp.Match("^https?://", cfg.Backend) {
				cfg.Backend = fmt.Sprintf("http://%s", cfg.Backend)
			}

			u, err := url.Parse(cfg.Backend)
			if err != nil {
				return nil, fmt.Errorf("backend format error, protocol://host:port")
			}

			cfgX.Backend = config.BackendService{
				Protocol: u.Scheme,
				Host:     u.Hostname(),
				Port:     cast.ToInt64(u.Port()),
			}

			if cfgX.Backend.Port == 0 {
				switch cfgX.Backend.Protocol {
				case "http":
					cfgX.Backend.Port = 80
				case "https":
					cfgX.Backend.Port = 443
				}
			}
		}

		if cfg.BackendPrefix != "" {
			cfgX.Backend.Prefix = cfg.BackendPrefix
		}

		if cfg.BackendIsDisablePrefixRewrite {
			cfgX.Backend.IsDisablePrefixRewrite = cfg.BackendIsDisablePrefixRewrite
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
