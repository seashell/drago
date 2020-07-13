package command

import (
	"github.com/seashell/cobra"
)

func NewNetworksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "networks",
		Short: "Interfacts with drago networks",
		Long: `This command groups subcommands for interacting with networks. 
	  `,
	}

	cmd.AddCommand(NewNetworksListCmd())
	cmd.AddCommand(NewNetworksAddCmd())

	return cmd
}
