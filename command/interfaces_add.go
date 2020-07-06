package command

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/api"
)

func NewInterfacesAddCmd() *cobra.Command {
	// flags vars
	var name 		string
	var ipAddress 	string
	var listenPort 	string
	var hostID 		string
	var networkID 	string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a drago interface",
		Long: `Usage: drago interfaces add [options] [args]
		
		Add a interface. 
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

			ni := &api.InterfaceInput{
				Name: 		name,
				HostID:		hostID,
				NetworkID:	networkID,
				IPAddress:	ipAddress,
				ListenPort: listenPort,
			}

			nn,err := a.Interfaces().CreateInterface(ni); 
			if err != nil {
				fmt.Println("failed to add interface: ", err)
				os.Exit(1)
			}

			dump,_ := json.MarshalIndent(nn, "", "    ")
			fmt.Println(string(dump))
		},
	}

	// flags init
	cmd.Flags().StringVar(&name, "name", "", "Interface name")
	cmd.Flags().StringVar(&ipAddress, "ip-address", "", "Interface IP address in CIDR notation")
	cmd.Flags().StringVar(&listenPort, "listen-port", "", "Interface listen port")
	cmd.Flags().StringVar(&hostID, "host-id", "", "Interface target host ID")
	cmd.Flags().StringVar(&networkID, "network-id", "", "Interface target ntwork ID")

	// required flags
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("ip-address")
	cmd.MarkFlagRequired("host-id")
	cmd.MarkFlagRequired("network-id")

	return cmd
}
