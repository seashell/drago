package command

import (
	"github.com/seashell/cobra"
)

var _configFile string

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drago",
		Short: "A flexible configuration manager for Wireguard networks",
		Long:  `Usage: drago [-version] [-help] [-autocomplete-(un)install] <command> [args]`,
	}

	cmd.PersistentFlags().StringVarP(&_configFile, "config", "", "/etc/drago.d", "config file (default is /etc/drago.d)")

	cmd.AddCommand(NewAgentCmd())

	return cmd
}
