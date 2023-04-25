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

type MenuItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Icon       string `json:"icon"`
	Sort       int64  `json:"sort"`
	IsHidden   bool   `json:"hidden"`
	IsExpended bool   `json:"expended"`
	Layout     string `json:"layout"`
	IFrame     string `json:"iframe"`
	Redirect   string `json:"redirect"`
}

func GetMenu(ctx *zoox.Context, cfg *config.Config, provider string, token string) ([]MenuItem, int, error) {
	cacheKey := fmt.Sprintf("menus:%s", token)
	statusCode := 200

	var menus []MenuItem
	if err := ctx.Cache().Get(cacheKey, &menus); err == nil {
		return menus, statusCode, nil
	}

	if cfg.Services.Menus.Mode == "local" {
		for _, menu := range cfg.Services.Menus.Local {
			menus = append(menus, MenuItem{
				ID:         menu.Path,
				Name:       menu.Name,
				Path:       menu.Path,
				Icon:       menu.Icon,
				Sort:       menu.Sort,
				IsHidden:   menu.IsHidden,
				IsExpended: menu.IsExpanded,
				Layout:     menu.Layout,
				IFrame:     menu.IFrame,
				Redirect:   menu.Redirect,
			})
		}

		ctx.Cache().Set(cacheKey, &menus, cfg.SessionMaxAgeDuration)
		return menus, statusCode, nil
	}

	clientCfg, err := oauth2.Get(provider)
	if err != nil {
		statusCode := 500
		return nil, statusCode, err
	}

	response, err := fetch.Get(cfg.Services.Menus.Service, &fetch.Config{
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
		return nil, statusCode, fmt.Errorf("failed to get menus: (status: %d, response: %s)", response.Status, response.String())
	}

	if err := json.Unmarshal([]byte(response.Get("result").String()), &menus); err != nil {
		statusCode := 500
		return nil, statusCode, fmt.Errorf("failed to parse menus with response.result: %v(response: %s)", err, response.String())
	}

	if len(menus) != 0 {
		ctx.Cache().Set(cacheKey, &menus, cfg.SessionMaxAgeDuration)
	} else {
		// no menus => 403 => cache 30s
		ctx.Cache().Set(cacheKey, &menus, 30*time.Second)
	}

	return menus, statusCode, nil
}
