package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-zoox/connect/internal/config"
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

func GetMenu(ctx *zoox.Context, cfg *config.Config, provider string, token string) ([]MenuItem, error) {
	cacheKey := fmt.Sprintf("menus:%s", token)

	var menus []MenuItem
	if err := ctx.Cache().Get(cacheKey, &menus); err == nil {
		return menus, nil
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
		return menus, nil
	}

	clientCfg, err := oauth2.Get(provider)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if err := json.Unmarshal([]byte(response.Get("result").String()), &menus); err != nil {
		return nil, err
	}

	if len(menus) != 0 {
		ctx.Cache().Set(cacheKey, &menus, cfg.SessionMaxAgeDuration)
	} else {
		// no menus => 403 => cache 30s
		ctx.Cache().Set(cacheKey, &menus, 30*time.Second)
	}

	return menus, nil
}
