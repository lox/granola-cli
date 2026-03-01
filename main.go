package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/lox/granola-cli/cmd"
	"github.com/lox/granola-cli/internal/cli"
)

var version = "dev"

func main() {
	c := &cmd.CLI{}
	ctx := kong.Parse(c,
		kong.Name("granola"),
		kong.Description("A CLI for Granola meeting notes"),
		kong.UsageOnError(),
		kong.Vars{"version": version},
	)
	cli.SetAccessToken(c.Token)
	err := ctx.Run(&cmd.Context{Token: c.Token})
	ctx.FatalIfErrorf(err)
	os.Exit(0)
}
