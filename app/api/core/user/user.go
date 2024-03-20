package user

import (
	"fmt"
	"net/http"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/errors"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/zoox"
)

// New ...
func New(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		if cfg.Auth.Mode == "none" {
			ctx.JSON(200, zoox.H{
				"code":    0,
				"message": "",
				"result": &service.User{
					ID:       "1",
					Nickname: "anonymous",
					Username: "anonymous",
					Email:    "anonymous@gozoox.com",
				},
			})
			return
		}

		token := service.GetToken(ctx)
		if token == "" {
			ctx.Fail(fmt.Errorf("[api.core.user] token is missing (1)"), errors.FailedToGetToken.Code, errors.FailedToGetToken.Message)
			return
		}

		user, statusCode, err := service.GetUser(ctx, cfg, token)
		if err != nil {
			// @TODO
			if statusCode == http.StatusUnauthorized {
				service.DelToken(ctx)
			}

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

// Login ...
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

		token, err := service.Login(ctx, cfg, user.Type, user.Username, user.Password)
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

// GetUsers ...
func GetUsers(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		page := ctx.Query().Get("page").String()
		pageSize := ctx.Query().Get("pageSize").String()
		if pageSize == "" {
			pageSize = ctx.Query().Get("page_size").String()
		}
		if page == "" {
			page = "1"
		}
		if pageSize == "" {
			pageSize = "10"
		}

		token := service.GetToken(ctx)
		if token == "" {
			ctx.Fail(fmt.Errorf("[api.core.user] token is missing (2)"), errors.FailedToGetToken.Code, errors.FailedToGetToken.Message)
			return
		}

		provider := service.GetProvider(ctx)
		if provider == "" {
			ctx.Fail(fmt.Errorf("provider is missing"), errors.FailedToGetOAuth2Provider.Code, errors.FailedToGetOAuth2Provider.Message)
			return
		}

		data, total, _, err := service.GetUsers(ctx, cfg, provider, token, page, pageSize)
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
