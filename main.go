package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-zoox/cli"
	internal "github.com/go-zoox/connect/app"
	"github.com/go-zoox/connect/app/config"
	"github.com/go-zoox/fs"
)

func main() {
	app := cli.NewSingleProgram(&cli.SingleProgramConfig{
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
	})

	app.Command(func(c *cli.Context) error {
		// port := c.String("port")
		// if os.Getenv("PORT") != "" {
		// 	port = os.Getenv("PORT")
		// }

		configFile := c.String("config")
		if os.Getenv("CONFIG") != "" {
			configFile = os.Getenv("CONFIG")
		}

		if configFile == "" {
			dotConfig := fs.JoinPath(fs.CurrentDir(), ".config.yml")
			if fs.IsExist(dotConfig) {
				configFile = dotConfig
			} else {
				log.Fatal(fmt.Errorf("config file(.config.yml) not found"))
			}
		}

		app := internal.New()
		var cfg *config.Config
		var err error
		if cfg, err = config.Load(configFile); err != nil {
			log.Fatal(fmt.Errorf("failed to load config (%s, %s)", configFile, err))
		}

		if err := app.Start(cfg); err != nil {
			log.Fatal(fmt.Errorf("failed to start server(err: %s)", err))
		}

		return nil
	})

	app.Run()
}
