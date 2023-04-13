package commands

import (
	"fmt"
	"log"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/connect/internal"
	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/fs"
)

func Server() *cli.Command {
	return &cli.Command{
		Name:        "serve",
		Usage:       "Start Connect Server",
		Description: "Start the connect server with config",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:    "port",
				Value:   8080,
				Usage:   "The port to listen on",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
			},
			&cli.StringFlag{
				Name:    "config",
				Usage:   "The config file",
				Aliases: []string{"c"},
				EnvVars: []string{"CONFIG"},
			},
		},
		Action: func(c *cli.Context) error {
			configFile := c.String("config")
			if configFile == "" {
				dotConfig := fs.JoinPath(fs.CurrentDir(), "conf/config.yml")
				if fs.IsExist(dotConfig) {
					configFile = dotConfig
				} else {
					log.Fatal(fmt.Errorf("config file(conf/config.yml) not found"))
				}
			}

			var cfg *config.Config
			var err error
			if cfg, err = config.Load(configFile); err != nil {
				log.Fatal(fmt.Errorf("failed to load config (%s, %s)", configFile, err))
			}

			cfg.Port = c.Int64("port")

			app := internal.New()
			if err := app.Start(cfg); err != nil {
				log.Fatal(fmt.Errorf("failed to start server(err: %s)", err))
			}

			return nil
		},
	}
}
