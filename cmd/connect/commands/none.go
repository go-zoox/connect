package commands

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/connect/app"
	"github.com/go-zoox/connect/apps/none"
	"github.com/go-zoox/core-utils/fmt"
)

// None ...
func None() *cli.Command {
	return &cli.Command{
		Name:  "none",
		Usage: "Start Connect Server using None Auth",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:    "port",
				Value:   8080,
				Usage:   "The port to listen on",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
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
			cfg, err := none.Create(&none.Config{
				Port:     c.Int64("port"),
				Frontend: c.String("frontend"),
				Backend:  c.String("backend"),
				Upstream: c.String("upstream"),
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
