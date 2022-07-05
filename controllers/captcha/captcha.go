package captcha

import (
	"github.com/go-zoox/connect/config"
	"github.com/go-zoox/connect/services"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) func(*zoox.Context) {
	return func(ctx *zoox.Context) {
		services.Captcha.Generate(cfg, ctx)
	}
}
