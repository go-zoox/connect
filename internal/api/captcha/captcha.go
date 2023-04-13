package captcha

import (
	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/connect/internal/service"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) func(*zoox.Context) {
	return func(ctx *zoox.Context) {
		service.GenerateCaptcha(cfg, ctx)
	}
}
