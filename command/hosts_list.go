package command

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/api"
)

func NewHostsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List drago hosts",
		Long: `Usage: drago hosts list [options] [args]
		
		List all registered hosts. 
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

			hl,err := a.Hosts().ListHosts()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			dump,_ := json.MarshalIndent(hl.Items, "", "    ")
			fmt.Println(string(dump))
		},
	}
	return cmd
}
