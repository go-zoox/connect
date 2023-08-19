package page

import (
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/headers"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/zoox"
)

// Page ...
type Page struct {
	frontend string
	cfg      *config.Config
}

// New ...
func New(cfg *config.Config) *Page {
	frontend := cfg.Frontend.String()

	return &Page{
		frontend: frontend,
		cfg:      cfg,
	}
}

func (p *Page) isHealth() bool {
	response, err := fetch.Get(p.frontend, &fetch.Config{
		Headers: map[string]string{
			"accept": "text/html",
		},
	})

	if err != nil || response.Status != 200 {
		logger.Debug("Check health: (URI: %s, error: %s)", p.frontend, err)
		return false
	}

	return true
}

// RenderStatic ...
func (p *Page) RenderStatic() func(ctx *zoox.Context) {
	return zoox.WrapH(proxy.NewSingleHost(p.frontend))
}

// RenderPage ...
func (p *Page) RenderPage() func(ctx *zoox.Context) {
	cfg := p.cfg

	return func(ctx *zoox.Context) {
		// request id
		ctx.Request.Header.Set(headers.XRequestID, ctx.RequestID())

		if cfg.Mode == "production" {
			zoox.WrapH(proxy.NewSingleHost(p.frontend, &proxy.SingleHostConfig{
				ChangeOrigin: cfg.Frontend.ChangeOrigin,
			}))(ctx)
			return
		}

		if !p.isHealth() {
			ctx.String(200, cfg.LoadingHTML)
			return
		}

		ctx.Request.Header.Set("cache-control", "no-cache")

		zoox.WrapH(proxy.NewSingleHost(p.frontend, &proxy.SingleHostConfig{
			ChangeOrigin: cfg.Frontend.ChangeOrigin,
		}))(ctx)
	}
}

// Health ...
func (p *Page) Health(cfg *config.Config) func(ctx *zoox.Context) {
	return func(ctx *zoox.Context) {
		if !p.isHealth() {
			ctx.Status(503)
			return
		}

		ctx.Status(200)
	}
}
