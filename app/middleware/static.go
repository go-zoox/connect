package middleware

import (
	"time"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/zoox"
	zw "github.com/go-zoox/zoox/middleware"
)

// DefaultStaticFileMaxAge ...
const DefaultStaticFileMaxAge = 7 * 24 * 60 * 60

// Static ...
func Static(cfg *config.Config) zoox.HandlerFunc {
	var staticFileMaxAge int64 = DefaultStaticFileMaxAge
	if cfg.StaticCacheControlMaxAge != 0 {
		staticFileMaxAge = cfg.StaticCacheControlMaxAge
	}

	return zw.StaticCache(&zw.StaticCacheConfig{
		MaxAge: time.Duration(staticFileMaxAge) * time.Second,
	})
}
