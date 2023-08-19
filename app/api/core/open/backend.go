package open

import (
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/proxy/utils/rewriter"
	"github.com/go-zoox/zoox"
)

// New ...
func New(cfg *config.Config) func(*zoox.Context) {
	// @TODO
	if cfg.Upstream.IsValid() {
		cfg.Backend.Protocol = cfg.Upstream.Protocol
		cfg.Backend.Host = cfg.Upstream.Host
		cfg.Backend.Port = cfg.Upstream.Port
	}

	backend := cfg.Backend.String()
	var rewrites rewriter.Rewriters
	if !cfg.Backend.IsDisablePrefixRewrite {
		rewrites = rewriter.Rewriters{
			{
				From: fmt.Sprintf("^%s/open/(.*)", cfg.Backend.Prefix),
				To:   "/open/$1",
			},
		}
	} else {
		rewrites = rewriter.Rewriters{
			{
				From: fmt.Sprintf("^%s/open/(.*)", cfg.Backend.Prefix),
				To:   fmt.Sprintf("%s/open/$1", cfg.Backend.Prefix),
			},
		}
	}

	return zoox.WrapH(proxy.NewSingleHost(backend, &proxy.SingleHostConfig{
		// Rewrites: rewriter.Rewriters{
		// 	{
		// 		From: "^/api/open/(.*)",
		// 		To:   fmt.Sprintf("%s/open/$1", prefix),
		// 	},
		// },
		Rewrites:     rewrites,
		ChangeOrigin: cfg.Backend.ChangeOrigin,
	}))
}
