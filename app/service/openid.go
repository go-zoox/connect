package service

import (
	"encoding/json"
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/zoox"
)

// OpenID ...
type OpenID struct {
	OpenID string `json:"openID"`
}

// GetOpenID ...
func GetOpenID(ctx *zoox.Context, cfg *config.Config, provider string, email string) (string, int, error) {
	cacheKey := fmt.Sprintf("open_id:%s", email)
	statusCode := 200

	var instance = new(OpenID)
	if err := ctx.Cache().Get(cacheKey, instance); err == nil {
		return instance.OpenID, statusCode, nil
	}

	if cfg.Services.OpenID.Mode == "local" {
		appD := cfg.Services.OpenID.Local

		instance = &OpenID{
			OpenID: appD.OpenID,
		}

		ctx.Cache().Set(cacheKey, instance, cfg.SessionMaxAgeDuration)
		return instance.OpenID, statusCode, nil
	}

	oauth2Provider := GetProvider(ctx)
	if provider == "" {
		statusCode = 400
		return "", statusCode, fmt.Errorf("oauth2 provider is missing")
	}

	clientCfg, err := oauth2.Get(oauth2Provider)
	if err != nil {
		statusCode = 500
		return "", statusCode, err
	}

	response, err := fetch.Get(cfg.Services.OpenID.Service, &fetch.Config{
		Headers: map[string]string{
			"x-real-ip":       ctx.Get("x-forwarded-for"),
			"x-forwarded-for": ctx.Get("x-forwarded-for"),
			//
			"Accept":          "application/json",
			"X-Client-ID":     clientCfg.ClientID,
			"X-Client-Secret": clientCfg.ClientSecret,
		},
		Query: map[string]string{
			"email":    email,
			"provider": provider,
		},
	})
	if err != nil {
		statusCode = 500
		return "", statusCode, err
	}

	if response.Status != 200 {
		statusCode := response.Status
		return "", statusCode, fmt.Errorf("failed to get openid: (status: %d, response: %s)", response.Status, response.String())
	}

	if response.Get("result").String() != "" {
		if err := json.Unmarshal([]byte(response.Get("result").String()), &instance); err != nil {
			statusCode := 500
			return "", statusCode, fmt.Errorf("failed to parse open_id with response.result: %v(response: %s)", err, response.String())
		}
	}

	logger.Info("[service.GetOpenID][%s: %s] open_id: %s", email, provider, response.String())
	ctx.Cache().Set(cacheKey, instance, cfg.SessionMaxAgeDuration)
	return instance.OpenID, statusCode, nil
}
