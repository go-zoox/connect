package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/user"
	"github.com/go-zoox/crypto/jwt"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
)

// User ...
type User = user.User

// GetUser ...
func GetUser(ctx *zoox.Context, cfg *config.Config, token string) (*User, int, error) {
	cacheKey := fmt.Sprintf("user:%s", token)
	statusCode := 200

	u := new(User)
	if err := ctx.Cache().Get(cacheKey, u); err == nil {
		return u, statusCode, nil
	}

	if cfg.Services.User.Mode == "local" {
		userD := cfg.Services.User.Local

		u = &User{
			ID:          userD.ID,
			Username:    userD.Username,
			Nickname:    userD.Nickname,
			Email:       userD.Email,
			Permissions: userD.Permissions,
		}

		ctx.Cache().Set(cacheKey, u, cfg.SessionMaxAgeDuration)
		return u, statusCode, nil
	}

	response, err := fetch.Get(cfg.Services.User.Service, &fetch.Config{
		Headers: map[string]string{
			"accept":        "application/json",
			"authorization": fmt.Sprintf("Bearer %s", token),
		},
	})
	if err != nil {
		statusCode := 500
		return nil, statusCode, err
	}

	if response.Status != 200 {
		statusCode := response.Status
		return nil, statusCode, fmt.Errorf("failed to get user: (status: %d, response: %s)", response.Status, response.String())
	}

	userStr := response.Get("result").String()
	if err := json.Unmarshal([]byte(userStr), &u); err != nil {
		statusCode := 500
		return nil, statusCode, fmt.Errorf("failed to parse user with response.result: %v(response: %s)", err, response.String())
	}
	if u.ID == "" {
		u.ID = response.Get("result._id").String()
	}

	// Get OpenID: feishu
	u.FeishuOpenID, _, err = GetOpenID(ctx, cfg, "feishu", u.Email)
	if err != nil {
		time.Sleep(3 * time.Second)
		ctx.Logger.Warn("[service.user] failed to get feishu open id: %#v", err)
	}

	// if len(user.Permissions) != 0 {
	// 	ctx.Cache().Set(cacheKey, user, cfg.SessionMaxAgeDuration)
	// } else {
	// 	// no permission => 403 => cache 30s
	// 	ctx.Cache().Set(cacheKey, user, 30*time.Second)
	// }

	logger.Info("[service.GetUser] user: %s(%s)", u.Nickname, u.Email)
	ctx.Cache().Set(cacheKey, u, cfg.SessionMaxAgeDuration)

	return u, statusCode, nil
}

// Login ...
func Login(ctx *zoox.Context, cfg *config.Config, typ string, username string, password string) (string, error) {
	if cfg.Auth.Mode != "password" {
		panic("unsupported auth mode in login service")
	}

	if cfg.Password.Mode == "local" {
		if username != cfg.Password.Local.Username || password != cfg.Password.Local.Password {
			return "", fmt.Errorf("用户名或密码错误")
		}

		j := jwt.New(cfg.SecretKey)
		token, err := j.Sign(map[string]interface{}{
			"username": username,
		})
		if err != nil {
			return "", err
		}

		return token, nil
	}

	response, err := fetch.Post(cfg.Password.Service, &fetch.Config{
		Headers: map[string]string{
			"x-real-ip":       ctx.Get("x-forwarded-for"),
			"x-forwarded-for": ctx.Get("x-forwarded-for"),
			//
			"accept":       "application/json",
			"content-type": "application/json",
		},
		Body: map[string]string{
			"type":     typ,
			"username": username,
			"password": password,
		},
	})
	if err != nil {
		return "", err
	}

	return response.Get("access_token").String(), nil
}
