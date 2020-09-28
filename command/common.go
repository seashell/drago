package command

import (
	"flag"
	"fmt"

	cli "github.com/seashell/drago/pkg/cli"
)

// manyStrings
type manyStrings []string

func (s *manyStrings) Set(val string) error {
	*s = append(*s, val)
	return nil
}

func (s *manyStrings) String() string {
	var out []string
	out = *s
	return fmt.Sprintf("%v", out)
}

// RootFlagSet :
type RootFlagSet struct {
	*flag.FlagSet
	envPaths    manyStrings
	configPaths manyStrings
}

// FlagSet declares flags that are common to all commands,
// returning a RootFlagSet struct that will hold their values after
// flag.Parse() is called by the command
func FlagSet(name string) *RootFlagSet {

	flags := &RootFlagSet{
		FlagSet: flag.NewFlagSet(name, flag.ContinueOnError),
	}

	flags.Usage = func() {}

	flags.Var(&flags.envPaths, "env", "")
	flags.Var(&flags.configPaths, "config", "")

	// TODO: direct output to UI
	flags.SetOutput(nil)

	return flags
}

// GlobalOptions returns the global usage options string.
func GlobalOptions() string {
	text := `
  --config=<path>
    Path to a HCL file containing valid Drago configurations.
    Overrides the DRAGO_CONFIG_PATH environment variable if set.

  --env=<path>
    Path to a an env file
    Overrides the DRAGO_ENV_FILE_PATH environment variable if set.
`
	return text
}

// DefaultErrorMessage returns the default error message for this command
func DefaultErrorMessage(cmd cli.NamedCommand) string {
	return fmt.Sprintf("For additional help try 'drago %s --help'", cmd.Name())
}
