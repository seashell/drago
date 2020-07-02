package command

import (
	"context"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/version"
)

// Used for context
type ctxKeyType string

var (
	// Used for flags
	configFile string
)

func NewRootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "drago",
		Short:   "A flexible configuration manager for Wireguard networks",
		Long:    `Usage: drago [--version] [--help] [--autocomplete-(un)install] <command> [args]`,
		Version: version.GetVersion().VersionNumber(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config := LoadConfigFromFile(configFile)
			ctx := context.WithValue(cmd.Context(), ctxKeyType("config"), config)
			cmd.SetContext(ctx)
		},
	}

	cmd.SetVersionTemplate("Drago {{.Version}}")

	cmd.PersistentFlags().StringVarP(&configFile, "config", "", "/etc/drago.d", "config file (default is /etc/drago.d)")

	cmd.AddCommand(NewAgentCmd())

	return cmd
}
