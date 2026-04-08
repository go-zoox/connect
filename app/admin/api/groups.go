package api

import (
	"net/http"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/zoox"
)

// MarkerGroups is a stable substring in JSON responses used by router tests.
const MarkerGroups = "__admin_builtin_groups__"

// Groups returns groups (stub; wired for route registration tests).
func Groups(cfg *config.Config) zoox.HandlerFunc {
	_ = cfg
	return func(ctx *zoox.Context) {
		ctx.JSON(http.StatusOK, zoox.H{
			"code":    200,
			"message": "",
			"result": zoox.H{
				"kind": MarkerGroups,
			},
		})
	}
}
