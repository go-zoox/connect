package service

import (
	"encoding/json"
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/connect/pkg/cache"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/zoox"
)

type OpenID struct {
	OpenID string `json:"openID"`
}

func GetOpenID(ctx *zoox.Context, cfg *config.Config, provider string, email string) (string, error) {
	cacheKey := fmt.Sprintf("open_id:%s", email)

	var instance = new(OpenID)
	if err := cache.Get(cacheKey, instance); err == nil {
		return instance.OpenID, nil
	}

	if cfg.Services.OpenID.Mode == "local" {
		appD := cfg.Services.OpenID.Local

		instance = &OpenID{
			OpenID: appD.OpenID,
		}

		cache.Set(cacheKey, instance, cfg.SessionMaxAgeDuration)
		return instance.OpenID, nil
	}

	oauth2Provider := GetProvider(ctx)
	if provider == "" {
		return "", fmt.Errorf("oauth2 provider is missing")
	}

	clientCfg, err := oauth2.Get(oauth2Provider)
	if err != nil {
		return "", err
	}

	response, err := fetch.Get(cfg.Services.OpenID.Service, &fetch.Config{
		Headers: map[string]string{
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
		return "", err
	}

	if response.Get("result").String() != "" {
		if err := json.Unmarshal([]byte(response.Get("result").String()), &instance); err != nil {
			return "", fmt.Errorf("unmarshal open_id: %s (response: %s)", err, response.String())
		}
	}

	logger.Info("[service.GetOpenID][%s: %s] open_id: %s", email, provider, response.String())
	cache.Set(cacheKey, instance, cfg.SessionMaxAgeDuration)
	return instance.OpenID, nil
}
