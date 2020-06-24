//go:generate go generate github.com/seashell/drago/ui

package main

import (
	"context"
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

	ctx := context.TODO()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
