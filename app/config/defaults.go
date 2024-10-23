package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-zoox/random"
)

// ApplyDefault applies default config
func (c *Config) ApplyDefault() {
	if os.Getenv("PORT") != "" {
		v, err := strconv.Atoi(os.Getenv("PORT"))
		if err == nil {
			c.Port = int64(v)
		}
	}

	if os.Getenv("MODE") != "" {
		c.Mode = os.Getenv("MODE")
	}

	if os.Getenv("SECRET_KEY") != "" {
		c.SecretKey = os.Getenv("SECRET_KEY")
	}
	if c.SecretKey == "" {
		c.SecretKey = random.String(16)
	}

	if os.Getenv("SESSION_MAX_AGE") != "" {
		v, err := strconv.Atoi(os.Getenv("SESSION_MAX_AGE"))
		if err == nil {
			c.SessionMaxAge = int64(v)
		}
	}

	if os.Getenv("LOG_LEVEL") != "" {
		c.LogLevel = os.Getenv("LOG_LEVEL")
	} else {
		c.LogLevel = "error"
	}

	if os.Getenv("AUTH_MODE") != "" {
		c.Auth.Mode = os.Getenv("AUTH_MODE")
	}

	if os.Getenv("AUTH_PROVIDER") != "" {
		c.Auth.Provider = os.Getenv("AUTH_PROVIDER")
	}

	if os.Getenv("AUTH_IGNORE_PATHS") != "" {
		c.Auth.IgnorePaths = append(c.Auth.IgnorePaths, strings.Split(os.Getenv("AUTH_IGNORE_PATHS"), ",")...)
	}

	if os.Getenv("AUTH_IS_IGNORE_PATHS_DISABLED") == "true" {
		c.Auth.IsIgnorePathsDisabled = true
	}

	if os.Getenv("AUTH_IS_IGNORE_WHEN_HEADER_AUTHORIZATION_FOUND") == "true" {
		c.Auth.IsIgnoreWhenHeaderAuthorizationFound = true
	}

	if c.Frontend.Host == "" && os.Getenv("FRONTEND") != "" {
		v := fixUpstream(os.Getenv("FRONTEND"))
		u, err := url.Parse(v)
		if err != nil {
			panic(fmt.Sprintf("invalid FRONTEND service(%s): %s", os.Getenv("FRONTEND"), err.Error()))
		}

		port, _ := strconv.Atoi(u.Port())
		c.Frontend = FrontendService{
			Protocol: u.Scheme,
			Host:     u.Hostname(),
			Port:     int64(port),
		}

		if c.Frontend.Port == 0 {
			switch c.Frontend.Protocol {
			case "http":
				c.Frontend.Port = 80
			case "https":
				c.Frontend.Port = 443
			}
		}
	}

	if c.Backend.Host == "" && os.Getenv("BACKEND") != "" {
		v := fixUpstream(os.Getenv("BACKEND"))
		u, err := url.Parse(v)
		if err != nil {
			panic(fmt.Sprintf("invalid BACKEND service(%s): %s", os.Getenv("BACKEND"), err.Error()))
		}

		port, _ := strconv.Atoi(u.Port())
		c.Backend = BackendService{
			Protocol:               u.Scheme,
			Host:                   u.Hostname(),
			Port:                   int64(port),
			Prefix:                 "/api",
			IsDisablePrefixRewrite: false,
		}

		if c.Backend.Port == 0 {
			switch c.Backend.Protocol {
			case "http":
				c.Backend.Port = 80
			case "https":
				c.Backend.Port = 443
			}
		}
	}

	if os.Getenv("UPSTREAM") != "" {
		v := fixUpstream(os.Getenv("UPSTREAM"))
		u, err := url.Parse(v)
		if err != nil {
			panic(fmt.Sprintf("invalid UPSTREAM service(%s): %s", os.Getenv("UPSTREAM"), err.Error()))
		}

		port, _ := strconv.Atoi(u.Port())
		c.Upstream = UpstreamService{
			Protocol: u.Scheme,
			Host:     u.Hostname(),
			Port:     int64(port),
		}

		if c.Upstream.Port == 0 {
			switch c.Upstream.Protocol {
			case "http":
				c.Upstream.Port = 80
			case "https":
				c.Upstream.Port = 443
			}
		}
	}

	if os.Getenv("DISABLE_PREFIX_REWRITE") != "" {
		c.Backend.IsDisablePrefixRewrite = true
	}

	if c.Port == 0 {
		c.Port = 8080
	}

	if c.Mode == "" {
		c.Mode = "development"
	}

	// default auth mode => oauth2
	if c.Auth.Mode == "" {
		c.Auth.Mode = "oauth2"
	}
	// default oauth2 provider => doreamon
	if c.Auth.Provider == "" {
		c.Auth.Provider = "doreamon"
	}

	// default services mode => service
	if c.Services.App.Mode == "" {
		c.Services.App.Mode = "service"
	}
	if c.Services.App.Service == "" {
		c.Services.App.Service = "https://api.zcorky.com/oauth/app"
	}

	if c.Services.User.Mode == "" {
		c.Services.User.Mode = "service"
	}
	if c.Services.User.Service == "" {
		c.Services.User.Service = "https://api.zcorky.com/user"
	}

	if c.Services.Menus.Mode == "" {
		c.Services.Menus.Mode = "service"
	}
	if c.Services.Menus.Service == "" {
		c.Services.Menus.Service = "https://api.zcorky.com/menus"
	}

	if c.Services.Users.Mode == "" {
		c.Services.Users.Mode = "service"
	}
	if c.Services.Users.Service == "" {
		c.Services.Users.Service = "https://api.zcorky.com/users"
	}

	if c.Services.OpenID.Mode == "" {
		c.Services.OpenID.Mode = "service"
	}
	if c.Services.OpenID.Service == "" {
		c.Services.OpenID.Service = "https://api.zcorky.com/oauth/app/user/open_id"
	}

	if c.SessionMaxAge == 0 {
		c.SessionMaxAge = DefaultMaxSessionAgeInSecond
	}

	if c.Backend.Prefix == "" {
		c.Backend.Prefix = "/api"
	}

	c.SetSessionMaxAgeDuration(time.Duration(c.SessionMaxAge) * time.Second)

	// built in apis
	if c.BuiltInAPIs.App == "" {
		if os.Getenv("BUILT_IN_APIS_APP") != "" {
			c.BuiltInAPIs.App = os.Getenv("BUILT_IN_APIS_APP")
		} else {
			c.BuiltInAPIs.App = "/app"
		}
	}
	if c.BuiltInAPIs.User == "" {
		if os.Getenv("BUILT_IN_APIS_USER") != "" {
			c.BuiltInAPIs.User = os.Getenv("BUILT_IN_APIS_USER")
		} else {
			c.BuiltInAPIs.User = "/user"
		}
	}
	if c.BuiltInAPIs.Menus == "" {
		if os.Getenv("BUILT_IN_APIS_MENUS") != "" {
			c.BuiltInAPIs.Menus = os.Getenv("BUILT_IN_APIS_MENUS")
		} else {
			c.BuiltInAPIs.Menus = "/menus"
		}
	}
	if c.BuiltInAPIs.Users == "" {
		if os.Getenv("BUILT_IN_APIS_USERS") != "" {
			c.BuiltInAPIs.Users = os.Getenv("BUILT_IN_APIS_USERS")
		} else {
			c.BuiltInAPIs.Users = "/users"
		}
	}
	if c.BuiltInAPIs.Config == "" {
		if os.Getenv("BUILT_IN_APIS_CONFIG") != "" {
			c.BuiltInAPIs.Config = os.Getenv("BUILT_IN_APIS_CONFIG")
		} else {
			c.BuiltInAPIs.Config = "/config"
		}
	}
	if c.BuiltInAPIs.QRCode == "" {
		if os.Getenv("BUILT_IN_APIS_QRCODE") != "" {
			c.BuiltInAPIs.QRCode = os.Getenv("BUILT_IN_APIS_QRCODE")
		} else {
			c.BuiltInAPIs.QRCode = "/qrcode"
		}
	}
	if c.BuiltInAPIs.Login == "" {
		if os.Getenv("BUILT_IN_APIS_LOGIN") != "" {
			c.BuiltInAPIs.Login = os.Getenv("BUILT_IN_APIS_LOGIN")
		} else {
			c.BuiltInAPIs.Login = "/login"
		}
	}
	if c.BuiltInAPIs.Public == "" {
		if os.Getenv("BUILT_IN_APIS_BUILT_IN") != "" {
			c.BuiltInAPIs.Public = os.Getenv("BUILT_IN_APIS_BUILT_IN")
		} else {
			c.BuiltInAPIs.Public = "/_"
		}
	}
}
