package command

import (
	"github.com/seashell/cobra"
)

func NewInterfacesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interfaces",
		Short: "Interfacts with drago interfaces",
		Long: `Usage: drago interfaces <subcomand> [options] [args]
		
		This command groups subcommands for interacting with interfaces. 
	  `,
	}

	cmd.AddCommand(NewInterfacesListCmd())
	cmd.AddCommand(NewInterfacesAddCmd())
	return cmd
}
