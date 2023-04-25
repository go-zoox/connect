package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/zoox"
)

type UserX struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

func GetUsers(ctx *zoox.Context, cfg *config.Config, provider string, token string, page, pageSize string) ([]*User, int64, int, error) {
	key := fmt.Sprintf("users:%s:%s:%s:%s", provider, token, page, pageSize)
	statusCode := 200

	var users []*User
	// if err = cache.Get(key, &users); err == nil {
	// 	return users, nil
	// }

	if cfg.Services.App.Mode == "local" {
		return nil, 0, statusCode, errors.New("unsupport in local mode")
	}

	clientCfg, err := oauth2.Get(provider)
	if err != nil {
		statusCode = 500
		return nil, 0, statusCode, err
	}

	response, err := fetch.Get(cfg.Services.Users.Service, &fetch.Config{
		Headers: map[string]string{
			"x-real-ip":       ctx.Get("x-forwarded-for"),
			"x-forwarded-for": ctx.Get("x-forwarded-for"),
			//
			"accept":          "application/json",
			"authorization":   fmt.Sprintf("Bearer %s", token),
			"x-client-id":     clientCfg.ClientID,
			"x-client-secret": clientCfg.ClientSecret,
		},
		Query: map[string]string{
			"page":     page,
			"pageSize": pageSize,
		},
	})
	if err != nil {
		statusCode = 500
		return nil, 0, statusCode, err
	}

	if response.Status != 200 {
		statusCode := response.Status
		return nil, 0, statusCode, fmt.Errorf("failed to get users: (status: %d, response: %s)", response.Status, response.String())
	}

	var usersX []*UserX
	if err := json.Unmarshal([]byte(response.Get("result.data").String()), &usersX); err != nil {
		statusCode := 500
		return nil, 0, statusCode, fmt.Errorf("failed to parse menus with response.result.data: %v(response: %s)", err, response.String())
	}

	total := response.Get("result.total").Int()

	for _, u := range usersX {
		users = append(users, &User{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Email:    u.Email,
		})
	}

	ctx.Cache().Set(key, &users, 10*time.Second)

	return users, total, statusCode, nil
}
