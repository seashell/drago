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
	"fmt"

	log "github.com/sirupsen/logrus"

	consul "github.com/edufschmidt/dragonair/pkg/backends/consul"

	"github.com/spf13/cobra"
)

type AgentConfig struct {
	consulAddr string
}

func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{}
}

type Agent struct {
	config *AgentConfig

	consulClient *consul.ConsulClient

	// 	shutdown     bool
	// 	shutdownCh   chan struct{}
	// 	shutdownLock sync.Mutex
}

func NewAgent(config *AgentConfig) (*Agent, error) {
	return &Agent{
		config: config,
	}, nil
}

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Runs a dragonair agent",
	Long: `
	Usage: dragonair agent [options]
	
	Starts the dragonair agent and runs until an interrupt is received.
	The agent may be a client and/or server.
  
	The dragonair agent's configuration primarily comes from the config
	files used, but a subset of the options may also be passed directly
	as CLI arguments.
  `,
	Run: func(cmd *cobra.Command, args []string) {

		configFile, _ := cmd.Flags().GetString("config")

		fmt.Println("==> Loaded configuration from ", configFile)

		fmt.Println("==> Starting dragonair agent...")
		fmt.Println("==> dragonair agent configuration:")

		fmt.Println("")
		fmt.Println("       Server ", false)
		fmt.Println("")

		fmt.Println("==> dragonair agent started! Log data will stream in below:")
		fmt.Println("")

		log.Info("Log data")

		//node := node.New(nil)
		//node.Run()

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

	agentCmd.Flags().String("config", "", "Path to custom config file.")

	rootCmd.AddCommand(agentCmd)
}
