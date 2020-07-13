package command

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/api"
)

func NewLinksListCmd() *cobra.Command {
	// flags vars
	var networkID 	string
	var hostID      string 
	var interfaceID	string 

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List drago links",
		Long: `Usage: drago links list [options] [args]
		
		List all registered links. 
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

			f := api.ListLinksInput{
				SourceHostIDFilter:			hostID,
				NetworkIDFilter:			networkID,
				SourceInterfaceIDFilter:	interfaceID,
			}

			hl,err := a.Links().ListLinks(f)
			if err != nil {
				fmt.Println("failed to list links: ",err)
				os.Exit(1)
			}

			dump,_ := json.MarshalIndent(hl.Items, "", "    ")
			fmt.Println(string(dump))
		},
	}

	// flags init
	cmd.Flags().StringVar(&hostID, "host-id", "", "Link host ID")
	cmd.Flags().StringVar(&networkID, "network-id", "", "Link network ID")
	cmd.Flags().StringVar(&networkID, "interface-id", "", "Link interface ID")

	return cmd
}
