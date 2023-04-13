package service

import (
	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/crypto/jwt"
	"github.com/go-zoox/zoox"
)

var tokenKey = "gz_ut"
var providerKey = "gz_provider"

func GenerateToken(cfg *config.Config, data map[string]any) (string, error) {
	j := jwt.New(cfg.SecretKey)

	if token, err := j.Sign(data); err != nil {
		return "", err
	} else {
		return token, nil
	}
}

func VerifyToken(cfg *config.Config, ctx *zoox.Context, token string) bool {
	if token := GetToken(ctx); token == "" {
		return false
	} else {
		j := jwt.New(cfg.SecretKey)
		if _, err := j.Verify(token); err != nil {
			return false
		} else {
			return true
		}
	}
}

func GetToken(ctx *zoox.Context) string {
	cookieToken := ctx.Cookie().Get(tokenKey)
	if cookieToken != "" {
		return cookieToken
	}

	headerToken := ctx.Get("authorization")
	if headerToken != "" {
		// Bear token
		if len(headerToken) > 6 && headerToken[:6] == "Bearer" {
			return headerToken[7:]
		}

		// not standard
		return headerToken
	}

	queryToken := ctx.Query().Get("access_token").String()
	if queryToken != "" {
		return queryToken
	}

	return ""
}

func SetToken(ctx *zoox.Context, cfg *config.Config, value string) {
	ctx.Cookie().Set(tokenKey, value, cfg.SessionMaxAgeDuration)
}

func DelToken(ctx *zoox.Context) {
	ctx.Cookie().Del(tokenKey)
}

// @TODO
func GetProvider(ctx *zoox.Context) string {
	return ctx.Cookie().Get(providerKey)
}

func SetProvider(ctx *zoox.Context, cfg *config.Config, value string) {
	ctx.Cookie().Set(providerKey, value, cfg.SessionMaxAgeDuration)
}

func DelProvider(ctx *zoox.Context) {
	ctx.Cookie().Del(providerKey)
}
