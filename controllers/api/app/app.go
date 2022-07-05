package app

import (
	"github.com/go-zoox/connect/config"
	"github.com/go-zoox/connect/services"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) func(*zoox.Context) {
	return func(ctx *zoox.Context) {
		token := services.Token.Get(ctx)
		app, err := services.App.Get(cfg, token)
		if err != nil {
			panic(err)
		}

		ctx.JSON(200, zoox.H{
			"code":    200,
			"message": "",
			"result":  app,
		})
	}
}
