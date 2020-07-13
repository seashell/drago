package command

import (
	"context"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/version"
)

const (
	drago_addr_env	string	=	"DRAGO_ADDR"
	drago_token_env	string	=	"DRAGO_TOKEN"
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
			if configFile != "" {
				config := LoadConfigFromFile(configFile)
				ctx := context.WithValue(cmd.Context(), ctxKeyType("config"), config)
				cmd.SetContext(ctx)
			}
		},
	}

	cmd.SetVersionTemplate("Drago {{.Version}}")

	cmd.PersistentFlags().StringVarP(&configFile, "config", "", "", "config file path")

	cmd.AddCommand(NewAgentCmd())
	cmd.AddCommand(NewNetworksCmd())
	cmd.AddCommand(NewHostsCmd())
	cmd.AddCommand(NewInterfacesCmd())
	cmd.AddCommand(NewLinksCmd())
	cmd.AddCommand(NewTokensCmd())
	return cmd
}
