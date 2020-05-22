//go:generate go generate github.com/seashell/drago/ui

package main

import (
	"fmt"
	"os"

	"github.com/seashell/drago/command"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	log.SetLevel(log.DebugLevel)

	rootCmd := command.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
