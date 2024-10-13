package service

import (
	"encoding/json"
	"fmt"

	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/zoox"
)

// App ...
type App struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Logo        string      `json:"logo"`
	Settings    AppSettings `json:"settings"`
}

// AppSettings ...
type AppSettings struct {
	Functions any `json:"functions"`
}

// GetApp ...
func GetApp(ctx *zoox.Context, cfg *config.Config, provider string, token string) (*App, int, error) {
	var app = new(App)
	if err := ctx.Cache().Get("app", app); err == nil {
		return app, 200, nil
	}

	if cfg.Services.App.Mode == "local" {
		appD := cfg.Services.App.Local

		app = &App{
			Name:        appD.Name,
			Description: appD.Description,
			Logo:        appD.Logo,
			Settings:    AppSettings(appD.Settings),
		}

		ctx.Cache().Set("app", app, cfg.GetSessionMaxAgeDuration())
		return app, 200, nil
	}

	clientCfg, err := oauth2.Get(provider)
	if err != nil {
		return nil, 500, err
	}

	response, err := fetch.Get(cfg.Services.App.Service, &fetch.Config{
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
		return nil, 500, err
	}

	if response.Status != 200 {
		statusCode := response.Status
		return nil, statusCode, fmt.Errorf("failed to get app: (status: %d, response: %s)", response.Status, response.String())
	}

	if err := json.Unmarshal([]byte(response.Get("result").String()), &app); err != nil {
		return nil, 500, fmt.Errorf("unmarshal app: %s (response: %s)", err, response.String())
	}

	ctx.Cache().Set("app", app, cfg.GetSessionMaxAgeDuration())
	return app, 200, nil
}
