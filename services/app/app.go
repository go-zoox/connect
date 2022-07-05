package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-zoox/connect/cache"
	"github.com/go-zoox/connect/config"
	"github.com/go-zoox/fetch"
)

type App struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Logo        string      `json:"logo"`
	Settings    AppSettings `json:"settings"`
}

type AppSettings struct {
	Functions []string `json:"functions"`
}

type AppService struct {
}

func New() *AppService {
	return &AppService{}
}

func (s *AppService) Get(cfg *config.Config, token string) (a *App, err error) {
	var app = new(App)
	if err = cache.Get("app", app); err == nil {
		return app, nil
	}

	if cfg.Services.App.Mode == "local" {
		appD := cfg.Services.App.Local

		app = &App{
			Name:        appD.Name,
			Description: appD.Description,
			Logo:        appD.Logo,
			Settings:    AppSettings(appD.Settings),
		}

		cache.Set("app", app, 2*time.Hour)
		return app, nil
	}

	clientID := cfg.Auth.OAuth2.ClientID
	clientSecret := cfg.Auth.OAuth2.ClientSecret

	response, err := fetch.Get(cfg.Services.App.Service, &fetch.Config{
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

	if err := json.Unmarshal([]byte(response.Get("result").String()), &app); err != nil {
		return nil, err
	}

	cache.Set("app", app, 2*time.Hour)
	return app, nil
}
