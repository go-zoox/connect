package app

import (
	"fmt"
	"net/http"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/errors"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/zoox"
)

// New ...
func New(cfg *config.Config) func(*zoox.Context) {
	return func(ctx *zoox.Context) {
		token := service.GetToken(ctx)
		if token == "" {
			ctx.Fail(fmt.Errorf("[api.core.app] token is missing"), errors.FailedToGetToken.Code, errors.FailedToGetToken.Message)
			return
		}

		provider := service.GetProvider(ctx)
		if provider == "" {
			ctx.Fail(fmt.Errorf("provider is missing"), errors.FailedToGetOAuth2Provider.Code, errors.FailedToGetOAuth2Provider.Message)
			return
		}

		app, statusCode, err := service.GetApp(ctx, cfg, provider, token)
		if err != nil {
			// @TODO
			if statusCode == http.StatusUnauthorized {
				service.DelToken(ctx)
			}

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
