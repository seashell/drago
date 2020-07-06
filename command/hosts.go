package command

import (
	"github.com/seashell/cobra"
)

func NewHostsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hosts",
		Short: "Interfacts with drago hosts",
		Long: `Usage: drago hosts <subcomand> [options] [args]
		
		This command groups subcommands for interacting with hosts. 
	  `,
	}

	cmd.AddCommand(NewHostsListCmd())
	cmd.AddCommand(NewHostsAddCmd())
	return cmd
}
