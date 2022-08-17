package api

// import (
// 	"github.com/go-zoox/zoox"

// 	apiMenus "github.com/go-zoox/connect/controllers/api/menus"
// 	apiUser "github.com/go-zoox/connect/controllers/api/user"

// 	apiApp "github.com/go-zoox/connect/controllers/api/app"

// 	apiConfig "github.com/go-zoox/connect/controllers/api/config"
// )

// func init() {
// 	zoox.DefaultGroup("/api", func(r *zoox.RouterGroup) {
// 		r.Get("/app", apiApp.New())
// 		r.Get("/user", apiUser.New())
// 		r.Get("/menus", apiMenus.New())
// 		r.Get("/config", apiConfig.New())

// 		//
// 		r.Post("/login", apiUser.Login())

// 		// r.Any("/api/*", apiBackend.New())
// 	})
// }
