package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/zoox"
)

// PermissionItem ...
type PermissionItem = string

// GetPermission ...
func GetPermission(ctx *zoox.Context, cfg *config.Config, provider string, token string) ([]PermissionItem, int, error) {
	cacheKey := fmt.Sprintf("permissions:%s", token)
	statusCode := 200

	var permissions []PermissionItem
	if err := ctx.Cache().Get(cacheKey, &permissions); err == nil {
		return permissions, statusCode, nil
	}

	if cfg.Services.Permissions.Mode == "local" {
		for _, permission := range cfg.Services.Permissions.Local {
			permissions = append(permissions, PermissionItem(permission))
		}

		ctx.Cache().Set(cacheKey, &permissions, cfg.GetSessionMaxAgeDuration())
		return permissions, statusCode, nil
	}

	clientCfg, err := oauth2.Get(provider)
	if err != nil {
		statusCode := 500
		return nil, statusCode, err
	}

	response, err := fetch.Get(cfg.Services.Permissions.Service, &fetch.Config{
		Headers: map[string]string{
			"x-real-ip":       ctx.Get("x-forwarded-for"),
			"x-forwarded-for": ctx.Get("x-forwarded-for"),
			//
			"accept":          "application/json",
			"authorization":   fmt.Sprintf("Bearer %s", token),
			"x-client-id":     clientCfg.ClientID,
			"x-client-secret": clientCfg.ClientSecret,
		},
	})
	if err != nil {
		statusCode := 500
		return nil, statusCode, err
	}

	if response.Status != 200 {
		statusCode := response.Status
		return nil, statusCode, fmt.Errorf("failed to get permissions: (status: %d, response: %s)", response.Status, response.String())
	}

	if err := json.Unmarshal([]byte(response.Get("result").String()), &permissions); err != nil {
		statusCode := 500
		return nil, statusCode, fmt.Errorf("failed to parse permissions with response.result: %v(response: %s)", err, response.String())
	}

	if len(permissions) != 0 {
		ctx.Cache().Set(cacheKey, &permissions, cfg.GetSessionMaxAgeDuration())
	} else {
		// no permissions => 403 => cache 30s
		ctx.Cache().Set(cacheKey, &permissions, 30*time.Second)
	}

	return permissions, statusCode, nil
}
