package command

import (
	"flag"
	"fmt"
	"strings"

	cli "github.com/seashell/drago/pkg/cli"
)

type RootFlagSet struct {
	*flag.FlagSet
	// Attributes containing values parsed from user input (e.g., flags,
	// environment variables, etc) which are not directly required by the
	// command implementation, but are important to provide them with common
	// functionality such as the API client, and avoid replicating code.
	//
	// Address of the Drago server
	address string

	// Secret token for authentication
	token string
}

// FlagSet declares flags that are common to all commands,
// returning a RootFlagSet struct that will hold their values after
// flag.Parse() is called by the command
func FlagSet(name string) *RootFlagSet {

	flags := &RootFlagSet{
		FlagSet: flag.NewFlagSet(name, flag.ContinueOnError),
	}

	flags.Usage = func() {}

	flags.StringVar(&flags.address, "address", "", "")
	flags.StringVar(&flags.token, "token", "", "")

	// TODO: direct output to UI
	flags.SetOutput(nil)

	return flags
}

// GlobalOptions returns the global usage options string.
func GlobalOptions() string {
	text := `
  --address=<addr>
    The address of the Drago server.
    Overrides the DRAGO_ADDR environment variable if set.
    Default = http://127.0.0.1:8080

  --token
    The SecretID of an ACL token to use to authenticate API requests with.
    Overrides the DRAGO_TOKEN environment variable if set.
`
	return strings.TrimSpace(text)
}

// DefaultErrorMessage returns the default error message for this command
func DefaultErrorMessage(cmd cli.NamedCommand) string {
	return fmt.Sprintf("For additional help try 'drago %s --help'", cmd.Name())
}
