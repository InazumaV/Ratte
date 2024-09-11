package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pre "github.com/x-cray/logrus-prefixed-formatter"
)

var version string
var buildDate string

var command = &cobra.Command{
	Use: "Ratte",
}

func Execute() {
	err := command.Execute()
	if err != nil {
		log.WithField("err", err).Error("Execute command failed")
	}
}

func main() {
	log.SetFormatter(&pre.TextFormatter{
		TimestampFormat: "01-02 15:04:05",
		FullTimestamp:   true,
	})
	log.Info("Ratte")
	log.Info("Version: ", version)
	log.Info("Build date: ", buildDate)
	Execute()
}
