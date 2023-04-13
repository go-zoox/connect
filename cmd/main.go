package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/connect"
	"github.com/go-zoox/connect/cmd/commands"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:        "connect",
		Usage:       "The Connector",
		Description: "Connect between auth with apps/services",
		Version:     connect.Version,
	})

	app.Register("server", commands.Server())
	app.Register("doremaon", commands.Doreamon())

	app.Run()
}
