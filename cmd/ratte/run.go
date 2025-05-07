package main

import (
	"github.com/InazumaV/Ratte/boot"
	"github.com/InazumaV/Ratte/conf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var runCommand = cobra.Command{
	Use:   "server",
	Short: "Run Ratte",
	Run:   runHandle,
	Args:  cobra.NoArgs,
}
var config string

func init() {
	runCommand.PersistentFlags().
		StringVarP(&config, "config", "c",
			"./config.json5", "config file path")
	command.AddCommand(&runCommand)
}

func runHandle(_ *cobra.Command, _ []string) {
	c := conf.New(config)
	log.WithField("path", config).Info("Load config...")
	err := c.Load(nil)
	if err != nil {
		log.WithError(err).Fatal("Load config failed")
	}
	log.WithField("path", config).Info("Loaded.")

	log.Info("Init...")
	b, err := boot.InitBoot(c)
	if err != nil {
		log.WithError(err).Fatal("Init failed")
	}
	log.Info("Init done.")

	log.Info("Start...")
	err = b.Start()
	if err != nil {
		log.WithError(err).Fatal("Start failed")
	}
	c.SetErrorHandler(func(err error) {
		log.WithFields(map[string]interface{}{
			"error":    err,
			"Services": "ConfWatcher",
		}).Error("watch config error")
	})
	c.SetEventHandler(func(event uint, target ...string) {
		l := log.WithFields(map[string]interface{}{
			"Service": "ConfWatcher",
		})
		l.Info("Config changed, restart...")
		b, err := boot.InitBoot(c)
		if err != nil {
			l.WithError(err).Fatal("Init boot failed")
		}
		err = b.Start()
		if err != nil {
			l.WithError(err).Fatal("Start boot failed")
		}
		l.Info("Restart Done.")
	})
	err = c.Watch()
	if err != nil {
		log.WithError(err).Fatal("Watch config failed")
	}
	log.Info("Start done.")

	runtime.GC()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sig
	// clear
	log.Info("Shutting down...")
	b.Close()
	log.Info("Done.")
}
