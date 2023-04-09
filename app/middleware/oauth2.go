package middleware

import (
	"fmt"
	"regexp"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/oauth2"
	oc "github.com/go-zoox/oauth2/create"
	"github.com/go-zoox/random"
	"github.com/go-zoox/zoox"
)

func OAuth2(cfg *config.Config) zoox.HandlerFunc {
	loginRegExp := regexp.MustCompile("^/login/([^/]+)$")
	logoutRegExp := regexp.MustCompile("^/logout/([^/]+)$")
	loginCallbackRegExp := regexp.MustCompile("^/login/([^/]+)/callback$")

	return func(ctx *zoox.Context) {
		// login => /login/:provider
		if loginRegExp.MatchString(ctx.Path) {
			provider := loginRegExp.FindStringSubmatch(ctx.Path)[1]
			if clientCfg, err := oauth2.Get(provider); err != nil {
				panic(err)
			} else {
				client, err := oc.Create(clientCfg.Name, clientCfg)
				if err != nil {
					panic(err)
				}

				service.SetProvider(ctx, cfg, provider)
				state := random.String(8)
				ctx.Session().Set("oauth2_state", state)

				client.Authorize(state, func(loginURL string) {
					ctx.Redirect(loginURL)
				})
			}
			return
		}

		// login callback => /login/:provider/callback
		if loginCallbackRegExp.MatchString(ctx.Path) {
			code := ctx.Query().Get("code").String()
			state := ctx.Query().Get("state").String()
			provider := loginCallbackRegExp.FindStringSubmatch(ctx.Path)[1]

			if ctx.Session().Get("oauth2_state") != state {
				logger.Infof("state not match: expect %s, but got %s", ctx.Session().Get("oauth2_state"), state)

				// panic("oauth2_state is not match")
				ctx.Redirect(fmt.Sprintf("/login/%s", provider))
				return
			}

			if clientCfg, err := oauth2.Get(provider); err != nil {
				panic(err)
			} else {
				client, err := oc.Create(clientCfg.Name, clientCfg)
				if err != nil {
					panic(err)
				}

				service.SetProvider(ctx, cfg, provider)

				client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
					if err != nil {
						panic(err)
					}

					service.SetToken(ctx, cfg, token.AccessToken)

					from := ctx.Session().Get("from")
					if from != "" {
						ctx.Session().Del("from")
						ctx.Redirect(from)
					} else {
						ctx.Redirect("/")
					}
				})
			}
			return
		}

		// logout => /logout/:provider
		if logoutRegExp.MatchString(ctx.Path) {
			provider := service.GetProvider(ctx)
			// cannot get provider mean not oauth2
			if provider == "" {
				ctx.Next()
				return
			}

			clientCfg, err := oauth2.Get(provider)
			if err != nil {
				panic(err)
			}

			client, err := oc.Create(clientCfg.Name, clientCfg)
			if err != nil {
				panic(err)
			}
			client.Logout(func(logoutURL string) {
				// clear token
				service.DelToken(ctx)
				// clear provider
				ctx.Session().Del("provider")

				from := ctx.Query().Get("from").String()
				if from != "" {
					ctx.Session().Set("from", from)
				}

				ctx.Redirect(logoutURL)
			})
			return
		}

		ctx.Next()
	}
}
