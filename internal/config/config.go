package config

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"time"

	goconfig "github.com/go-zoox/config"
)

type Config struct {
	Port                  int64  `config:"port"`
	Mode                  string `config:"mode"`
	SecretKey             string `config:"secret_key"`
	SessionMaxAge         int64  `config:"session_max_age"`
	SessionMaxAgeDuration time.Duration
	LogLevel              string            `config:"log_level"`
	Frontend              ConfigPartService `config:"frontend"`
	Backend               ConfigPartService `config:"backend"`
	//
	OAuth2   []ConfigPartAuthOAuth2 `config:"oauth2"`
	Password ConfigPartAuthPassword `config:"password"`
	//
	Auth ConfigPartAuth `config:"auth"`
	//
	Services ConfigPartServices `config:"services"`
	//
	LoadingHTML string `config:"loading_html"`
	IndexHTML   string `config:"index_html"`
}

type ConfigPartService struct {
	Scheme string `config:"scheme"`
	Host   string `config:"host"`
	Port   int64  `config:"port"`
}

type ConfigPartAuth struct {
	Mode     string `config:"mode"`
	Provider string `config:"provider"`
}

type ConfigPartAuthPassword struct {
	Mode    string                      `config:"mode"`
	Local   ConfigPartAuthPasswordLocal `config:"local"`
	Service string                      `config:"service"`
}

type ConfigPartAuthPasswordLocal struct {
	Username string `config:"username"`
	Password string `config:"password"`
}

type ConfigPartAuthOAuth2 struct {
	Name         string `config:"name"`
	ClientID     string `config:"client_id"`
	ClientSecret string `config:"client_secret"`
	RedirectURI  string `config:"redirect_uri"`
	Scope        string `config:"scope"`
}

type ConfigPartServices struct {
	App   ConfigPartServicesApp   `config:"app"`
	User  ConfigPartServicesUser  `config:"user"`
	Menus ConfigPartServicesMenus `config:"menus"`
	Users ConfigPartServicesUsers `config:"users"`
}

type ConfigPartServicesApp struct {
	Mode  string `config:"mode"`
	Local struct {
		Name        string `config:"name"`
		Logo        string `config:"logo"`
		Description string `config:"description"`
		Settings    struct {
			Functions []string `config:"functions"`
		} `config:"settings"`
	} `config:"local"`
	Service string `config:"service"`
}

type ConfigPartServicesUser struct {
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

type ConfigPartServicesMenus struct {
	Mode    string     `config:"mode"`
	Local   []MenuItem `config:"local"`
	Service string     `config:"service"`
}

type ConfigPartServicesUsers struct {
	Mode    string     `config:"mode"`
	Local   []MenuItem `config:"local"`
	Service string     `config:"service"`
}

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

//
// var isLoaded = false
var cfg Config

func Load() (*Config, error) {
	if err := goconfig.Load(&cfg); err != nil {
		return nil, err
	}

	if cfg.SecretKey == "" {
		return nil, errors.New("secret_key is empty")
	}

	if cfg.SessionMaxAge == 0 {
		// cfg.SessionMaxAge = 7200 * 1000 // 2 hours
		cfg.SessionMaxAgeDuration = 2 * time.Hour
	} else {
		cfg.SessionMaxAgeDuration = time.Duration(cfg.SessionMaxAge) * time.Millisecond
	}

	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "ERROR"
	}
	if cfg.Mode == "" {
		cfg.Mode = "development"
	}

	// default auth mode => oauth2
	if cfg.Auth.Mode == "" {
		cfg.Auth.Mode = "oauth2"
	}
	// default oauth2 provider => doreamon
	if cfg.Auth.Provider == "" {
		cfg.Auth.Provider = "doreamon"
	}

	// default services mode => service
	if cfg.Services.App.Mode == "" {
		cfg.Services.App.Mode = "service"
	}
	if cfg.Services.User.Mode == "" {
		cfg.Services.User.Mode = "service"
	}
	if cfg.Services.Menus.Mode == "" {
		cfg.Services.Menus.Mode = "service"
	}

	applyEnv()

	// isLoaded = true

	return &cfg, nil
}

func applyEnv() {
	if os.Getenv("PORT") != "" {
		v, err := strconv.Atoi(os.Getenv("PORT"))
		if err == nil {
			cfg.Port = int64(v)
		}
	}

	if os.Getenv("MODE") != "" {
		cfg.Mode = os.Getenv("MODE")
	}

	if os.Getenv("SECRET_KEY") != "" {
		cfg.SecretKey = os.Getenv("SECRET_KEY")
	}

	if os.Getenv("SESSION_MAX_AGE") != "" {
		v, err := strconv.Atoi(os.Getenv("SESSION_MAX_AGE"))
		if err == nil {
			cfg.SessionMaxAge = int64(v)
		}
	}

	if os.Getenv("LOG_LEVEL") != "" {
		cfg.LogLevel = os.Getenv("LOG_LEVEL")
	}

	if os.Getenv("AUTH_MODE") != "" {
		cfg.Auth.Mode = os.Getenv("AUTH_MODE")
	}

	if os.Getenv("FRONTEND") != "" {
		u, err := url.Parse(os.Getenv("FRONTEND"))
		if err == nil {
			panic("invalid FRONTEND service: " + os.Getenv("FRONTEND"))
		}

		port, _ := strconv.Atoi(u.Port())
		cfg.Frontend = ConfigPartService{
			Scheme: u.Scheme,
			Host:   u.Hostname(),
			Port:   int64(port),
		}
	}

	if os.Getenv("BACKEND") != "" {
		u, err := url.Parse(os.Getenv("BACKEND"))
		if err == nil {
			panic("invalid BACKEND service: " + os.Getenv("BACKEND"))
		}

		port, _ := strconv.Atoi(u.Port())
		cfg.Frontend = ConfigPartService{
			Scheme: u.Scheme,
			Host:   u.Hostname(),
			Port:   int64(port),
		}
	}

}
