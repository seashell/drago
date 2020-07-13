package command

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/api"
)

func NewNetworksAddCmd() *cobra.Command {
	// flags vars
	var name 			string
	var ipAddressRange 	string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a drago network",
		Long: `Usage: drago networks add [options] [args]
		
		Add a network. 
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

			ni := &api.NetworkInput{
				Name: 			name,
				IPAddressRange:	ipAddressRange,
			}

			nn,err := a.Networks().CreateNetwork(ni); 
			if err != nil {
				fmt.Println("failed to add network: ", err)
				os.Exit(1)
			}
			
			dump,_ := json.MarshalIndent(nn, "", "    ")
			fmt.Println(string(dump))
		},
	}

	// flags init
	cmd.Flags().StringVar(&name, "name", "", "network name")
	cmd.Flags().StringVar(&ipAddressRange, "ip-address-range", "", "network ip address range in CIDR notation")

	// required flags
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("ip-address-range")
	
	return cmd
}
