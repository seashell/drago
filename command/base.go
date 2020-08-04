package command

import (
	cli "github.com/seashell/drago/pkg/cli"
)

// TODO: move UI to the structs implementing cli.Command interface
type BaseCommand struct {
	UI cli.UI
}
