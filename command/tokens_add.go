package command

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/seashell/cobra"
	"github.com/seashell/drago/api"
)

func NewTokensAddCmd() *cobra.Command {
	// flags vars
	var hostID	string
	var labels *[]string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a drago token",
		Long: `Usage: drago tokens add [options] [args]
		
		Add a token. 
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

			ni := &api.TokenInput{
				Type: 		"client",
				Subject:	hostID,
				Labels:		*labels,
			}

			nn,err := a.Tokens().CreateToken(ni); 
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			
			dump,_ := json.MarshalIndent(nn, "", "    ")
			fmt.Println(string(dump))
		},
	}

	// flags init
	cmd.Flags().StringVar(&hostID, "host-id", "", "token host ID")
	labels = cmd.Flags().StringSlice("label", nil, "token label")
	// required flags
	cmd.MarkFlagRequired("host-id")
	
	return cmd
}
