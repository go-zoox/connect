package oauth2

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-zoox/connect/config"
	gooauth2 "github.com/go-zoox/oauth2"
	goaDoreamon "github.com/go-zoox/oauth2/doreamon"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) zoox.HandlerFunc {
	client, err := goaDoreamon.New(cfg.OAuth2.ClientID, cfg.OAuth2.ClientSecret, cfg.OAuth2.RedirectURI)
	if err != nil {
		panic(err)
	}

	sessionId := "gz_uid"
	tokenId := "gz_ut"

	return func(ctx *zoox.Context) {
		if ctx.Path == "/logout" {
			if ctx.Cookie.Get(sessionId) != "" {
				ctx.Cookie.Set(sessionId, "", 0)
			}

			client.Logout(func(logoutUrl string) {
				ctx.Redirect(logoutUrl)
			})
			return
		}

		if ctx.Path == "/login" {
			if ctx.Cookie.Get(sessionId) != "" {
				ctx.Redirect("/")
			}

			client.Authorize("ops", func(loginUrl string) {
				ctx.Redirect(loginUrl)
			})

			// ctx.Next()
			return
		}

		if ctx.Path == "/login/callback" {
			code := ctx.Query("code")
			state := ctx.Query("state")

			client.Callback(code, state, func(user *gooauth2.User, token *gooauth2.Token, err error) {
				if err != nil {
					panic(err)
				}

				ctx.Cookie.Set(sessionId, user.ID, 2*time.Hour)
				ctx.Cookie.Set(tokenId, token.AccessToken, 2*time.Hour)

				ctx.Redirect("/")
			})

			return
		}

		if ctx.Cookie.Get(sessionId) == "" {
			excludes := []string{
				"^/__umi_ping$",
				"^/robots.txt$",
				"^/sockjs-node",
				"\\.(css|js|ico|jpg|png|jpeg|webp|gif|socket|ws)$",
			}
			for _, exclude := range excludes {
				matched, err := regexp.MatchString(exclude, ctx.Path)
				if err == nil && matched {
					ctx.Next()
					return
				} else if err != nil {
					panic(err)
				}
			}

			ctx.Redirect("/login")
			return
		}

		ctx.Next()

		fmt.Println("xxx:", ctx.StatusCode, ctx.Writer.Status())
	}
}
