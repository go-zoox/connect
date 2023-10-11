package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-zoox/jwt"
)

// CreateGinMiddlewareOptions ...
type CreateGinMiddlewareOptions struct {
	RequireAuth bool
}

const ContextUserKeyForGin = "zoox.connect::user"

// CreateGinMiddleware ...
func CreateGinMiddleware(secretKey string, opts ...*CreateGinMiddlewareOptions) gin.HandlerFunc {
	var signer jwt.Jwt
	var optsX *CreateGinMiddlewareOptions
	if len(opts) > 0 && opts[0] != nil {
		optsX = opts[0]
	}

	return func(ctx *gin.Context) {
		if signer == nil {
			signer = jwt.New(secretKey)
		}

		token := ctx.GetHeader("x-connect-token")
		if token != "" {
			user := &User{}
			if err := user.Decode(signer, token); err != nil {
				// if ctx.AcceptJSON() {
				// 	ctx.JSON(http.StatusUnauthorized, gin.H{
				// 		"code":    401001,
				// 		"message": "Unauthorized",
				// 	})
				// } else {
				// 	ctx.Status(401)
				// }

				ctx.Status(401)
				return
			}

			ctx.Set(ContextUserKeyForGin, user)
		}

		if optsX != nil && optsX.RequireAuth {
			if _, ok := ctx.Get(ContextUserKeyForGin); !ok {
				// if ctx.AcceptJSON() {
				// 	ctx.JSON(http.StatusUnauthorized, gin.H{
				// 		"code":    401001,
				// 		"message": "Unauthorized",
				// 	})
				// } else {
				// 	ctx.Status(401)
				// }

				ctx.Status(401)
				return
			}
		}

		ctx.Next()
	}
}
