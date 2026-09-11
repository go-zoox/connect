package api

import (
	"net/http"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/zoox"
)

// MarkerRoles is a stable substring in JSON responses used by router tests.
const MarkerRoles = "__admin_builtin_roles__"

// Roles returns roles (stub; wired for route registration tests).
func Roles(cfg *config.Config) zoox.HandlerFunc {
	_ = cfg
	return func(ctx *zoox.Context) {
		ctx.JSON(http.StatusOK, zoox.H{
			"code":    200,
			"message": "",
			"result": zoox.H{
				"kind": MarkerRoles,
			},
		})
	}
}
