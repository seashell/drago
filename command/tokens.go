package command

import (
	"github.com/seashell/cobra"
)

func NewTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Interfacts with drago tokens",
		Long: `Usage: drago tokens <subcomand> [options] [args]
		
		This command groups subcommands for interacting with tokens. 
	  `,
	}
	cmd.AddCommand(NewTokensAddCmd())
	return cmd
}
