package command

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/api"
)

func NewNetworksListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List drago networks",
		Long: `Usage: drago networks list [options] [args]
		
		List all registered networks. 
	  `,
	  	Run: func(cmd *cobra.Command, args []string) {
			//create api instance
			serverAddr := os.Getenv(drago_addr_env)
			serverToken := os.Getenv(drago_token_env)
			a,err := api.NewClient(&api.Config{
				Address:	serverAddr,
				Token:		serverToken,
			})
			if err != nil {
				fmt.Println("failed to initialize API: ",err)
				os.Exit(1)
			}

			nl,err := a.Networks().ListNetworks()
			if err != nil {
				fmt.Println("failed to list networks: ",err)
				os.Exit(1)
			}

			dump,_ := json.MarshalIndent(nl.Items, "", "    ")
			fmt.Println(string(dump))
		},
	}
	return cmd
}
