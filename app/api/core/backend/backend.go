package user

import (
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/proxy/utils/rewriter"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) func(*zoox.Context) {
	backend := cfg.Backend.String()
	rewrites := rewriter.Rewriters{}
	if !cfg.Backend.IsDisablePrefixRewrite {
		rewrites = rewriter.Rewriters{
			{
				From: fmt.Sprintf("^%s/(.*)", cfg.Backend.Prefix),
				To:   "/$1",
			},
		}
	}

	return zoox.WrapH(proxy.NewSingleTarget(backend, &proxy.SingleTargetConfig{
		// Rewrites: map[string]string{
		// 	"^/api/(.*)": "/$1",
		// },
		Rewrites:     rewrites,
		ChangeOrigin: cfg.Backend.ChangeOrigin,
	}))
}
