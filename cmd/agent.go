package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sort"
	"time"

	"github.com/seashell/drago/agent"
	version "github.com/seashell/drago/version"

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

		var config agent.AgentConfig
		viper.Unmarshal(&config)


		// Create configuration info structure
		info := make(map[string]string)
		info["client"] = viper.GetString("client.enabled")
		info["server"] = viper.GetString("server.enabled")
		info["interface"] = viper.GetString("client.iface")
		info["data dir"] = viper.GetString("client.data_dir")
		if viper.GetBool("ui") {
			info["web ui"] = "http://localhost:8080"
		} else {
			info["web ui"] = "false"
		}
		info["version"] =  version.GetVersion().VersionNumber()	

		// Sort the keys for output
		infoKeys := make([]string, 0, len(info))
		for key := range info {
			infoKeys = append(infoKeys, key)
		}
		sort.Strings(infoKeys)

		// Agent configuration output
		padding := 18
		fmt.Println("==> Drago agent configuration:")
		fmt.Println("")

		for _, k := range infoKeys {
			fmt.Println(fmt.Sprintf(
				"%s%s: %s",
				strings.Repeat(" ", padding-len(k)),
				strings.Title(k),
				info[k]))
		}
		fmt.Println("")

		a, err := agent.NewAgent(config)
		if err != nil {
			panic("Error initializing agent")
		}

		fmt.Println("")
		fmt.Println("==> drago agent started! Log data will stream in below:")
		fmt.Println("")

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
	//viper.SetDefault("iface", "wg0")
	viper.SetDefault("network", "192.168.2.0/24")

	viper.SetDefault("bind_addr", "127.0.0.1")
	viper.SetDefault("server_addr", "192.168.2.1/24")

	// Bind viper configs to cobra flags
	viper.BindPFlag("server", agentCmd.Flags().Lookup("server"))
	viper.BindPFlag("client", agentCmd.Flags().Lookup("client"))
	viper.BindPFlag("ui", agentCmd.Flags().Lookup("ui"))

	rootCmd.AddCommand(agentCmd)
}
