package commands

import (
	"log"
	"os"

	"github.com/go-zoox/core-utils/fmt"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/connect/app"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/fs"
)

// Server ...
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
			if os.Getenv("CONFIG") != "" {
				configFile = os.Getenv("CONFIG")
			}

			if configFile == "" {
				dotConfig := fs.JoinPath(fs.CurrentDir(), "conf/config.yml")
				if fs.IsExist(dotConfig) {
					configFile = dotConfig
				} else {
					log.Fatal(fmt.Errorf("config file(conf/config.yml) not found"))
				}
			}

			app := app.New()
			var cfg *config.Config
			var err error
			if cfg, err = config.Load(configFile); err != nil {
				log.Fatal(fmt.Errorf("failed to load config (%s, %s)", configFile, err))
			}

			if c.IsSet("port") {
				cfg.Port = c.Int64("port")
			}

			if err := app.Start(cfg); err != nil {
				log.Fatal(fmt.Errorf("failed to start server(err: %s)", err))
			}

			return nil
		},
	}
}
