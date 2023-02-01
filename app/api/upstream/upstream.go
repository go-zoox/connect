package upstream

import (
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/zoox"
)

// var Page = New()

type upstream struct {
	upstream string
	cfg      *config.Config
}

func New(cfg *config.Config) *upstream {
	protocol := cfg.Upstream.Protocol
	host := cfg.Upstream.Host
	port := cfg.Upstream.Port

	if protocol == "" {
		protocol = "http"
	}

	if host == "" {
		host = "127.0.0.1"
	}

	if port == 0 {
		port = 8000
	}
	target := fmt.Sprintf(
		"%s://%s:%d",
		protocol,
		host,
		port,
	)

	return &upstream{
		upstream: target,
		cfg:      cfg,
	}
}

func (p *upstream) isHealth() bool {
	response, err := fetch.Get(p.upstream, &fetch.Config{
		Headers: map[string]string{
			"accept": "text/html",
		},
	})

	if err != nil || response.Status != 200 {
		logger.Debug("Check health: (URI: %s, error: %s)", p.upstream, err)
		return false
	}

	return true
}

func (p *upstream) RenderStatic() func(ctx *zoox.Context) {
	return zoox.WrapH(proxy.NewSingleTarget(p.upstream))
}

func (p *upstream) RenderPage() func(ctx *zoox.Context) {
	cfg := p.cfg

	return func(ctx *zoox.Context) {
		// request id
		ctx.Request.Header.Set(zoox.RequestIDHeader, ctx.RequestID())

		if cfg.Mode == "production" {
			// ctx.Status(200)
			// ctx.String(200, cfg.IndexHTML)
			zoox.WrapH(proxy.NewSingleTarget(p.upstream))(ctx)
			return
		}

		if !p.isHealth() {
			// ctx.Render(200, "loading.html", nil)
			ctx.String(200, cfg.LoadingHTML)
			return
		}

		ctx.Request.Header.Set("cache-control", "no-cache")

		zoox.WrapH(proxy.NewSingleTarget(p.upstream))(ctx)
	}
}

func (p *upstream) Health(cfg *config.Config) func(ctx *zoox.Context) {
	return func(ctx *zoox.Context) {
		if !p.isHealth() {
			ctx.Status(503)
			return
		}

		ctx.Status(200)
	}
}
