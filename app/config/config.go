package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-zoox/core-utils/strings"

	goconfig "github.com/go-zoox/config"
	"github.com/go-zoox/random"
)

// Config ...
type Config struct {
	Port      int64  `config:"port"`
	Mode      string `config:"mode"`
	SecretKey string `config:"secret_key"`
	// SessionMaxAge is the max age of session, unit: second, default: 86400
	SessionMaxAge         int64 `config:"session_max_age"`
	SessionMaxAgeDuration time.Duration
	LogLevel              string `config:"log_level"`
	// S1: Connect => Frontend + Backend
	Frontend FrontendService `config:"frontend"`
	Backend  BackendService  `config:"backend"`
	// S2: Connect => Upstream
	Upstream UpstreamService `config:"upstream"`
	//
	OAuth2   []AuthOAuth2 `config:"oauth2"`
	Password AuthPassword `config:"password"`
	//
	Auth Auth `config:"auth"`
	//
	Services Services `config:"services"`
	//
	LoadingHTML string `config:"loading_html"`
	IndexHTML   string `config:"index_html"`
	//
	Routes []Route `config:"routes"`
	//
	BuiltInAPIs BuiltInAPIs `config:"built_in_apis"`
}

// FrontendService ...
type FrontendService struct {
	Protocol string `config:"protocol"`
	Host     string `config:"host"`
	Port     int64  `config:"port"`
	//
	ChangeOrigin bool `config:"change_origin"`
}

// BackendService ...
type BackendService struct {
	Protocol string `config:"protocol"`
	Host     string `config:"host"`
	Port     int64  `config:"port"`
	// Prefix is the backend prefix, default: /api
	Prefix                 string `config:"prefix"`
	IsDisablePrefixRewrite bool   `config:"is_disable_prefix_rewrite"`
	//
	ChangeOrigin bool `config:"change_origin"`
}

// UpstreamService ...
type UpstreamService struct {
	Protocol string `config:"protocol"`
	Host     string `config:"host"`
	Port     int64  `config:"port"`
	//
	ChangeOrigin bool `config:"change_origin"`
}

// Auth ...
type Auth struct {
	Mode     string `config:"mode"`
	Provider string `config:"provider"`
}

// AuthPassword ...
type AuthPassword struct {
	Mode    string            `config:"mode"`
	Local   AuthPasswordLocal `config:"local"`
	Service string            `config:"service"`
}

// AuthPasswordLocal ...
type AuthPasswordLocal struct {
	Username string `config:"username"`
	Password string `config:"password"`
}

// AuthOAuth2 ...
type AuthOAuth2 struct {
	Name         string `config:"name"`
	ClientID     string `config:"client_id"`
	ClientSecret string `config:"client_secret"`
	RedirectURI  string `config:"redirect_uri"`
	Scope        string `config:"scope"`
}

// Services ...
type Services struct {
	App    ServicesApp    `config:"app"`
	User   ServicesUser   `config:"user"`
	Menus  ServicesMenus  `config:"menus"`
	Users  ServicesUsers  `config:"users"`
	OpenID ServicesOpenID `config:"open_id"`
}

// ServicesApp ...
type ServicesApp struct {
	Mode  string `config:"mode"`
	Local struct {
		Name        string `config:"name"`
		Logo        string `config:"logo"`
		Description string `config:"description"`
		Settings    struct {
			Functions any `config:"functions"`
		} `config:"settings"`
	} `config:"local"`
	Service string `config:"service"`
}

// ServicesUser ...
type ServicesUser struct {
	Mode  string `config:"mode"`
	Local struct {
		ID          string   `config:"id"`
		Username    string   `config:"username"`
		Avatar      string   `config:"avatar"`
		Nickname    string   `config:"nickname"`
		Email       string   `config:"email"`
		Permissions []string `config:"permissions"`
	} `config:"local"`
	Service string `config:"service"`
}

// ServicesMenus ...
type ServicesMenus struct {
	Mode    string     `config:"mode"`
	Local   []MenuItem `config:"local"`
	Service string     `config:"service"`
}

// ServicesUsers ...
type ServicesUsers struct {
	Mode    string     `config:"mode"`
	Local   []MenuItem `config:"local"`
	Service string     `config:"service"`
}

// ServicesOpenID ...
type ServicesOpenID struct {
	Mode  string `config:"mode"`
	Local struct {
		OpenID string `config:"open_id"`
	} `config:"local"`
	Service string `config:"service"`
}

// Route ...
type Route struct {
	Path    string       `config:"path"`
	Backend RouteBackend `config:"backend"`
}

// RouteBackend ...
type RouteBackend struct {
	ServiceName     string `config:"service_name"`
	ServicePort     int64  `config:"service_port"`
	ServiceProtocol string `config:"service_protocol"`
	//
	DisableRewrite bool `config:"disable_rewrite"`
	//
	SecretKey string `config:"secret_key"`
}

// BuiltInAPIs ...
type BuiltInAPIs struct {
	App    string `config:"app"`
	User   string `config:"user"`
	Menus  string `config:"menus"`
	Users  string `config:"users"`
	Config string `config:"config"`
	//
	QRCode string `config:"qrcode"`
	//
	Login string `config:"login"`
}

// MenuItem ...
type MenuItem struct {
	ID         string `config:"id"`
	Name       string `config:"name"`
	Path       string `config:"path"`
	Icon       string `config:"icon"`
	Sort       int64  `config:"sort"`
	IsHidden   bool   `config:"hidden"`
	IsExpanded bool   `config:"expanded"`
	Layout     string `config:"layout"`
	IFrame     string `config:"iframe"`
	Redirect   string `config:"redirect"`
}

// String ...
func (s *FrontendService) String() string {
	if s.Protocol == "" {
		s.Protocol = "http"
	}

	if s.Host == "" {
		s.Host = "127.0.0.1"
	}

	if s.Port == 0 {
		s.Port = 8000
	}

	if cfg.Upstream.Protocol == "https" && cfg.Upstream.Port == 0 {
		cfg.Upstream.Port = 443
	}

	if s.Protocol == "https" && s.Port == 443 {
		return fmt.Sprintf("%s://%s", s.Protocol, s.Host)
	}

	return fmt.Sprintf("%s://%s:%d", s.Protocol, s.Host, s.Port)
}

// String ...
func (s *BackendService) String() string {
	if s.Protocol == "" {
		s.Protocol = "http"
	}

	if s.Host == "" {
		s.Host = "127.0.0.1"
	}

	if s.Port == 0 {
		s.Port = 8001
	}

	if cfg.Upstream.Protocol == "https" && cfg.Upstream.Port == 0 {
		cfg.Upstream.Port = 443
	}

	if s.Protocol == "https" && s.Port == 443 {
		return fmt.Sprintf("%s://%s", s.Protocol, s.Host)
	}

	return fmt.Sprintf("%s://%s:%d", s.Protocol, s.Host, s.Port)
}

// String ...
func (s *UpstreamService) String() string {
	if s.Protocol == "" {
		s.Protocol = "http"
	}

	if cfg.Upstream.Protocol == "https" && cfg.Upstream.Port == 0 {
		cfg.Upstream.Port = 443
	}

	if s.Protocol == "https" && s.Port == 443 {
		return fmt.Sprintf("%s://%s", s.Protocol, s.Host)
	}

	return fmt.Sprintf("%s://%s:%d", s.Protocol, s.Host, s.Port)
}

// IsValid ...
func (s *UpstreamService) IsValid() bool {
	return s.Host != "" && s.Port != 0
}

// String ...
func (s *RouteBackend) String() string {
	if s.ServiceProtocol == "" {
		s.ServiceProtocol = "http"
	}

	if s.ServicePort == 0 {
		panic(fmt.Errorf("service_port is required"))
	}

	return fmt.Sprintf("%s://%s:%d", s.ServiceProtocol, s.ServiceName, s.ServicePort)
}

var cfg Config

// Load loads config from file
func Load(cfgFile string) (*Config, error) {
	if err := goconfig.Load(&cfg, &goconfig.LoadOptions{
		FilePath: cfgFile,
	}); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadFromService loads config from service
func LoadFromService(fn func() (string, error)) (*Config, error) {
	var cfg Config
	if err := goconfig.LoadFromService(&cfg, fn); err != nil {
		// return nil, fmt.Errorf("load config from service error: %s", err)
		return nil, err
	}

	return &cfg, nil
}

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
	}

	if os.Getenv("AUTH_MODE") != "" {
		c.Auth.Mode = os.Getenv("AUTH_MODE")
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

		if c.Frontend.Protocol == "https" && c.Frontend.Port == 0 {
			c.Frontend.Port = 443
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

		if c.Backend.Protocol == "https" && c.Backend.Port == 0 {
			c.Backend.Port = 443
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

		if c.Upstream.Protocol == "https" && c.Upstream.Port == 0 {
			c.Upstream.Port = 443
		}
	}

	if os.Getenv("DISABLE_PREFIX_REWRITE") != "" {
		c.Backend.IsDisablePrefixRewrite = true
	}

	if c.Port == 0 {
		c.Port = 8080
	}
	if c.LogLevel == "" {
		c.LogLevel = "ERROR"
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
		c.SessionMaxAge = 86400
	}

	if c.Backend.Prefix == "" {
		c.Backend.Prefix = "/api"
	}

	c.SessionMaxAgeDuration = time.Duration(c.SessionMaxAge) * time.Second

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
}

// fixUpstream fix upstream url
// e.g: localhost:8080 => http://localhost:8080
func fixUpstream(upstream string) string {
	parts := strings.Split(upstream, ":")
	if len(parts) <= 2 {
		return "http://" + upstream
	}

	return upstream
}
