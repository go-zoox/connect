package service

import (
	"os"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/crypto/jwt"
	"github.com/go-zoox/zoox"
)

var tokenKey = "gz_ut"
var refreshTokenKey = "gz_rt"
var providerKey = "gz_provider"

// X_CONNECT_TOKEN_COMPITABLE_WITH_AUTHORIZATION_HEADER
var isXConnectTokenCompitableWithAuthorizationHeader = os.Getenv("X_CONNECT_TOKEN_COMPITABLE_WITH_AUTHORIZATION_HEADER") == "true"

// X_CONNECT_TOKEN_COMPITABLE_WITH_QUERY
var isXConnectTokenCompitableWithQuery = os.Getenv("X_CONNECT_TOKEN_COMPITABLE_WITH_QUERY") == "true"

// GenerateToken ...
func GenerateToken(cfg *config.Config, data map[string]any) (string, error) {
	j := jwt.New(cfg.SecretKey)

	token, err := j.Sign(data)
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyToken ...
func VerifyToken(cfg *config.Config, ctx *zoox.Context, token string) bool {
	if token := GetToken(ctx); token == "" {
		return false
	}

	j := jwt.New(cfg.SecretKey)
	if _, err := j.Verify(token); err != nil {
		return false
	}

	return true
}

// GetToken ...
func GetToken(ctx *zoox.Context) string {
	sessionToken := ctx.Session().Get(tokenKey)
	if sessionToken != "" {
		return sessionToken
	}

	if ctx.Header().Get("X-Connect-Token") != "" {
		return ctx.Header().Get("X-Connect-Token")
	}

	if ctx.Query().Get("X-Connect-Token").String() != "" {
		return ctx.Query().Get("X-Connect-Token").String()
	}

	if isXConnectTokenCompitableWithAuthorizationHeader {
		headerToken := ctx.Get("authorization")
		if headerToken != "" {
			// Bear token
			if len(headerToken) > 6 && headerToken[:6] == "Bearer" {
				return headerToken[7:]
			}

			// not standard
			return headerToken
		}
	}

	if isXConnectTokenCompitableWithQuery {
		queryToken := ctx.Query().Get("access_token").String()
		if queryToken != "" {
			return queryToken
		}
	}

	return ""
}

// SetToken ...
func SetToken(ctx *zoox.Context, cfg *config.Config, value string) {
	ctx.Session().Set(tokenKey, value)
}

// DelToken ...
func DelToken(ctx *zoox.Context) {
	ctx.Session().Del(tokenKey)
}

// GetRefreshToken ...
func GetRefreshToken(ctx *zoox.Context) string {
	return ctx.Session().Get(refreshTokenKey)
}

// SetRefreshToken ...
func SetRefreshToken(ctx *zoox.Context, cfg *config.Config, value string) {
	ctx.Session().Set(refreshTokenKey, value)
}

// DelRefreshToken ...
func DelRefreshToken(ctx *zoox.Context) {
	ctx.Session().Del(refreshTokenKey)
}

// GetProvider ...
func GetProvider(ctx *zoox.Context) string {
	return ctx.Session().Get(providerKey)
}

// SetProvider ...
func SetProvider(ctx *zoox.Context, cfg *config.Config, value string) {
	ctx.Session().Set(providerKey, value)
}

// DelProvider ...
func DelProvider(ctx *zoox.Context) {
	ctx.Session().Del(providerKey)
}
