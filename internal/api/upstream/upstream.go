package upstream

import (
	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/headers"
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
	return &upstream{
		upstream: cfg.Upstream.String(),
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
		ctx.Request.Header.Set(headers.XRequestID, ctx.RequestID())

		if cfg.Mode == "production" {
			zoox.WrapH(proxy.NewSingleTarget(p.upstream, &proxy.SingleTargetConfig{
				ChangeOrigin: cfg.Upstream.ChangeOrigin,
			}))(ctx)
			return
		}

		if !p.isHealth() {
			ctx.String(200, cfg.LoadingHTML)
			return
		}

		ctx.Request.Header.Set("cache-control", "no-cache")

		zoox.WrapH(proxy.NewSingleTarget(p.upstream, &proxy.SingleTargetConfig{
			ChangeOrigin: cfg.Upstream.ChangeOrigin,
		}))(ctx)
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
