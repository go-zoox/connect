package app

import (
	"embed"
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/router"
	"github.com/go-zoox/oauth2"

	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

//go:embed web
var StaticFS embed.FS

type Connect struct {
	core *zoox.Application
	cfg  *config.Config
	//
	apiHandler    func(*zoox.Application)
	staticHandler func(*zoox.Application) string
}

func New() *Connect {
	app := defaults.Application()

	// @TODO Test with api docs at 127.0.0.1
	// app.Use(middleware.CORS())

	return &Connect{
		core: app,
	}
}

func (e *Connect) RegisterApi(fn func(*zoox.Application)) {
	e.apiHandler = fn
}

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
	e.cfg = cfg

	// 1. register loading
	e.registerLoading()
	// 2. register oauth2
	e.registerOauth2()

	router.New(e.core, e.cfg)
}

func (e *Connect) Start(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	e.handle(cfg)

	return e.core.Run(fmt.Sprintf(":%d", e.cfg.Port))
}
