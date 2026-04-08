package api

import (
	"net/http"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/zoox"
)

// User returns the current admin user placeholder (stub).
func User(cfg *config.Config) zoox.HandlerFunc {
	_ = cfg
	return func(ctx *zoox.Context) {
		ctx.JSON(http.StatusOK, zoox.H{
			"code":    0,
			"message": "",
			"result": zoox.H{
				"id":       "0",
				"username": "admin",
			},
		})
	}
}
