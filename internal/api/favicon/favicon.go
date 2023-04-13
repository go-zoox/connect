package favicon

import (
	"fmt"
	"io"

	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/connect/internal/errors"
	"github.com/go-zoox/connect/internal/service"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/zoox"
)

func Get(cfg *config.Config) func(*zoox.Context) {
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

		logo := app.Logo
		response, err := fetch.Stream(logo)
		if err != nil {
			ctx.Fail(err, 404, "no favicon found")
			return
		}

		ctx.Set("Content-Type", "image/x-icon")
		ctx.Set("Cache-Control", "public, max-age=31536000")

		if _, err := io.Copy(ctx.Writer, response.Stream); err != nil {
			ctx.Fail(err, 500, "failed to copy favicon")
			return
		}
	}
}
