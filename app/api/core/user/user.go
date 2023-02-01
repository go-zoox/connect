package user

import (
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/errors"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		token := service.GetToken(ctx)
		if token == "" {
			ctx.Fail(fmt.Errorf("token is missing"), errors.FailedToGetToken.Code, errors.FailedToGetToken.Message)
			return
		}

		user, err := service.GetUser(ctx, cfg, token)
		if err != nil {
			ctx.Fail(err, errors.FailedToGetUser.Code, errors.FailedToGetUser.Message)
			return
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
			ctx.Fail(err, errors.InvalidJSON.Code, errors.InvalidJSON.Message)
			return
		}

		// ctx.Logger.Info("Login User: %v", user)

		if ok := service.ValidateCaptcha(cfg, ctx, user.Captcha); !ok {
			ctx.Fail(fmt.Errorf("invalid captcha"), errors.InvalidCaptcha.Code, errors.InvalidCaptcha.Message)
			return
		}

		token, err := service.Login(cfg, user.Type, user.Username, user.Password)
		if err != nil {
			// panic(errors.Wrap(err, "user login service failed"))
			// ctx.JSON(400, zoox.H{
			// 	"code":    400123,
			// 	"message": err.Error(),
			// })

			ctx.Fail(err, errors.UserLoginFailed.Code, errors.UserLoginFailed.Message)
			return
		}

		service.SetToken(ctx, cfg, token)

		ctx.Status(200)
	}
}

func GetUsers(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		page := ctx.Query().Get("page")
		pageSize := ctx.Query().Get("pageSize")
		if pageSize == "" {
			pageSize = ctx.Query().Get("page_size")
		}
		if page == "" {
			page = "1"
		}
		if pageSize == "" {
			pageSize = "10"
		}

		token := service.GetToken(ctx)
		if token == "" {
			ctx.Fail(fmt.Errorf("token is missing"), errors.FailedToGetToken.Code, errors.FailedToGetToken.Message)
			return
		}

		provider := service.GetProvider(ctx)
		if provider == "" {
			ctx.Fail(fmt.Errorf("provider is missing"), errors.FailedToGetOAuth2Provider.Code, errors.FailedToGetOAuth2Provider.Message)
			return
		}

		data, total, err := service.GetUsers(cfg, provider, token, page.String(), pageSize.String())
		if err != nil {
			ctx.Fail(err, errors.FailedToGetUsers.Code, errors.FailedToGetUsers.Message)
			return
		}

		ctx.Success(zoox.H{
			"data":  data,
			"total": total,
		})
	}
}
