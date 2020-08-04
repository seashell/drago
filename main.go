//go:generate go generate github.com/seashell/drago/ui

package main

import (
	"context"
	"fmt"
	"os"

	command "github.com/seashell/drago/command"
	cli "github.com/seashell/drago/pkg/cli"
	version "github.com/seashell/drago/version"
)

func init() {
	os.Setenv("TZ", "UTC")
}

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {

	cli := setupCLI()

	code, err := cli.Run(context.Background(), args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return code
}

// setupCLI
func setupCLI() *cli.CLI {

	base := command.BaseCommand{
		UI: &cli.SimpleUI{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	cli := cli.New(&cli.Config{
		Name:    "drago",
		Version: version.GetVersion().VersionNumber(),
		Commands: map[string]cli.Command{
			"agent": &command.AgentCommand{
				BaseCommand: base,
			},
			"host": &command.HostCommand{
				BaseCommand: base,
			},
			"host create": &command.HostCreateCommand{
				BaseCommand: base,
			},
			"host delete": &command.HostDeleteCommand{
				BaseCommand: base,
			},
			"host list": &command.HostListCommand{
				BaseCommand: base,
			},
			"token": &command.TokenCommand{
				BaseCommand: base,
			},
			"token create": &command.TokenCreateCommand{
				BaseCommand: base,
			},
			"token delete": &command.TokenDeleteCommand{
				BaseCommand: base,
			},
			"token info": &command.TokenInfoCommand{
				BaseCommand: base,
			},
			"token list": &command.TokenListCommand{
				BaseCommand: base,
			},
		},
	})

	return cli
}
