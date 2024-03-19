package commands

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/connect/app"
	"github.com/go-zoox/connect/apps/feishu"
	"github.com/go-zoox/core-utils/fmt"
)

// Feishu ...
func Feishu() *cli.Command {
	return &cli.Command{
		Name:  "feishu",
		Usage: "Start Connect Server using Feishu",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:    "port",
				Value:   8080,
				Usage:   "The port to listen on",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
			},
			&cli.StringFlag{
				Name:    "secret-key",
				Usage:   "Secret Key",
				EnvVars: []string{"SESSION_KEY", "SECRET_KEY"},
			},
			&cli.Int64Flag{
				Name:    "session-max-age",
				Usage:   "Session Max Age",
				EnvVars: []string{"SESSION_MAX_AGE"},
				Value:   86400,
			},
			&cli.StringFlag{
				Name:     "client-id",
				Usage:    "Feishu Client ID",
				EnvVars:  []string{"CLIENT_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "client-secret",
				Usage:    "Feishu Client Secret",
				EnvVars:  []string{"CLIENT_SECRET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "redirect-uri",
				Usage:    "Feishu Client Secret",
				EnvVars:  []string{"REDIRECT_URI"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "frontend",
				Usage:   "frontend service",
				EnvVars: []string{"FRONTEND"},
			},
			&cli.StringFlag{
				Name:    "backend",
				Usage:   "backend service",
				EnvVars: []string{"BACKEND"},
			},
			&cli.StringFlag{
				Name:    "upstream",
				Usage:   "upstream service",
				EnvVars: []string{"UPSTREAM"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "Debug mode show config info",
				EnvVars: []string{"DEBUG"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "backend-prefix",
				Usage:   "backend prefix",
				EnvVars: []string{"BACKEND_PREFIX"},
				Value:   "",
			},
			&cli.BoolFlag{
				Name:    "backend-disable-prefix-rewrite",
				Usage:   "backend disable prefix rewrite",
				EnvVars: []string{"BACKEND_DISABLE_PREFIX_REWRITE"},
			},
		},
		Action: func(c *cli.Context) error {
			cfg, err := feishu.Create(&feishu.Config{
				Port:          c.Int64("port"),
				SecretKey:     c.String("secret-key"),
				SessionMaxAge: c.Int64("session-max-age"),
				ClientID:      c.String("client-id"),
				ClientSecret:  c.String("client-secret"),
				RedirectURI:   c.String("redirect-uri"),
				Frontend:      c.String("frontend"),
				Backend:       c.String("backend"),
				Upstream:      c.String("upstream"),
				//
				BackendPrefix:                 c.String("backend-prefix"),
				BackendIsDisablePrefixRewrite: c.Bool("backend-disable-prefix-rewrite"),
			})
			if err != nil {
				return err
			}

			if c.Bool("debug") {
				fmt.PrintJSON("config:", cfg)
			}

			return app.New().Start(cfg)
		},
	}
}
