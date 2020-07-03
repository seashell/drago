package command

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/api"
)

func NewLinksAddCmd() *cobra.Command {
	// flags vars
	var fromInterfaceID 	string
	var toInterfaceID 		string
	var allowedIPs			*[]string
	var persistentKeepalive	int
	
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a drago link",
		Long: `Usage: drago links add [options] [args]
		
		Add a link. 
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

			ni := &api.LinkInput{
				FromInterfaceID: 		fromInterfaceID,
				ToInterfaceID:			toInterfaceID,
				AllowedIPs:				*allowedIPs,
				PersistentKeepalive:	persistentKeepalive,
			}

			nn,err := a.Links().CreateLink(ni); 
			if err != nil {
				fmt.Println("failed to add link: ", err)
				os.Exit(1)
			}

			dump,_ := json.MarshalIndent(nn, "", "    ")
			fmt.Println(string(dump))
		},
	}

	// flags init
	cmd.Flags().StringVar(&fromInterfaceID, "source-interface-id", "", "Link source interface ID")
	cmd.Flags().StringVar(&toInterfaceID, "target-interface-id", "", "Link target interface ID")
	allowedIPs = cmd.Flags().StringSlice("allowed-ip", nil, "Link allowed IPs")
	cmd.Flags().IntVar(&persistentKeepalive, "persistante-keepalive", 0, "Link persistante keepalive settings")

	// required flags
	cmd.MarkFlagRequired("source-interface-id")
	cmd.MarkFlagRequired("target-interface-id")
	cmd.MarkFlagRequired("allowed-ip")

	return cmd
}
