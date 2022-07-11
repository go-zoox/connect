package connect

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/go-zoox/connect/config"
	captcha "github.com/go-zoox/connect/controllers/captcha"
	"github.com/go-zoox/connect/middlewares/auth"
	"github.com/go-zoox/connect/services"
	jwtsigner "github.com/go-zoox/jwt"

	apiApp "github.com/go-zoox/connect/controllers/api/app"
	apiBackend "github.com/go-zoox/connect/controllers/api/backend"
	apiConfig "github.com/go-zoox/connect/controllers/api/config"
	apiMenus "github.com/go-zoox/connect/controllers/api/menus"
	apiUser "github.com/go-zoox/connect/controllers/api/user"

	"github.com/go-zoox/connect/controllers/page"
	"github.com/go-zoox/zoox"
	z "github.com/go-zoox/zoox/default"
)

//go:embed public
var StaticFS embed.FS

type Connect struct {
	core *zoox.Application
	cfg  *config.Config
	//
	apiHandler    func(*zoox.Application)
	staticHandler func(*zoox.Application) string
}

func New() *Connect {
	app := z.Default()

	return &Connect{
		core: app,
	}
}

func (e *Connect) RegisterApi(fn func(*zoox.Application)) {
	e.apiHandler = fn
}

func (e *Connect) RegisterStatic(fn func(*zoox.Application) string) {
	e.staticHandler = fn
}

func (e *Connect) handle(cfg *config.Config) {
	// env
	// if cfg.Mode != "" {
	// 	e.config.Mode = cfg.Mode
	// } else if os.Getenv("MODE") != "" {
	// 	e.config.Mode = os.Getenv("MODE")
	// }

	e.cfg = cfg

	// @TODO
	loadingHTML, _ := StaticFS.ReadFile("public/loading.html")
	cfg.LoadingHTML = string(loadingHTML)

	e.core.SecretKey = cfg.SecretKey

	e.core.Use(auth.New(cfg))

	// if e.staticHandler != nil {
	// 	e.cfg.Mode = "production"

	// 	cfg.IndexHTML = e.staticHandler(e.core)
	// } else {
	// 	// indexHTML, _ := StaticFS.ReadFile("public/static/index.html")
	// 	// cfg.IndexHTML = string(indexHTML)

	// 	// staticfs, _ := fs.Sub(StaticFS, "public/static")
	// 	// e.core.StaticFS("/static/", http.FS(staticfs))

	// 	// dev mode will not use static
	// }

	e.core.Get("/captcha", captcha.New(cfg))

	pg := page.New(cfg)

	jwt := jwtsigner.NewHS256(cfg.SecretKey)

	// api
	api := e.core.Group("/api")
	{
		api.Get("/app", apiApp.New(cfg))
		api.Get("/user", apiUser.New(cfg))
		api.Get("/menus", apiMenus.New(cfg))
		api.Get("/config", apiConfig.New(cfg))
		//
		api.Post("/login", apiUser.Login(cfg))
		//
		api.Get("/page/health", pg.Health(cfg))
		//
		if e.apiHandler == nil {
			api.Any(
				"/*",
				func(ctx *zoox.Context) {
					token := services.Token.Get(ctx)
					user, err := services.User.Get(cfg, token)
					if err != nil {
						ctx.JSON(http.StatusUnauthorized, err)
						return
					}

					timestamp := time.Now().UnixMilli()
					jwt.Set("user", map[string]string{
						"id":       user.ID,
						"nickname": user.Nickname,
						"avatar":   user.Avatar,
						"email":    user.Email,
					})
					jwtToken, err := jwt.Sign()
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, err)
						return
					}

					ctx.Request.Header.Set("X-Connect-Timestamp", fmt.Sprintf("%d", timestamp))
					ctx.Request.Header.Set("X-Connect-Token", jwtToken)

					ctx.Next()
				},
				apiBackend.New(cfg),
			)
		} else {
			e.apiHandler(e.core)
		}
	}

	e.core.Fallback(pg.RenderPage())
}

func (e *Connect) Start(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	e.handle(cfg)
	e.core.Run(fmt.Sprintf(":%d", e.cfg.Port))

	return nil
}
