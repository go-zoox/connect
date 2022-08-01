package page

import (
	"fmt"

	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/zoox"
)

// var Page = New()

type page struct {
	frontend string
	cfg      *config.Config
}

func New(cfg *config.Config) *page {
	scheme := cfg.Frontend.Scheme
	host := cfg.Frontend.Host
	port := cfg.Frontend.Port

	if scheme == "" {
		scheme = "http"
	}

	if host == "" {
		host = "127.0.0.1"
	}

	if port == 0 {
		port = 8000
	}
	frontend := fmt.Sprintf(
		"%s://%s:%d",
		scheme,
		host,
		port,
	)

	return &page{
		frontend: frontend,
		cfg:      cfg,
	}
}

func (p *page) isHealth() bool {
	response, err := fetch.Get(p.frontend, &fetch.Config{
		Headers: map[string]string{
			"accept": "text/html",
		},
	})

	if err != nil || response.Status != 200 {
		logger.Debug("Check health: ", p.frontend, err)
		return false
	}

	return true
}

func (p *page) RenderStatic() func(ctx *zoox.Context) {
	return zoox.WrapH(proxy.NewSingleTarget(p.frontend))
}

func (p *page) RenderPage() func(ctx *zoox.Context) {
	cfg := p.cfg

	return func(ctx *zoox.Context) {
		// request id
		ctx.Request.Header.Set(zoox.RequestIDHeader, ctx.RequestID())

		if cfg.Mode == "production" {
			// ctx.Status(200)
			// ctx.String(200, cfg.IndexHTML)
			zoox.WrapH(proxy.NewSingleTarget(p.frontend))(ctx)
			return
		}

		if !p.isHealth() {
			// ctx.Render(200, "loading.html", nil)
			ctx.String(200, cfg.LoadingHTML)
			return
		}

		ctx.Request.Header.Set("cache-control", "no-cache")

		zoox.WrapH(proxy.NewSingleTarget(p.frontend))(ctx)
	}
}

func (p *page) Health(cfg *config.Config) func(ctx *zoox.Context) {
	return func(ctx *zoox.Context) {
		if !p.isHealth() {
			ctx.Status(503)
			return
		}

		ctx.Status(200)
	}
}
