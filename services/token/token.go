package token

import (
	"time"

	"github.com/go-zoox/connect/config"
	"github.com/go-zoox/jwt"
	"github.com/go-zoox/zoox"
)

var tokenKey = "gz_ut"

type TokenService struct {
}

func New() *TokenService {
	return &TokenService{}
}

func (ts *TokenService) Generate(cfg *config.Config, data map[string]any) (string, error) {
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

func (ts *TokenService) Verify(cfg *config.Config, ctx *zoox.Context, token string) bool {
	if token := ts.Get(ctx); token == "" {
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

func (ts *TokenService) Get(ctx *zoox.Context) string {
	return ctx.Cookie.Get(tokenKey)
}

func (ts *TokenService) Set(ctx *zoox.Context, token string) {
	ctx.Cookie.Set(tokenKey, token, 2*time.Hour)
}

func (ts *TokenService) Clear(ctx *zoox.Context) {
	ctx.Cookie.Del(tokenKey)
}
