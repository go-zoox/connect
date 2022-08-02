package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/connect/internal/service"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
)

func Auth(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		if isIgnoreAuthoried(ctx.Path) {
			ctx.Next()
			return
		}

		if ctx.Path == "/login" {
			from := ctx.Query().Get("from")
			if from != "" {
				ctx.Session().Set("from", from)
			}

			// auth mode from oauth2 => local password
			if cfg.Auth.Mode == "password" && service.GetProvider(ctx) != "" {
				service.DelProvider(ctx)
				service.DelToken(ctx)
			}

			// if user has login, redirect to from
			token := service.GetToken(ctx)
			if token != "" {
				if user, err := service.GetUser(cfg, token); err == nil && user != nil {
					if from != "" {
						ctx.Session().Del("from")
						ctx.Redirect(from)
					} else {
						ctx.Redirect("/")
					}

					return
				}
			}

			// oauth2 login
			if cfg.Auth.Mode == "oauth2" {
				ctx.Redirect(fmt.Sprintf("/login/%s", cfg.Auth.Provider))
				return
			}

			// local login, fallback to login page render
			ctx.Next()
			return
		} else if ctx.Path == "/logout" {
			service.DelToken(ctx)

			from := ctx.Query().Get("from")
			if from != "" {
				ctx.Session().Set("from", from)
			}

			ctx.Redirect(fmt.Sprintf("/login?from=%s", from))
			return
		}

		if ctx.AcceptJSON() {
			ctx.Fail(errors.New("api auth failed"), http.StatusUnauthorized, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// auth mode from oauth2 => local password
		if cfg.Auth.Mode == "password" && service.GetProvider(ctx) != "" {
			service.DelProvider(ctx)
			service.DelToken(ctx)
		}

		provider := service.GetProvider(ctx)

		token := service.GetToken(ctx)
		if token == "" {
			ctx.Redirect("/login?from=" + url.QueryEscape(ctx.Request.RequestURI))
			return
		} else if user, err := service.GetUser(cfg, token); err != nil && user == nil {
			logger.Error("[middleware][auth] cannot get user: %v", err)

			time.Sleep(3 * time.Second)

			ctx.Redirect("/login?from=" + url.QueryEscape(ctx.Request.RequestURI))
			return
		} else if app, err := service.GetApp(cfg, provider, token); err != nil && app == nil {
			logger.Error("[middleware][auth] cannot get app: %v", err)

			time.Sleep(3 * time.Second)

			ctx.Redirect("/login?from=" + url.QueryEscape(ctx.Request.RequestURI))
			return
		}

		ctx.Next()
	}
}

func isIgnoreAuthoried(path string) bool {
	excludes := []string{
		"^/api/login$",
		"^/api/app$",
		"^/captcha$",
		"^/favicon.ico",
		"^/__umi_ping$",
		"^/__umiDev/routes$",
		"^/robots.txt$",
		"^/sockjs-node",
		"\\.(css|js|ico|jpg|png|jpeg|webp|gif|socket|ws|map)$",
		"\\.hot-update.json$",
		"^/api/open/",
	}
	for _, exclude := range excludes {
		matched, err := regexp.MatchString(exclude, path)
		if err == nil && matched {
			return true
		}
	}

	return false
}
