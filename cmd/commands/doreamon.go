package commands

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/connect/app"
	"github.com/go-zoox/connect/apps/doreamon"
	"github.com/go-zoox/core-utils/fmt"
)

func Doreamon() *cli.Command {
	return &cli.Command{
		Name:  "doreamon",
		Usage: "Start Connect Server using Doreamon",
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
				Name:    "client-id",
				Usage:   "Doreamon Client ID",
				EnvVars: []string{"CLIENT_ID"},
			},
			&cli.StringFlag{
				Name:    "client-secret",
				Usage:   "Doreamon Client Secret",
				EnvVars: []string{"CLIENT_SECRET"},
			},
			&cli.StringFlag{
				Name:    "redirect-uri",
				Usage:   "Doreamon Client Secret",
				EnvVars: []string{"REDIRECT_URI"},
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
		},
		Action: func(c *cli.Context) error {
			cfg, err := doreamon.Create(&doreamon.Config{
				Port:          c.Int64("port"),
				SecretKey:     c.String("secret-key"),
				SessionMaxAge: c.Int64("session-max-age"),
				ClientID:      c.String("client-id"),
				ClientSecret:  c.String("client-secret"),
				RedirectURI:   c.String("redirect-uri"),
				Frontend:      c.String("frontend"),
				Backend:       c.String("backend"),
				Upstream:      c.String("upstream"),
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
