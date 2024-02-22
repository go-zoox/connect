package middleware

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/oauth2"
	oc "github.com/go-zoox/oauth2/create"
	"github.com/go-zoox/random"
	"github.com/go-zoox/zoox"
)

// OAuth2 ...
func OAuth2(cfg *config.Config) zoox.HandlerFunc {
	loginRegExp := regexp.MustCompile("^/login/([^/]+)$")
	logoutRegExp := regexp.MustCompile("^/logout/([^/]+)$")
	registerRegExp := regexp.MustCompile("^/register/([^/]+)$")
	loginCallbackRegExp := regexp.MustCompile("^/login/([^/]+)/callback$")

	return func(ctx *zoox.Context) {
		// login => /login/:provider
		if loginRegExp.MatchString(ctx.Path) {
			provider := loginRegExp.FindStringSubmatch(ctx.Path)[1]
			if clientCfg, err := oauth2.Get(provider); err != nil {
				panic(fmt.Errorf("failed to get oauth2 client config with provider(%s): %s", provider, err))
			} else {
				client, err := oc.Create(clientCfg.Name, clientCfg)
				if err != nil {
					panic(fmt.Errorf("failed to create oauth2 client with provider(%s): %s", provider, err))
				}

				service.SetProvider(ctx, cfg, provider)
				state := random.String(8)
				ctx.Session().Set("oauth2_state", state)

				logger.Infof("[oauth2:start] provider(%s) - state(%s) - from(%s)", provider, state, ctx.Session().Get("from"))

				client.Authorize(state, func(loginURL string) {
					logger.Infof("[oauth2:authorize] url: %s", loginURL)

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

			logger.Infof("[oauth2:callback] provider(%s) - code(%s) - state(%s)", provider, code, state)

			if ctx.Session().Get("oauth2_state") != state {
				logger.Infof("state not match: expect %s, but got %s", ctx.Session().Get("oauth2_state"), state)

				// panic("oauth2_state is not match")
				ctx.Redirect(fmt.Sprintf("/login/%s", provider))
				return
			}

			if clientCfg, err := oauth2.Get(provider); err != nil {
				panic(fmt.Errorf("failed to get oauth2 client config with provider(%s): %s", provider, err))
			} else {
				client, err := oc.Create(clientCfg.Name, clientCfg)
				if err != nil {
					panic(fmt.Errorf("failed to create oauth2 client with provider(%s): %s", provider, err))
				}

				service.SetProvider(ctx, cfg, provider)

				client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
					if err != nil {
						panic(fmt.Errorf("failed to run oauth2 callback with provider(%s): %s", provider, err))
					}

					service.SetToken(ctx, cfg, token.AccessToken)
					service.SetRefreshToken(ctx, cfg, token.RefreshToken)

					from := ctx.Session().Get("from")
					logger.Infof("[oauth2:done] from %s", from)
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

		// register => /register/:provider
		if registerRegExp.MatchString(ctx.Path) {
			provider := registerRegExp.FindStringSubmatch(ctx.Path)[1]
			if clientCfg, err := oauth2.Get(provider); err != nil {
				panic(err)
			} else {
				client, err := oc.Create(clientCfg.Name, clientCfg)
				if err != nil {
					panic(err)
				}

				client.Register(func(registerURL string) {
					ur, err := url.Parse(registerURL)
					if err != nil {
						ctx.Redirect(fmt.Sprintf("/error?code=%d&message=%s", 500, url.QueryEscape("register url parse error")))
					}

					query := ur.Query()
					for k, v := range ctx.Request.URL.Query() {
						query.Add(k, v[0])
					}
					ur.RawQuery = query.Encode()

					ctx.Redirect(ur.String())
				})
			}
			return
		}

		ctx.Next()
	}
}
