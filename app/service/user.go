package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/pkg/cache"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/jwt"
)

type User struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Avatar      string   `json:"avatar"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
}

func GetUser(cfg *config.Config, token string) (*User, error) {
	cacheKey := fmt.Sprintf("user:%s", token)

	user := new(User)
	if err := cache.Get(cacheKey, user); err == nil {
		return user, nil
	}

	if cfg.Services.User.Mode == "local" {
		userD := cfg.Services.User.Local

		user = &User{
			ID:          userD.ID,
			Username:    userD.Username,
			Nickname:    userD.Nickname,
			Email:       userD.Email,
			Permissions: userD.Permissions,
		}

		cache.Set(cacheKey, user, cfg.SessionMaxAgeDuration)
		return user, nil
	}

	response, err := fetch.Get(cfg.Services.User.Service, &fetch.Config{
		Headers: map[string]string{
			"accept":        "application/json",
			"authorization": fmt.Sprintf("Bearer %s", token),
		},
	})
	if err != nil {
		return nil, err
	}

	userStr := response.Get("result").String()
	if err := json.Unmarshal([]byte(userStr), &user); err != nil {
		return nil, err
	}
	if user.ID == "" {
		user.ID = response.Get("result._id").String()
	}

	if len(user.Permissions) != 0 {
		cache.Set(cacheKey, user, cfg.SessionMaxAgeDuration)
	} else {
		// no permission => 403 => cache 30s
		cache.Set(cacheKey, user, 30*time.Second)
	}
	return user, nil
}

func Login(cfg *config.Config, typ string, username string, password string) (string, error) {
	if cfg.Auth.Mode != "password" {
		panic("unsupported auth mode in login service")
	}

	if cfg.Password.Mode == "local" {
		if username != cfg.Password.Local.Username || password != cfg.Password.Local.Password {
			// return "", fmt.Errorf("username or password are not matched")
			return "", fmt.Errorf("用户名或密码错误")
		}

		j := jwt.NewHS256(cfg.SecretKey)
		j.Set("username", username)
		token, err := j.Sign()
		if err != nil {
			return "", err
		}

		return token, nil
	}

	response, err := fetch.Post(cfg.Password.Service, &fetch.Config{
		Headers: map[string]string{
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
