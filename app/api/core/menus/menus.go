package menus

import (
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/errors"
	"github.com/go-zoox/connect/app/service"
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

		menus, statusCode, err := service.GetMenu(ctx, cfg, provider, token)
		if err != nil {
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
