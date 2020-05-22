package command

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drago",
		Short: "A flexible configuration manager for Wireguard networks",
		Long:  `Usage: drago [-version] [-help] [-autocomplete-(un)install] <command> [args]`,
	}

	var config string
	cmd.PersistentFlags().StringVarP(&config, "config", "", "/etc/drago.d", "config file (default is /etc/drago.d)")

	cmd.AddCommand(NewAgentCmd())

	return cmd
}
