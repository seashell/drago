package command

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/api"
)

func NewHostsAddCmd() *cobra.Command {
	// flags vars
	var name 				string
	var advertiseAddress 	string
	var labels				*[]string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a drago host",
		Long: `Usage: drago hosts add [options] [args]
		
		Add a host. 
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

			ni := &api.HostInput{
				Name: 				name,
				Labels:				*labels,
				AdvertiseAddress:	advertiseAddress,
			}

			nn,err := a.Hosts().CreateHost(ni); 
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			dump,_ := json.MarshalIndent(nn, "", "    ")
			fmt.Println(string(dump))
		},
	}

		// flags init
		cmd.Flags().StringVar(&name, "name", "", "host name")
		cmd.Flags().StringVar(&advertiseAddress, "advertise-address", "", "host advertise IP address")
		labels = cmd.Flags().StringSlice("label", nil, "host labels")
	
		// required flags
		cmd.MarkFlagRequired("name")

	return cmd
}
