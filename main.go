//go:generate yarn --cwd ./ui build

package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"

	command "github.com/seashell/drago/command"
	cli "github.com/seashell/drago/pkg/cli"
	version "github.com/seashell/drago/version"
)

//go:embed ui/build/*
var uiefs embed.FS

// Credits to https://github.com/pulumi/pulumi/blob/master/pkg/cmd/pulumi/main.go
func panicHandler() {
	if panicPayload := recover(); panicPayload != nil {

		stack := string(debug.Stack())
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "================================================================================")
		fmt.Fprintln(os.Stderr, "Drago has encountered a fatal error. This is a bug!")
		fmt.Fprintln(os.Stderr, "We would appreciate a report: https://github.com/seashell/drago/issues/")
		fmt.Fprintln(os.Stderr, "Please provide all of the below text in your report.")
		fmt.Fprintln(os.Stderr, "================================================================================")

		fmt.Fprintf(os.Stderr, "Drago Version:       %s\n", version.Version)
		fmt.Fprintf(os.Stderr, "Go Version:          %s\n", runtime.Version())
		fmt.Fprintf(os.Stderr, "Go Compiler:         %s\n", runtime.Compiler)
		fmt.Fprintf(os.Stderr, "Architecture:        %s\n", runtime.GOARCH)
		fmt.Fprintf(os.Stderr, "Operating System:    %s\n", runtime.GOOS)
		fmt.Fprintf(os.Stderr, "Panic:               %s\n\n", panicPayload)
		fmt.Fprintln(os.Stderr, stack)
	}
}

func init() {
	os.Setenv("TZ", "UTC")
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {

	cli := setupCLI()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)

	go func() {
		<-signalCh
		fmt.Fprintf(os.Stderr, "Received signal. Interrupting...\n")
		cancel()
	}()

	defer panicHandler()

	code, err := cli.Run(ctx, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return code
}

func setupCLI() *cli.CLI {

	ui := &cli.SimpleUI{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	subfs, err := fs.Sub(uiefs, "ui/build")
	if err != nil {
		panic(err)
	}

	uifs := http.FS(subfs)

	cli := cli.New(&cli.Config{
		Name: "drago",
		Commands: map[string]cli.Command{
			"agent":       &command.AgentCommand{UI: ui, StaticFS: uifs},
			"node":        &command.NodeCommand{UI: ui},
			"node status": &command.NodeStatusCommand{UI: ui},
		},
		Version: version.GetVersion().VersionNumber(),
	})

	return cli
}
