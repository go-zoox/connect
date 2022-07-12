package app

import (
	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/connect/internal/errors"
	"github.com/go-zoox/connect/internal/service"
	"github.com/go-zoox/zoox"
)

type Config struct {
}

func New(cfg *config.Config) func(*zoox.Context) {
	return func(ctx *zoox.Context) {
		token := service.GetToken(ctx)
		if token == "" {
			ctx.Fail(errors.FailedToGetToken.Code, errors.FailedToGetToken.Message)
			return
		}

		provider := service.GetProvider(ctx)
		if provider == "" {
			ctx.Fail(errors.FailedToGetOAuth2Provider.Code, errors.FailedToGetOAuth2Provider.Message)
			return
		}

		ctx.JSON(200, zoox.H{
			"code":    200,
			"message": "",
			"result":  nil,
		})
	}
}
