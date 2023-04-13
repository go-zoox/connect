package app

import (
	"fmt"

	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/connect/internal/errors"
	"github.com/go-zoox/connect/internal/service"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) func(*zoox.Context) {
	return func(ctx *zoox.Context) {
		token := service.GetToken(ctx)
		if token == "" {
			ctx.Fail(fmt.Errorf("token is missing"), errors.FailedToGetToken.Code, errors.FailedToGetToken.Message)
			return
		}

		provider := service.GetProvider(ctx)
		if provider == "" {
			ctx.Fail(fmt.Errorf("provider is missing"), errors.FailedToGetOAuth2Provider.Code, errors.FailedToGetOAuth2Provider.Message)
			return
		}

		app, err := service.GetApp(ctx, cfg, provider, token)
		if err != nil {
			ctx.Fail(err, errors.FailedToGetApps.Code, errors.FailedToGetApps.Message)
			return
		}

		ctx.JSON(200, zoox.H{
			"code":    200,
			"message": "",
			"result":  app,
		})
	}
}
