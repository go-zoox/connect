package main

import (
	"github.com/go-zoox/connect/internal"
	"github.com/go-zoox/connect/internal/config"
)

func main() {
	app := internal.New()
	var cfg *config.Config
	var err error
	if cfg, err = config.Load(); err != nil {
		panic(err)
	}

	if err := app.Start(cfg); err != nil {
		panic(err)
	}
}
