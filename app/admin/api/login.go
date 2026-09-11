package api

import (
	"fmt"
	"net/http"

	"github.com/go-zoox/connect/app/admin/repository"
	adminsvc "github.com/go-zoox/connect/app/admin/service"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/app/errors"
	"github.com/go-zoox/connect/app/service"
	"github.com/go-zoox/gormx"
	"github.com/go-zoox/zoox"
)

// Login validates credentials against the admin users table (gormx DB) and,
// on success, issues a session token the same way as password-mode core user login (JWT + SetToken).
func Login(cfg *config.Config) zoox.HandlerFunc {
	return func(ctx *zoox.Context) {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := ctx.BindJSON(&body); err != nil {
			ctx.Fail(err, errors.InvalidJSON.Code, errors.InvalidJSON.Message)
			return
		}

		db := gormx.GetDB()
		if db == nil {
			ctx.Fail(fmt.Errorf("admin database not initialized"), 500201, "admin database not available", http.StatusInternalServerError)
			return
		}

		auth := adminsvc.NewAuthService(repository.NewUserRepo(db))
		res, err := auth.Login(body.Username, body.Password)
		if err != nil {
			ctx.Fail(err, 500202, "login failed", http.StatusInternalServerError)
			return
		}
		if res.Status == adminsvc.AuthLoginInvalidCredentials {
			ctx.Fail(fmt.Errorf("invalid credentials"), errors.AdminLoginInvalid.Code, errors.AdminLoginInvalid.Message, http.StatusUnauthorized)
			return
		}

		token, err := service.GenerateToken(cfg, map[string]any{
			"username": res.User.Username,
		})
		if err != nil {
			ctx.Fail(err, 500203, "failed to issue token", http.StatusInternalServerError)
			return
		}

		service.SetToken(ctx, cfg, token)
		ctx.Status(http.StatusOK)
	}
}
