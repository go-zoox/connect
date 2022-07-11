package auth

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-zoox/connect/config"
	"github.com/go-zoox/connect/services"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/oauth2"
	goaDoreamon "github.com/go-zoox/oauth2/doreamon"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) zoox.HandlerFunc {
	var client *oauth2.Client
	var err error
	if cfg.Auth.Mode == "oauth2" {
		client, err = goaDoreamon.New(&goaDoreamon.DoreamonConfig{
			ClientID:     cfg.Auth.OAuth2.ClientID,
			ClientSecret: cfg.Auth.OAuth2.ClientSecret,
			RedirectURI:  cfg.Auth.OAuth2.RedirectURI,
			Scope:        cfg.Auth.OAuth2.Scope,
		})
		if err != nil {
			panic(err)
		}
	}

	return func(ctx *zoox.Context) {
		if isIgnoreAuthoried(ctx.Path) {
			ctx.Next()
			return
		}

		if ctx.Path == "/login" {
			from := ctx.Query("from")
			if from != "" {
				ctx.Session.Set("from", from)
			}

			if cfg.Auth.Mode == "oauth2" {
				client.Authorize("ops", func(loginURL string) {
					ctx.Redirect(loginURL)
				})
				return
			}

			ctx.Redirect("/")
			return
		} else if ctx.Path == "/login/callback" {
			code := ctx.Query("code")
			state := ctx.Query("state")

			client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
				if err != nil {
					panic(err)
				}

				services.Token.Set(ctx, token.AccessToken)

				from := ctx.Session.Get("from")
				if from != "" {
					ctx.Session.Del("from")
					ctx.Redirect(from)
				} else {
					ctx.Redirect("/")
				}
			})

			return
		} else if ctx.Path == "/logout" {
			from := ctx.Query("from")

			services.Token.Clear(ctx)

			if cfg.Auth.Mode == "oauth2" {
				client.Logout(func(logoutURL string) {
					ctx.Redirect(logoutURL)
				})
				return
			}

			ctx.Redirect(fmt.Sprintf("/login?from=%s", from))
			return
		}

		token := services.Token.Get(ctx)
		if user, err := services.User.Get(cfg, token); err != nil && user == nil {
			logger.Error("[middleware][auth] cannot get user: ", err)

			time.Sleep(3 * time.Second)

			ctx.Redirect("/login")
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
		"^/__umi_ping$",
		"^/robots.txt$",
		"^/sockjs-node",
		"\\.(css|js|ico|jpg|png|jpeg|webp|gif|socket|ws|map)$",
		"\\.hot-update.json$",
	}
	for _, exclude := range excludes {
		matched, err := regexp.MatchString(exclude, path)
		if err == nil && matched {
			return true
		}
	}

	return false
}
