package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/connect"
	"github.com/go-zoox/connect/cmd/connect/commands"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:        "connect",
		Usage:       "The Connector",
		Description: "Connect between auth with apps/services",
		Version:     connect.Version,
	})

	app.Register("server", commands.Server())

	app.Register("none", commands.None())

	app.Register("doremaon", commands.Doreamon())

	app.Register("github", commands.GitHub())

	app.Run()
}
