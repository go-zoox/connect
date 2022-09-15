package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-zoox/connect/app/api/captcha"
	"github.com/go-zoox/connect/app/api/favicon"
	"github.com/go-zoox/connect/app/api/page"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/middleware"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/crypto/jwt"
	"github.com/go-zoox/zoox"

	apiApp "github.com/go-zoox/connect/app/api/core/app"
	apiBackend "github.com/go-zoox/connect/app/api/core/backend"
	apiConfig "github.com/go-zoox/connect/app/api/core/config"
	apiMenus "github.com/go-zoox/connect/app/api/core/menus"
	apiUser "github.com/go-zoox/connect/app/api/core/user"

	//
	apiOpen "github.com/go-zoox/connect/app/api/core/open"
)

func New(app *zoox.Application, cfg *config.Config) {
	app.Use(middleware.OAuth2(cfg))
	app.Use(middleware.Auth(cfg))

	// if e.staticHandler != nil {
	// 	e.cfg.Mode = "production"

	// 	cfg.IndexHTML = e.staticHandler(app)
	// } else {
	// 	// indexHTML, _ := StaticFS.ReadFile("web/static/index.html")
	// 	// cfg.IndexHTML = string(indexHTML)

	// 	// staticfs, _ := fs.Sub(StaticFS, "web/static")
	// 	// app.StaticFS("/static/", http.FS(staticfs))

	// 	// dev mode will not use static
	// }

	app.Get("/favicon.ico", favicon.Get(cfg))

	app.Get("/captcha", captcha.New(cfg))

	pg := page.New(cfg)

	// api
	api := app.Group("/api")
	{
		api.Get("/app", apiApp.New(cfg))
		api.Get("/user", apiUser.New(cfg))
		api.Get("/menus", apiMenus.New(cfg))
		api.Get("/users", apiUser.GetUsers(cfg))
		api.Get("/config", apiConfig.New(cfg))
		//
		api.Post("/login", apiUser.Login(cfg))
		//
		api.Get("/page/health", pg.Health(cfg))
		// open
		api.Any("/open/*", apiOpen.New(cfg))
		//
		api.Any(
			"/*",
			func(ctx *zoox.Context) {
				signer := jwt.New(cfg.SecretKey)

				token := service.GetToken(ctx)
				user, err := service.GetUser(cfg, token)
				if err != nil {
					ctx.JSON(http.StatusUnauthorized, err)
					return
				}

				timestamp := time.Now().UnixMilli()
				jwtToken, err := signer.Sign(map[string]interface{}{
					"user_id":       user.ID,
					"user_nickname": user.Nickname,
					"user_avatar":   user.Avatar,
					"user_email":    user.Email,
				})
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, err)
					return
				}

				ctx.Request.Header.Set("X-Connect-Timestamp", fmt.Sprintf("%d", timestamp))
				ctx.Request.Header.Set("X-Connect-Token", jwtToken)

				// request id
				ctx.Request.Header.Set(zoox.RequestIDHeader, ctx.RequestID())

				ctx.Next()
			},
			apiBackend.New(cfg),
		)
	}

	app.Fallback(pg.RenderPage())
}
