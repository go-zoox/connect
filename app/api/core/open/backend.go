package open

import (
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/proxy/utils/rewriter"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) func(*zoox.Context) {
	// @TODO
	if cfg.Upstream.IsValid() {
		cfg.Backend.Protocol = cfg.Upstream.Protocol
		cfg.Backend.Host = cfg.Upstream.Host
		cfg.Backend.Port = cfg.Upstream.Port
	}

	backend := cfg.Backend.String()
	rewrites := rewriter.Rewriters{}
	if !cfg.Backend.IsDisablePrefixRewrite {
		rewrites = rewriter.Rewriters{
			{
				From: fmt.Sprintf("^%s/open/(.*)", cfg.Backend.Prefix),
				To:   "/open/$1",
			},
		}
	}

	return zoox.WrapH(proxy.NewSingleTarget(backend, &proxy.SingleTargetConfig{
		// Rewrites: map[string]string{
		// 	"^/api/open/(.*)": "/open/$1",
		// },
		// Rewrites: rewriter.Rewriters{
		// 	{
		// 		From: "^/api/open/(.*)",
		// 		To:   fmt.Sprintf("%s/open/$1", prefix),
		// 	},
		// },
		Rewrites: rewrites,
		// OnRequest: func(req *http.Request) error {
		// 	fmt.Println("open:", req.URL.Path)
		// 	return nil
		// },
		// OnResponse: func(res *http.Response) error {
		// 	if res.ContentLength == 0 {
		// 		if strings.Contains(res.Request.Header.Get("Accept"), "application/json") {
		// 			res.Header.Set("Content-Type", "application/json")
		// 		}
		// 	}

		// 	return nil
		// },
		ChangeOrigin: cfg.Backend.ChangeOrigin,
	}))
}
