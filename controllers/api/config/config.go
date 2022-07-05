package app

import (
	"github.com/go-zoox/connect/config"
	"github.com/go-zoox/zoox"
)

type Config struct {
}

func New(cfg *config.Config) func(*zoox.Context) {
	return func(ctx *zoox.Context) {
		ctx.JSON(200, zoox.H{
			"code":    200,
			"message": "",
			"result":  nil,
		})
	}
}
