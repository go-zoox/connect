package main

import (
	_ "embed"

	"github.com/go-zoox/connect"
	"github.com/go-zoox/connect/config"
)

func main() {
	app := connect.New()
	var cfg *config.Config
	var err error
	if cfg, err = config.Load(); err != nil {
		panic(err)
	}

	if err := app.Start(cfg); err != nil {
		panic(err)
	}
}
