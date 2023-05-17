package middleware

import (
	"net/http"

	"github.com/go-zoox/jwt"
	"github.com/go-zoox/zoox"
)

// CreateZooxMiddlewareOptions ...
type CreateZooxMiddlewareOptions struct {
	RequireAuth bool
}

// CreateZooxMiddleware ...
func CreateZooxMiddleware(opts ...*CreateZooxMiddlewareOptions) zoox.Middleware {
	var signer jwt.Jwt
	var optsX *CreateZooxMiddlewareOptions
	if len(opts) > 0 && opts[0] != nil {
		optsX = opts[0]
	}

	return func(ctx *zoox.Context) {
		if signer == nil {
			signer = jwt.New(ctx.App.SecretKey)
		}

		token := ctx.Get("x-connect-token")
		if token != "" {
			user := &User{}
			if err := user.Decode(signer, token); err != nil {
				if ctx.AcceptJSON() {
					ctx.JSON(http.StatusUnauthorized, zoox.H{
						"code":    401000,
						"message": err.Error(),
					})
				} else {
					ctx.Status(401)
				}
				return
			}

			ctx.User().Set(user)
		}

		if optsX != nil && optsX.RequireAuth {
			if ctx.User().Get() == nil {
				if ctx.AcceptJSON() {
					ctx.JSON(http.StatusUnauthorized, zoox.H{
						"code":    401001,
						"message": "Unauthorized",
					})
				} else {
					ctx.Status(401)
				}
				return
			}
		}

		ctx.Next()
	}
}
