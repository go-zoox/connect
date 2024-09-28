package public

import (
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/zoox"
)

type LoginProviderMetadata struct {
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
}

func GetLoginProvider(c *config.Config) func(ctx *zoox.Context) {
	return func(ctx *zoox.Context) {
		provider := ctx.Param().Get("provider").String()
		if provider == "" {
			ctx.JSON(404, zoox.H{
				"code":    404,
				"message": "provider is missing",
			})
			return
		}

		cfg, err := oauth2.Get(provider)
		if err != nil {
			ctx.JSON(404, zoox.H{
				"code":    404,
				"message": fmt.Sprintf("provider(%s) not found", provider),
			})
			return
		}

		ctx.JSON(200, &LoginProviderMetadata{
			ClientID:    cfg.ClientID,
			RedirectURI: cfg.RedirectURI,
		})
	}
}
