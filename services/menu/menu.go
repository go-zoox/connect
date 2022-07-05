package menu

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-zoox/connect/cache"
	"github.com/go-zoox/connect/config"
	"github.com/go-zoox/fetch"
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

type MenuService struct {
}

func New() *MenuService {
	return &MenuService{}
}

func (s *MenuService) Get(cfg *config.Config, token string) ([]MenuItem, error) {
	cacheKey := fmt.Sprintf("menus:%s", token)

	var menus []MenuItem
	if err := cache.Get(cacheKey, &menus); err == nil {
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

		cache.Set(cacheKey, &menus, 2*time.Hour)
		return menus, nil
	}

	clientID := cfg.Auth.OAuth2.ClientID
	clientSecret := cfg.Auth.OAuth2.ClientSecret

	response, err := fetch.Get(cfg.Services.Menus.Service, &fetch.Config{
		Headers: map[string]string{
			"accept":          "application/json",
			"authorization":   fmt.Sprintf("Bearer %s", token),
			"x-client-id":     clientID,
			"x-client-secret": clientSecret,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(response.Get("result").String()), &menus); err != nil {
		return nil, err
	}

	if len(menus) != 0 {
		cache.Set(cacheKey, &menus, 2*time.Hour)
	} else {
		// no menus => 403 => cache 30s
		cache.Set(cacheKey, &menus, 30*time.Second)
	}

	return menus, nil
}
