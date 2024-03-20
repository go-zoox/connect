package menus

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
			ctx.Fail(fmt.Errorf("[api.core.menus] token is missing"), errors.FailedToGetToken.Code, errors.FailedToGetToken.Message)
			return
		}

		provider := service.GetProvider(ctx)
		if provider == "" {
			ctx.Fail(fmt.Errorf("provider is missing"), errors.FailedToGetOAuth2Provider.Code, errors.FailedToGetOAuth2Provider.Message)
			return
		}

		menus, statusCode, err := service.GetMenu(ctx, cfg, provider, token)
		if err != nil {
			// @TODO
			if statusCode == http.StatusUnauthorized {
				service.DelToken(ctx)
			}

			ctx.Fail(err, errors.FailedToGetMenus.Code, errors.FailedToGetMenus.Message, statusCode)
			return
		}

		ctx.JSON(200, zoox.H{
			"code":    200,
			"message": "",
			"result":  menus,
		})
	}
}
