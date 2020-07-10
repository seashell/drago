package command

import (
	"github.com/seashell/cobra"
)

func NewLinksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links",
		Short: "Interfacts with drago links",
		Long: `Usage: drago links <subcomand> [options] [args]
		
		This command groups subcommands for interacting with links. 
	  `,
	}

	cmd.AddCommand(NewLinksListCmd())
	cmd.AddCommand(NewLinksAddCmd())
	return cmd
}
