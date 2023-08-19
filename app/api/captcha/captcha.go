package captcha

import (
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/zoox"
)

// New creates a new captcha handler.
func New(cfg *config.Config) func(*zoox.Context) {
	return func(ctx *zoox.Context) {
		service.GenerateCaptcha(cfg, ctx)
	}
}
