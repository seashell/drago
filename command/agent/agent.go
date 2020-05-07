package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/seashell/drago/agent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var Command = &cobra.Command{
	Use:   "agent",
	Short: "Runs a Drago agent",
	Long: `Usage: drago agent [options]

	Starts the Drago agent and runs until an interrupt is received.
	The agent may be a client and/or server.
  
	The Drago agent's configuration primarily comes from the config
	files used, but a subset of the options may also be passed directly
	as CLI arguments.
  `,
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile == "" {
			fmt.Println("==> No config file specified. Using default values.")
		}

		fmt.Println("==> Starting drago agent...")

		config := agent.AgentConfig{}
		viper.Unmarshal(&config)

		agent, err := agent.New(config)
		if err != nil {
			fmt.Printf("==> Error starting Drago agent: %s\n", err)
		}

		agent.Run()

		var wait time.Duration

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
	Command.Flags().BoolP("server", "s", false, "Start agent in server mode")
	Command.Flags().BoolP("client", "c", true, "Start agent in client mode")
	Command.Flags().BoolP("dev", "d", false, "Start agent in dev mode")

	viper.BindPFlag("server", Command.Flags().Lookup("server"))
	viper.BindPFlag("client", Command.Flags().Lookup("client"))
	viper.BindPFlag("dev", Command.Flags().Lookup("dev"))
}
