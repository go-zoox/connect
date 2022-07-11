package auth

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-zoox/connect/config"
	"github.com/go-zoox/connect/services"
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
		token := services.Token.Get(ctx)

		// 0. check the path need to be authenticated
		if isIgnoreAuthoried(ctx.Path) {
			ctx.Next()
			return
		}

		if ctx.Path == "/logout" {
			from := ctx.Query("from")

			if token != "" {
				services.Token.Clear(ctx)
			}

			if cfg.Auth.Mode == "oauth2" {
				client.Logout(func(logoutURL string) {
					ctx.Redirect(logoutURL)
				})
				return
			}

			ctx.Redirect(fmt.Sprintf("/login?from=%s", from))
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
		}

		// 1. check the user is authenticated
		// 1.1 authenticated
		if user, err := services.User.Get(cfg, token); err == nil && user != nil {
			if token != "" {
				if ctx.Path == "/login" {
					from := ctx.Query("from")
					if from == "" {
						from = ctx.Session.Get("from")
						ctx.Session.Del("from")
					}

					if from == "" {
						from = "/"
					}

					ctx.Redirect(from)
					return
				}

				from := ctx.Session.Get("from")
				if from != "" {
					ctx.Session.Del("from")
					ctx.Redirect(from)
					return
				}

				ctx.Next()
				return
			}
		} else {
			services.Token.Clear(ctx)
			token = ""
		}

		// 1.2 not authenticated => go to login page
		if ctx.Path != "/login" {
			if ctx.AcceptJSON() {
				ctx.Status(http.StatusUnauthorized)
				return
			}

			if cfg.Auth.Mode == "password" {
				ctx.Redirect(fmt.Sprintf("/login?from=%s", ctx.Request.RequestURI))
			} else if cfg.Auth.Mode == "oauth2" {
				from := ctx.Query("from")
				if from == "" {
					from = ctx.Path
				}

				if from == "" {
					from = "/"
				}

				// save the from url
				ctx.Session.Set("from", from)

				client.Authorize("ops", func(loginUrl string) {
					ctx.Redirect(loginUrl)
				})
			}

			return
		} else {
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

			ctx.Next()
			return
		}
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
