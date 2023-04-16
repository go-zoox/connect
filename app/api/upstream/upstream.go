package upstream

import (
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/zoox"
)

// var Page = New()

type upstream struct {
	upstream string
	cfg      *config.Config
	//
	handler zoox.HandlerFunc
}

func New(cfg *config.Config) *upstream {
	target := cfg.Upstream.String()

	handler := zoox.WrapH(proxy.NewSingleTarget(target, &proxy.SingleTargetConfig{
		ChangeOrigin: cfg.Upstream.ChangeOrigin,
	}))

	return &upstream{
		upstream: target,
		cfg:      cfg,
		//
		handler: handler,
	}
}

func (p *upstream) Handle(ctx *zoox.Context) {
	p.handler(ctx)
}
