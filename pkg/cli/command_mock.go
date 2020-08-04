package cli

import "context"

// MockCommand is an implementation of Command that can be used for testing.
// It is also used for automatically populating missing parent commands.
type MockCommand struct {
	HelpText      string
	SynopsisText  string
	RunReturnCode int
	// Set by the command
	RunCalled bool
	RunArgs   []string
}

func (c *MockCommand) Help() string {
	return c.HelpText
}

func (c *MockCommand) Run(ctx context.Context, args []string) int {
	c.RunCalled = true
	c.RunArgs = args
	return c.RunReturnCode
}

func (c *MockCommand) Synopsis() string {
	return c.SynopsisText
}
