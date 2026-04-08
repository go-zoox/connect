package api

import (
	"net/http"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/zoox"
)

// Users lists users (stub).
func Users(cfg *config.Config) zoox.HandlerFunc {
	_ = cfg
	return func(ctx *zoox.Context) {
		ctx.JSON(http.StatusOK, zoox.H{
			"code":    200,
			"message": "",
			"result":  []any{},
		})
	}
}
