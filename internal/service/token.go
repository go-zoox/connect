package service

import (
	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/jwt"
	"github.com/go-zoox/zoox"
)

var tokenKey = "gz_ut"
var providerKey = "gz_provider"

func GenerateToken(cfg *config.Config, data map[string]any) (string, error) {
	j := jwt.NewHS256(cfg.SecretKey)
	for k, v := range data {
		j.Set(k, v)
	}

	if token, err := j.Sign(); err != nil {
		return "", err
	} else {
		return token, nil
	}
}

func VerifyToken(cfg *config.Config, ctx *zoox.Context, token string) bool {
	if token := GetToken(ctx); token == "" {
		return false
	} else {
		j := jwt.NewHS256(cfg.SecretKey)
		if err := j.Verify(token); err != nil {
			return false
		} else {
			return true
		}
	}
}

func GetToken(ctx *zoox.Context) string {
	return ctx.Cookie.Get(tokenKey)
}

func SetToken(ctx *zoox.Context, cfg *config.Config, value string) {
	ctx.Cookie.Set(tokenKey, value, cfg.SessionMaxAgeDuration)
}

func DelToken(ctx *zoox.Context) {
	ctx.Cookie.Del(tokenKey)
}

// @TODO
func GetProvider(ctx *zoox.Context) string {
	return ctx.Cookie.Get(providerKey)
}

func SetProvider(ctx *zoox.Context, cfg *config.Config, value string) {
	ctx.Cookie.Set(providerKey, value, cfg.SessionMaxAgeDuration)
}

func DelProvider(ctx *zoox.Context) {
	ctx.Cookie.Del(providerKey)
}
