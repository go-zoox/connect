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
	"github.com/go-zoox/proxy/utils/rewriter"
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

// New ...
func New(app *zoox.Application, cfg *config.Config) {
	// manifest
	app.Get("/api/_/built_in_apis", func(ctx *zoox.Context) {
		ctx.JSON(http.StatusOK, cfg.BuiltInAPIs)
	})

	app.Use(middleware.Static(cfg))

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
	app.Group("/api", func(group *zoox.RouterGroup) {
		group.Use(zm.CacheControl(&zm.CacheControlConfig{
			Paths:  []string{"^/api/(app|menus|users|config)$"},
			MaxAge: 30 * time.Second,
		}))

		// /app
		group.Get(cfg.BuiltInAPIs.App, apiApp.New(cfg))
		// /user
		group.Get(cfg.BuiltInAPIs.User, apiUser.New(cfg))
		// /menus
		group.Get(cfg.BuiltInAPIs.Menus, apiMenus.New(cfg))
		// /users
		group.Get(cfg.BuiltInAPIs.Users, apiUser.GetUsers(cfg))
		// /config
		group.Get(cfg.BuiltInAPIs.Config, apiConfig.New(cfg))
		// /qrcode
		qrcodeBasePath := cfg.BuiltInAPIs.QRCode
		group.Get(fmt.Sprintf("%s/device/uuid", qrcodeBasePath), apiQRCode.GenerateDeviceUUID(cfg))
		group.Get(fmt.Sprintf("%s/device/status", qrcodeBasePath), apiQRCode.GetDeviceStatus(cfg))
		group.Post(fmt.Sprintf("%s/device/token", qrcodeBasePath), apiQRCode.GetDeviceToken(cfg))
		group.Get(fmt.Sprintf("%s/device/user", qrcodeBasePath), apiQRCode.GetUser(cfg))
		// /login
		group.Post(cfg.BuiltInAPIs.Login, apiUser.Login(cfg))
	})

	// backend api
	api := app.Group(cfg.Backend.Prefix)

	// open
	api.Any("/open/*", apiOpen.New(cfg))

	// routes
	for _, route := range cfg.Routes {
		app.Logger.Infof("[router] load route: %s => %s (secret_key: %s)", route.Path, route.Backend.String(), route.Backend.SecretKey)

		app.Proxy(route.Path, route.Backend.String(), func(cfgX *zoox.ProxyConfig) {
			if !route.Backend.DisableRewrite {
				cfgX.Rewrites = rewriter.Rewriters{
					{
						From: fmt.Sprintf("%s/(.*)", route.Path),
						To:   "/$1",
					},
				}
			}

			cfgX.OnRequestWithContext = func(ctx *zoox.Context) error {
				// app.Logger.Infof("[api][ms] %s => %s (secret_key: %s)", route.Path, route.Backend.String(), route.Backend.SecretKey)

				if route.Backend.SecretKey != "" {
					signer := jwt.New(route.Backend.SecretKey)

					token := service.GetToken(ctx)
					userIns, _, err := service.GetUser(ctx, cfg, token)
					if err != nil {
						return fmt.Errorf("failed to get user: %w", err)
					}
					timestamp := time.Now().UnixMilli()
					jwtToken, err := userIns.Encode(signer)
					if err != nil {
						return fmt.Errorf("failed to sign jwt token: %w", err)
					}

					ctx.Request.Header.Set("X-Connect-Timestamp", fmt.Sprintf("%d", timestamp))
					ctx.Request.Header.Set("X-Connect-Token", jwtToken)
				}

				// request id
				ctx.Request.Header.Set(headers.XRequestID, ctx.RequestID())
				return nil
			}

			cfgX.OnResponse = func(res *http.Response) error {
				res.Header.Set(headers.XPoweredBy, "go-zoox")
				return nil
			}
		})
	}

	// @TODO
	if cfg.Upstream.IsValid() {
		app.Logger.Infof("mode: upstream")
		app.Logger.Infof("upstream: %s", cfg.Upstream.String())

		up := upstream.New(cfg)
		app.Fallback(func(ctx *zoox.Context) {
			signer := jwt.New(cfg.SecretKey)

			token := service.GetToken(ctx)
			userIns, _, err := service.GetUser(ctx, cfg, token)
			if err != nil {
				// ctx.Logger.Errorf(err)
				fmt.Println("failed to get user:", err)
				ctx.Fail(err, 401002, "user not found", http.StatusUnauthorized)
				return
			}

			timestamp := time.Now().UnixMilli()
			jwtToken, err := userIns.Encode(signer)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}

			ctx.Request.Header.Set("X-Connect-Timestamp", fmt.Sprintf("%d", timestamp))
			ctx.Request.Header.Set("X-Connect-Token", jwtToken)

			// request id
			ctx.Request.Header.Set(headers.XRequestID, ctx.RequestID())

			up.Handle(ctx)
		})
		return
	}

	app.Logger.Infof("mode: frontend + backend")
	app.Logger.Infof("frontend: %s", cfg.Frontend.String())
	app.Logger.Infof("backend: %s", cfg.Backend.String())

	// proxy pass
	pg := page.New(cfg)
	// proxy pass => backend
	//
	api.Get("/page/health", pg.Health(cfg))
	api.Any(
		"/*",
		func(ctx *zoox.Context) {
			signer := jwt.New(cfg.SecretKey)

			token := service.GetToken(ctx)
			userIns, _, err := service.GetUser(ctx, cfg, token)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, err)
				return
			}

			timestamp := time.Now().UnixMilli()
			jwtToken, err := userIns.Encode(signer)
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
