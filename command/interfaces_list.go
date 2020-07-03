package command

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/api"
)

func NewInterfacesListCmd() *cobra.Command {
	// flags vars
	var hostID 		string
	var networkID 	string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List drago interfaces",
		Long: `Usage: drago interfaces list [options] [args]
		
		List all registered interfaces. 
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

			f := api.ListInterfacesInput{
				HostIDFilter: 		hostID,
				NetworkIDFilter:	networkID,
			}

			hl,err := a.Interfaces().ListInterfaces(f)
			if err != nil {
				fmt.Println("failed to list interfaces: ",err)
				os.Exit(1)
			}

			dump,_ := json.MarshalIndent(hl.Items, "", "    ")
			fmt.Println(string(dump))
		},
	}

	// flags init
	cmd.Flags().StringVar(&hostID, "host-id", "", "host ID")
	cmd.Flags().StringVar(&networkID, "network-id", "", "network ID")

	return cmd
}
