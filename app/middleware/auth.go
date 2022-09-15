package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
)

func Auth(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		if isIgnoreAuthoried(ctx.Path) {
			ctx.Next()
			return
		}

		provider := service.GetProvider(ctx)
		token := service.GetToken(ctx)

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

			// why visit /login
			//	=> token_expired/user_expired/app_expired
			//	=> it means you had checked the token + user + app
			//	=> so why check again?
			//	=> delete token before redirect to /login
			// solution:
			//	1. not check, go login
			//	2. delete token before redirect to /login
			//
			// // if user has login, redirect to from
			// @1 check token
			if token != "" {
				// @2 check user
				if user, err := service.GetUser(cfg, token); err == nil && user != nil {
					// @3 check app
					if app, err := service.GetApp(cfg, provider, token); err != nil && app == nil {
						if from != "" {
							ctx.Session().Del("from")
							ctx.Redirect(from)
						} else {
							ctx.Redirect("/")
						}

						return
					}
				}
			}

			// go to redirect login
			// @login_1 oauth2 login
			if cfg.Auth.Mode == "oauth2" {
				ctx.Redirect(fmt.Sprintf("/login/%s", cfg.Auth.Provider))
				return
			}

			// @login_2 local login, fallback to login page render
			ctx.Next()
			return
		} else if ctx.Path == "/logout" {
			from := ctx.Query().Get("from")
			if from != "" {
				ctx.Session().Set("from", from)
			}

			// delete token before
			service.DelToken(ctx)
			ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(from), "visit_logout"))
			return
		}

		// auth mode from oauth2 => local password
		if cfg.Auth.Mode == "password" && service.GetProvider(ctx) != "" {
			service.DelProvider(ctx)
			service.DelToken(ctx)
		}

		// visit real path
		// [visit real path] @1 check token
		if token == "" {
			// @TODO
			// sleep for a while to avoid too many requests
			time.Sleep(time.Second * 1)

			if ctx.AcceptJSON() {
				ctx.Fail(errors.New("api auth failed"), http.StatusUnauthorized, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(ctx.Request.RequestURI), "token_not_found"))
			return
		} else if user, err := service.GetUser(cfg, token); err != nil && user == nil {
			// [visit real path] @2 check user
			// @TODO
			// sleep for a while to avoid too many requests
			time.Sleep(time.Second * 1)

			if ctx.AcceptJSON() {
				ctx.Fail(errors.New("api auth failed"), http.StatusUnauthorized, "Unauthorized", http.StatusUnauthorized)
				return
			}

			logger.Error("[middleware][auth] cannot get user: %v", err)

			// remove token from session to avoid login check visit user again
			service.DelToken(ctx)
			ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(ctx.Request.RequestURI), "user_expired"))
			return
		} else if app, err := service.GetApp(cfg, provider, token); err != nil && app == nil {
			// [visit real path] @2 check app
			// @TODO
			// sleep for a while to avoid too many requests
			time.Sleep(time.Second * 1)

			if ctx.AcceptJSON() {
				ctx.Fail(errors.New("api auth failed"), http.StatusUnauthorized, "Unauthorized", http.StatusUnauthorized)
				return
			}

			logger.Error("[middleware][auth] cannot get app: %v", err)

			// remove token from session to avoid login check visit user + app again
			service.DelToken(ctx)
			ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(ctx.Request.RequestURI), "app_expired"))
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