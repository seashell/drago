package cli

import "context"

const (
	CommandReturnCodeHelp = -18511
)

// A Command is a runnable sub-command of a CLI
type Command interface {
	// Help should return long-form help text that includes the command-line
	// usage, a brief few sentences explaining the function of the command,
	// and the complete list of flags the command accepts.
	Help() string

	// Synopsis should return a one-line, short synopsis of the command.
	// This should be less than 50 characters ideally.
	Synopsis() string

	// Run should run the actual command with the given CLI instance and
	// command-line arguments. It should return the exit status when it is
	// finished.
	Run(ctx context.Context, args []string) int
}

// A NamedCommand is a runnable sub-command of a CLI with a name
type NamedCommand interface {
	Command
	Name() string
}
