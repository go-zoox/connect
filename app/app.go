package app

import (
	"embed"
	"net/http"
	"os"

	"github.com/go-zoox/chalk"
	"github.com/go-zoox/connect"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/router"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/debug"
	"github.com/go-zoox/oauth2"

	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

// StaticFS ...
//
//go:embed web
var StaticFS embed.FS

// Connect ...
type Connect struct {
	core *zoox.Application
	cfg  *config.Config
	//
	apiHandler    func(*zoox.Application)
	staticHandler func(*zoox.Application) string
}

// New creates a new Connect instance.
func New() *Connect {
	app := defaults.Application()

	return &Connect{
		core: app,
	}
}

// RegisterAPI registers the api handler.
func (e *Connect) RegisterAPI(fn func(*zoox.Application)) {
	e.apiHandler = fn
}

// RegisterStatic registers the static handler.
func (e *Connect) RegisterStatic(fn func(*zoox.Application) string) {
	e.staticHandler = fn
}

func (e *Connect) registerLoading() {
	loadingHTML, _ := StaticFS.ReadFile("web/loading.html")
	e.cfg.LoadingHTML = string(loadingHTML)
}

func (e *Connect) registerOauth2() {
	for _, cfg := range e.cfg.OAuth2 {
		oauth2.Register(cfg.Name, &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Scope:        cfg.Scope,
			RedirectURI:  cfg.RedirectURI,
		})
	}
}

func (e *Connect) handle(cfg *config.Config) {
	// @TODO
	cfg.ApplyDefault()

	if debug.IsDebugMode() {
		fmt.PrintJSON("connect config:", cfg)
	}

	e.cfg = cfg

	e.core.Config.Banner = fmt.Sprintf(`
  _____       ____                  _____                       __ 
 / ___/__    /_  / ___  ___ __ __  / ___/__  ___  ___  ___ ____/ /_
/ (_ / _ \    / /_/ _ \/ _ \\ \ / / /__/ _ \/ _ \/ _ \/ -_) __/ __/
\___/\___/   /___/\___/\___/_\_\  \___/\___/_//_/_//_/\__/\__/\__/ 
                                                               
The Lighweight, Powerful Auth Connect (Version: %s)

____________________________________O/_______
                                    O\
	`, chalk.Green(connect.Version))

	e.core.Config.LogLevel = cfg.LogLevel

	e.core.Config.SecretKey = cfg.SecretKey

	e.core.Config.Session.MaxAge = cfg.SessionMaxAgeDuration
	// iframe
	e.core.Config.Session.Secure = true
	e.core.Config.Session.SameSite = http.SameSiteNoneMode

	// 1. register loading
	e.registerLoading()
	// 2. register oauth2
	e.registerOauth2()

	router.New(e.core, e.cfg)
}

// Start starts the Connect server.
func (e *Connect) Start(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if os.Getenv("LOG_LEVEL") == "debug" {
		fmt.PrintJSON("config:", cfg)
	}

	e.handle(cfg)

	return e.core.Run(fmt.Sprintf(":%d", e.cfg.Port))
}
