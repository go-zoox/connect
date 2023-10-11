package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-zoox/jwt"
)

// CreateGinMiddlewareOptions ...
type CreateGinMiddlewareOptions struct {
	SecretKey string
	//
	RequireAuth bool
}

const ContextUserKeyForGin = "zoox.connect::user"

// CreateGinMiddleware ...
func CreateGinMiddleware(cfg *CreateGinMiddlewareOptions) gin.HandlerFunc {
	var signer jwt.Jwt

	return func(ctx *gin.Context) {
		if signer == nil {
			signer = jwt.New(cfg.SecretKey)
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

			ctx.Set(ContextUserKeyForGinMiddleware, user)
		}

		if cfg.RequireAuth {
			if _, ok := ctx.Get(ContextUserKeyForGinMiddleware); !ok {
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
