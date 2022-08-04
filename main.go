package main

import (
	"os"

	"github.com/go-zoox/connect/internal"
	"github.com/go-zoox/connect/internal/config"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/logger"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "Serve",
		Usage:       "The Serve",
		Description: "Server static files",
		// Version:     Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Value:   "8080",
				Usage:   "The port to listen on",
				Aliases: []string{"p"},
			},
			&cli.StringFlag{
				Name:    "config",
				Usage:   "The config file",
				Aliases: []string{"c"},
			},
		},
		Action: func(c *cli.Context) error {
			// port := c.String("port")
			// if os.Getenv("PORT") != "" {
			// 	port = os.Getenv("PORT")
			// }

			config_file := c.String("config")
			if os.Getenv("CONFIG") != "" {
				config_file = os.Getenv("CONFIG")
			}

			if config_file == "" {
				dotConfig := fs.JoinPath(fs.CurrentDir(), ".config.yml")
				if fs.IsExist(dotConfig) {
					config_file = dotConfig
				} else {
					panic("config file is required")
				}
			}

			app := internal.New()
			var cfg *config.Config
			var err error
			if cfg, err = config.Load(config_file); err != nil {
				panic(err)
			}

			if err := app.Start(cfg); err != nil {
				panic(err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal("%s", err.Error())
	}
}
