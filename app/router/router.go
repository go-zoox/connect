package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-zoox/connect/app/api/captcha"
	"github.com/go-zoox/connect/app/api/favicon"
	"github.com/go-zoox/connect/app/api/page"
	"github.com/go-zoox/connect/app/api/upstream"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/middleware"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/headers"
	"github.com/go-zoox/jwt"
	"github.com/go-zoox/zoox"
	zm "github.com/go-zoox/zoox/middleware"

	apiApp "github.com/go-zoox/connect/app/api/core/app"
	apiBackend "github.com/go-zoox/connect/app/api/core/backend"
	apiConfig "github.com/go-zoox/connect/app/api/core/config"
	apiMenus "github.com/go-zoox/connect/app/api/core/menus"
	apiQRCode "github.com/go-zoox/connect/app/api/core/qrcode"
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

	// api
	api := app.Group("/api", func(group *zoox.RouterGroup) {
		group.Use(zm.CacheControl(&zm.CacheControlConfig{
			Paths:  []string{"^/api/(app|user|menus|users|config)$"},
			MaxAge: 24 * time.Hour,
		}))

		group.Get("/app", apiApp.New(cfg))
		group.Get("/user", apiUser.New(cfg))
		group.Get("/menus", apiMenus.New(cfg))
		group.Get("/users", apiUser.GetUsers(cfg))
		group.Get("/config", apiConfig.New(cfg))
		// qrcode
		group.Get("/qrcode/device/uuid", apiQRCode.GenerateDeviceUUID(cfg))
		group.Get("/qrcode/device/status", apiQRCode.GetDeviceStatus(cfg))
		group.Post("/qrcode/device/token", apiQRCode.GetDeviceToken(cfg))
		group.Get("/qrcode/device/user", apiQRCode.GetUser(cfg))
		//
		group.Post("/login", apiUser.Login(cfg))
	})

	// @TODO
	if cfg.Upstream.Host != "" {
		pg := upstream.New(cfg)
		app.Fallback(func(ctx *zoox.Context) {
			signer := jwt.New(cfg.SecretKey)

			token := service.GetToken(ctx)
			user, err := service.GetUser(ctx, cfg, token)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, err)
				return
			}

			timestamp := time.Now().UnixMilli()
			jwtToken, err := signer.Sign(map[string]interface{}{
				"user_id":             user.ID,
				"user_nickname":       user.Nickname,
				"user_avatar":         user.Avatar,
				"user_email":          user.Email,
				"user_feishu_open_id": user.FeishuOpenID,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}

			ctx.Request.Header.Set("X-Connect-Timestamp", fmt.Sprintf("%d", timestamp))
			ctx.Request.Header.Set("X-Connect-Token", jwtToken)

			// request id
			ctx.Request.Header.Set(headers.XRequestID, ctx.RequestID())

			pg.RenderPage()(ctx)
		})
		return
	}

	// proxy pass
	pg := page.New(cfg)
	// proxy pass => backend
	//
	api.Get("/page/health", pg.Health(cfg))
	// open
	api.Any("/open/*", apiOpen.New(cfg))
	api.Any(
		"/*",
		func(ctx *zoox.Context) {
			signer := jwt.New(cfg.SecretKey)

			token := service.GetToken(ctx)
			user, err := service.GetUser(ctx, cfg, token)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, err)
				return
			}

			timestamp := time.Now().UnixMilli()
			jwtToken, err := signer.Sign(map[string]interface{}{
				"user_id":             user.ID,
				"user_nickname":       user.Nickname,
				"user_avatar":         user.Avatar,
				"user_email":          user.Email,
				"user_feishu_open_id": user.FeishuOpenID,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}

			ctx.Request.Header.Set("X-Connect-Timestamp", fmt.Sprintf("%d", timestamp))
			ctx.Request.Header.Set("X-Connect-Token", jwtToken)

			// request id
			ctx.Request.Header.Set(headers.XRequestID, ctx.RequestID())

			ctx.Next()
		},
		apiBackend.New(cfg),
	)
	// proxy pass => frontend
	app.Fallback(pg.RenderPage())
}
