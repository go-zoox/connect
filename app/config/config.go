package config

import (
	"fmt"
	"time"

	"github.com/go-zoox/core-utils/regexp"

	goconfig "github.com/go-zoox/config"
)

// Config ...
type Config struct {
	// Port is the port of the server, default: 8080
	Port int64 `config:"port"`

	// Mode is the mode of the server, default: development
	// options: development, production
	Mode string `config:"mode"`

	// SecretKey is the secret key of the server, default: random string
	SecretKey string `config:"secret_key"`

	// SessionMaxAge is the max age of session, unit: second, default: 86400
	SessionMaxAge int64 `config:"session_max_age"`
	//
	sessionMaxAgeDuration time.Duration

	// StaticCacheControlMaxAge is the max age of static cache control, unit: second, default: 7d
	StaticCacheControlMaxAge int64

	// LogLevel is the log level of the server, default: info
	LogLevel string `config:"log_level"`

	// Proxy Mode:
	// - 1. frontend + backend
	// - 2. upstream

	// S1: Connect => Frontend + Backend
	Frontend FrontendService `config:"frontend"`
	Backend  BackendService  `config:"backend"`
	// S2: Connect => Upstream
	Upstream UpstreamService `config:"upstream"`

	// Auth is the auth config
	Auth Auth `config:"auth"`
	// OAuth2 is the oauth2 config
	OAuth2 []AuthOAuth2 `config:"oauth2"`
	// AuthPassword is the auth password config
	Password AuthPassword `config:"password"`

	// Services is the services config, e.g: app, user, menus, users, open_id
	Services Services `config:"services"`

	// LoadingHTML is the loading html, default: <html><body>Loading...</body></html>
	// Only used in development mode
	LoadingHTML string `config:"loading_html"`

	// IndexHTML is the index html, default: <html><body>Index...</body></html>
	IndexHTML string `config:"index_html"`

	// Routes is the routes config, e.g: /api/app, /api/user
	Routes []Route `config:"routes"`

	// BuiltInAPIs is the built in apis config, e.g: app, user, menus, users, config
	BuiltInAPIs BuiltInAPIs `config:"built_in_apis"`
}

// GetSessionMaxAgeDuration gets the session max age duration.
func (c *Config) GetSessionMaxAgeDuration() time.Duration {
	if c.sessionMaxAgeDuration == 0 {
		return DefaultMaxSessionAgeInSecond * time.Second
	}

	return c.sessionMaxAgeDuration
}

// SetSessionMaxAgeDuration sets the session max age duration.
func (c *Config) SetSessionMaxAgeDuration(d time.Duration) {
	c.sessionMaxAgeDuration = d
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

// Services ...
type Services struct {
	App         ServicesApp         `config:"app"`
	User        ServicesUser        `config:"user"`
	Menus       ServicesMenus       `config:"menus"`
	Permissions ServicesPermissions `config:"permissions"`
	Users       ServicesUsers       `config:"users"`
	OpenID      ServicesOpenID      `config:"open_id"`
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

// ServicesPermissions ...
type ServicesPermissions struct {
	Mode    string           `config:"mode"`
	Local   []PermissionItem `config:"local"`
	Service string           `config:"service"`
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
	App         string `config:"app"`
	User        string `config:"user"`
	Menus       string `config:"menus"`
	Permissions string `config:"permissions"`
	Users       string `config:"users"`
	Config      string `config:"config"`
	//
	QRCode string `config:"qrcode"`
	//
	Login string `config:"login"`
	//
	Public string `config:"built_in"`
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

// PermissionItem ...
type PermissionItem = string

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

	if cfg.Upstream.Port == 0 {
		if cfg.Upstream.Protocol == "https" {
			cfg.Upstream.Port = 443
		} else if cfg.Upstream.Protocol == "http" {
			cfg.Upstream.Port = 80
		}
	}

	if s.Protocol == "https" && s.Port == 443 {
		return fmt.Sprintf("%s://%s", s.Protocol, s.Host)
	}

	if s.Protocol == "http" && s.Port == 80 {
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

// fixUpstream fix upstream url
// e.g: localhost:8080 => http://localhost:8080
func fixUpstream(upstream string) string {
	if !regexp.Match("^https?://", upstream) {
		return fmt.Sprintf("http://%s", upstream)
	}

	return upstream
}
