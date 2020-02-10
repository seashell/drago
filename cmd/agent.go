/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"time"

	agent "github.com/seashell/drago/pkg/agent"
	version "github.com/seashell/drago/pkg/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Runs a drago agent",
	Long: `
	Usage: drago agent [options]
	
	Starts the drago agent and runs until an interrupt is received.
	The agent may be a client and/or server.
  
	The drago agent's configuration primarily comes from the config
	files used, but a subset of the options may also be passed directly
	as CLI arguments.
  `,
	Run: func(cmd *cobra.Command, args []string) {

		if cfgFile == "" {
			fmt.Println("==> No config file specified. Using default values.")
		}

		fmt.Println("==> Starting drago agent...")
		fmt.Println("==> drago agent configuration:")

		fmt.Println("")
		fmt.Println("	Interface: ", viper.GetString("iface"))
		fmt.Println("	Address: ", viper.GetString("network"))
		fmt.Println("	Server: ", viper.GetBool("server"))
		fmt.Println("	Web UI: ", viper.GetBool("ui"))
		fmt.Println("	Version: ", version.GetVersion().VersionNumber())
		fmt.Println("")

		if viper.GetBool("ui") {
			fmt.Println("	UI: http://localhost:8080")
		}

		fmt.Println("")
		fmt.Println("==> drago agent started! Log data will stream in below:")
		fmt.Println("")

		var config agent.AgentConfig

		viper.Unmarshal(&config)

		a, err := agent.NewAgent(config)
		if err != nil {
			panic("Error creating agent")
		}

		var wait time.Duration

		a.Run()

		c := make(chan os.Signal, 1)

		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
		// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
		signal.Notify(c, os.Interrupt)

		// Block until we receive our signal.
		<-c

		// Create a deadline to wait for.
		_, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()

		os.Exit(0)
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// agentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// agentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Declare flags
	agentCmd.Flags().BoolP("server", "s", false, "Start agent in server mode")
	agentCmd.Flags().BoolP("client", "c", true, "Start agent in client mode")
	agentCmd.Flags().Bool("ui", true, "Serve web UI for configuration")

	// Set default values for configs not exposed through flags
	viper.SetDefault("iface", "wg0")
	viper.SetDefault("network", "192.168.2.0/24")

	viper.SetDefault("bind_addr", "127.0.0.1")
	viper.SetDefault("server_addr", "192.168.2.1/24")
	viper.SetDefault("client.data_dir", "/tmp/drago")

	// Bind viper configs to cobra flags
	viper.BindPFlag("server", agentCmd.Flags().Lookup("server"))
	viper.BindPFlag("client", agentCmd.Flags().Lookup("client"))
	viper.BindPFlag("ui", agentCmd.Flags().Lookup("ui"))

	rootCmd.AddCommand(agentCmd)
}
