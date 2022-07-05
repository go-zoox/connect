package user

import (
	"github.com/go-zoox/connect/config"
	"github.com/go-zoox/connect/errors"
	"github.com/go-zoox/connect/services"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		token := services.Token.Get(ctx)
		user, err := services.User.Get(cfg, token)
		if err != nil {
			panic(err)
		}

		ctx.JSON(200, zoox.H{
			"code":    0,
			"message": "",
			"result":  user,
		})
	}
}

func Login(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		type UserDTO struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Type     string `json:"type"`
			Captcha  string `json:"captcha"`
		}

		var user UserDTO
		if err := ctx.BindJSON(&user); err != nil {
			panic(err)
		}

		ctx.Logger.Info("Login User: %v", user)

		if ok := services.Captcha.Validate(cfg, ctx, user.Captcha); !ok {
			ctx.Fail(errors.InvalidCaptcha.Code, errors.InvalidCaptcha.Message)
			return
		}

		token, err := services.User.Login(cfg, user.Type, user.Username, user.Password)
		if err != nil {
			// panic(errors.Wrap(err, "user login service failed"))
			// ctx.JSON(400, zoox.H{
			// 	"code":    400123,
			// 	"message": err.Error(),
			// })

			ctx.Fail(errors.UserLoginFailed.Code, errors.UserLoginFailed.Message+": "+err.Error())
			return
		}

		services.Token.Set(ctx, token)

		ctx.Status(200)
	}
}
