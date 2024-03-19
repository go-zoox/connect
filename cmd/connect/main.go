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

	// full mode
	app.Register("server", commands.Server())

	// none mode
	app.Register("none", commands.None())

	// passport mode @TODO
	// app.Register("passport", commands.Passport())

	// doreamon mode (oauth2 + users + menus + permissions)
	app.Register("doremaon", commands.Doreamon())

	// oauth2 mode
	app.Register("github", commands.GitHub())
	app.Register("feishu", commands.Feishu())

	app.Run()
}
