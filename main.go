//go:generate go generate github.com/seashell/drago/server

package main

import (
	"github.com/seashell/drago/command"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	log.SetLevel(log.DebugLevel)

	command.Execute()
}
