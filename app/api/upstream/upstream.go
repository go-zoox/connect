package upstream

import (
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/zoox"
)

// var Page = New()

// Upstream ...
type Upstream struct {
	upstream string
	cfg      *config.Config
	//
	handler zoox.HandlerFunc
}

// New ...
func New(cfg *config.Config) *Upstream {
	target := cfg.Upstream.String()

	handler := zoox.WrapH(proxy.NewSingleHost(target, &proxy.SingleHostConfig{
		ChangeOrigin: cfg.Upstream.ChangeOrigin,
	}))

	return &Upstream{
		upstream: target,
		cfg:      cfg,
		//
		handler: handler,
	}
}

// Handle ...
func (p *Upstream) Handle(ctx *zoox.Context) {
	p.handler(ctx)
}
