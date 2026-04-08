package api

import (
	"net/http"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/zoox"
)

// MarkerUsers is a stable substring in JSON responses used by router tests.
const MarkerUsers = "__admin_builtin_users__"

// Users lists users (stub).
func Users(cfg *config.Config) zoox.HandlerFunc {
	_ = cfg
	return func(ctx *zoox.Context) {
		ctx.JSON(http.StatusOK, zoox.H{
			"code":    200,
			"message": "",
			"kind":    MarkerUsers,
			"result":  []any{},
		})
	}
}
