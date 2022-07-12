package menus

import (
	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/connect/internal/errors"
	"github.com/go-zoox/connect/internal/service"
	"github.com/go-zoox/zoox"
)

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

		menus, err := service.GetMenu(cfg, provider, token)
		if err != nil {
			ctx.Fail(errors.FailedToGetMenus.Code, errors.FailedToGetMenus.Message+": "+err.Error())
			return
		}

		ctx.JSON(200, zoox.H{
			"code":    200,
			"message": "",
			"result":  menus,
		})
	}
}
