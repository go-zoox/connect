package user

import (
	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/connect/internal/errors"
	"github.com/go-zoox/connect/internal/service"
	"github.com/go-zoox/zoox"
)

func New(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		token := service.GetToken(ctx)
		if token == "" {
			ctx.Fail(errors.FailedToGetToken.Code, errors.FailedToGetToken.Message)
			return
		}

		user, err := service.GetUser(cfg, token)
		if err != nil {
			ctx.Fail(errors.FailedToGetUser.Code, errors.FailedToGetUser.Message+": "+err.Error())
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
			ctx.Fail(errors.InvalidJSON.Code, errors.InvalidJSON.Message)
			return
		}

		// ctx.Logger.Info("Login User: %v", user)

		if ok := service.ValidateCaptcha(cfg, ctx, user.Captcha); !ok {
			ctx.Fail(errors.InvalidCaptcha.Code, errors.InvalidCaptcha.Message)
			return
		}

		token, err := service.Login(cfg, user.Type, user.Username, user.Password)
		if err != nil {
			// panic(errors.Wrap(err, "user login service failed"))
			// ctx.JSON(400, zoox.H{
			// 	"code":    400123,
			// 	"message": err.Error(),
			// })

			ctx.Fail(errors.UserLoginFailed.Code, errors.UserLoginFailed.Message+": "+err.Error())
			return
		}

		service.SetToken(ctx, cfg, token)

		ctx.Status(200)
	}
}

func GetUsers(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		page := ctx.Query("page")
		pageSize := ctx.Query("pageSize")
		if pageSize == "" {
			pageSize = ctx.Query("page_size")
		}
		if page == "" {
			page = "1"
		}
		if pageSize == "" {
			pageSize = "10"
		}

		token := service.GetToken(ctx)
		if token == "" {
			ctx.Fail(errors.FailedToGetToken.Code, errors.FailedToGetToken.Message)
			return
		}

		provider := service.GetProvider(ctx)
		if provider == "" {
			ctx.Fail(errors.FailedToGetOAuth2Provider.Code, errors.FailedToGetOAuth2Provider.Message)
			return
		}

		data, total, err := service.GetUsers(cfg, provider, token, page, pageSize)
		if err != nil {
			ctx.Fail(errors.FailedToGetUsers.Code, errors.FailedToGetUsers.Message+": "+err.Error())
			return
		}

		ctx.Success(zoox.H{
			"data":  data,
			"total": total,
		})
	}
}
