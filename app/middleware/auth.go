package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/connect/app/utils"
	"github.com/go-zoox/headers"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
)

// Auth ...
func Auth(cfg *config.Config) zoox.HandlerFunc {
	isIgnoreAuthoriedMatcher := utils.CreateIsPathIgnoreAuthoriedMatcher(func(opt *utils.CreateIsPathIgnoreAuthoriedMatcherOption) {
		opt.Excludes = append([]string{
			// "^/api/login$",
			fmt.Sprintf("^/api%s$", cfg.BuiltInAPIs.Login),
			// "^/api/app$",
			fmt.Sprintf("^/api%s$", cfg.BuiltInAPIs.App),
			// "^/api/qrcode/",
			fmt.Sprintf("^/api%s/", cfg.BuiltInAPIs.QRCode),
			// "^/api/_/",
			fmt.Sprintf("^/api%s/", cfg.BuiltInAPIs.Public),
		}, cfg.Auth.IgnorePaths...)
	})

	return func(ctx *zoox.Context) {
		if !cfg.Auth.IsIgnorePathsDisabled {
			if isIgnoreAuthoriedMatcher.Match(ctx.Path) {
				ctx.State().Set("@@ignore_auth", true)

				ctx.Next()
				return
			}
		}

		if cfg.Auth.IsIgnoreWhenHeaderAuthorizationFound {
			if ctx.Header().Get(headers.Authorization) != "" {
				ctx.Next()
				return
			}
		}

		provider := service.GetProvider(ctx)
		token := service.GetToken(ctx)

		if ctx.Path == "/login" {
			from := ctx.Query().Get("from").String()
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
				if user, _, err := service.GetUser(ctx, cfg, token); err == nil && user != nil {
					// @3 check app
					if app, _, err := service.GetApp(ctx, cfg, provider, service.GetToken(ctx)); err != nil && app == nil {
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
			from := ctx.Query().Get("from").String()
			if from != "" {
				ctx.Session().Set("from", from)
			}

			// delete token before
			service.DelToken(ctx)
			if provider != "" {
				ctx.Redirect(fmt.Sprintf("/logout/%s", provider))
			} else {
				ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(from), "visit_logout"))
			}
			return
		} else if ctx.Path == "/register" {
			from := ctx.Query().Get("from").String()
			if from != "" {
				ctx.Session().Set("from", from)
			}
			// go to redirect register
			// @register_1 oauth2 register
			if cfg.Auth.Mode == "oauth2" {
				ctx.Redirect(fmt.Sprintf("/register/%s?%s", cfg.Auth.Provider, ctx.Request.URL.RawQuery))
				return
			}

			// @register_2 local register, fallback to register page render
			ctx.Next()
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
				ctx.Fail(errors.New("[middlware.auth] token is missing"), http.StatusUnauthorized, "token is missing", http.StatusUnauthorized)
				return
			}

			ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(ctx.Request.RequestURI), "token_not_found"))
			return
		} else if user, _, err := service.GetUser(ctx, cfg, token); err != nil && user == nil {
			// [visit real path] @2 check user
			// @TODO
			// sleep for a while to avoid too many requests
			time.Sleep(time.Second * 1)

			if ctx.AcceptJSON() {
				ctx.Fail(fmt.Errorf("failed to get user (err: %s)", err), http.StatusUnauthorized, fmt.Sprintf("failed to get user (err: %s)", err), http.StatusUnauthorized)
				return
			}

			logger.Error("[middleware][auth] cannot get user: %v", err)

			// remove token from session to avoid login check visit user again
			service.DelToken(ctx)
			ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(ctx.Request.RequestURI), "user_expired"))
			return
		} else if app, _, err := service.GetApp(ctx, cfg, provider, service.GetToken(ctx)); err != nil || app == nil {
			// [visit real path] @2 check app
			// @TODO
			// sleep for a while to avoid too many requests
			time.Sleep(time.Second * 1)

			if ctx.AcceptJSON() {
				ctx.Fail(fmt.Errorf("failed to get app (err: %s)", err), http.StatusUnauthorized, fmt.Sprintf("failed to get app (err: %s)", err), http.StatusUnauthorized)
				return
			}

			logger.Error("[middleware][auth] cannot get app: %v", err)

			// remove token from session to avoid login check visit user + app again
			service.DelToken(ctx)
			ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(ctx.Request.RequestURI), "app_expired"))
			return
		}

		if len(cfg.Auth.AllowUsernames) > 0 {
			user, status, err := service.GetUser(ctx, cfg, token)
			if err != nil {
				ctx.Fail(err, status, fmt.Sprintf("failed to get user when validate allow usernames: %s", err))
				return
			}

			isAllowed := false
			for _, username := range cfg.Auth.AllowUsernames {
				if user.Username == username {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				ctx.Fail(fmt.Errorf("username(%s) is not allowed", user.Username), http.StatusForbidden, "Forbidden")
				return
			}
		}

		ctx.Next()
	}
}
