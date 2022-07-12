package user

import (
	"fmt"

	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) func(*zoox.Context) {
	scheme := cfg.Backend.Scheme
	host := cfg.Backend.Host
	port := cfg.Backend.Port

	if scheme == "" {
		scheme = "http"
	}

	if host == "" {
		host = "127.0.0.1"
	}

	if port == 0 {
		port = 8001
	}

	backend := fmt.Sprintf(
		"%s://%s:%d",
		scheme,
		host,
		port,
	)

	return zoox.WrapH(proxy.NewSingleTarget(backend, &proxy.SingleTargetConfig{
		Rewrites: map[string]string{
			"^/api/(.*)": "/$1",
		},
		// OnResponse: func(res *http.Response) error {
		// 	if res.ContentLength == 0 {
		// 		if strings.Contains(res.Request.Header.Get("Accept"), "application/json") {
		// 			res.Header.Set("Content-Type", "application/json")
		// 		}
		// 	}

		// 	return nil
		// },
	}))
}
