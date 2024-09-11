package main

import (
	"Ratte/acme"
	"Ratte/conf"
	"Ratte/handler"
	"Ratte/trigger"
	"github.com/Yuzuki616/Ratte-Interface/core"
	"github.com/Yuzuki616/Ratte-Interface/panel"
	"github.com/Yuzuki616/Ratte-Interface/plugin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
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

var (
	cores    map[string]*plugin.Client[core.Core]
	panels   map[string]*plugin.Client[panel.Panel]
	acmes    map[string]*acme.Acme
	handlers []*handler.Handler
	triggers []*trigger.Trigger
)

func runHandle(_ *cobra.Command, _ []string) {
	c := conf.New(config)
	log.WithField("path", config).Info("Load config...")
	err := c.Load(nil)
	if err != nil {
		log.WithError(err).Fatal("Load config failed")
	}
	log.WithField("path", config).Info("Loaded.")

	log.Info("Init core plugin...")
	// new cores
	cores = make(map[string]*plugin.Client[core.Core], len(c.Core))
	for _, co := range c.Core {
		c, err := plugin.NewClient[core.Core](&plugin.Config{
			Type: plugin.CoreType,
			Cmd:  exec.Command(co.Path),
		})
		if err != nil {
			log.WithError(err).WithField("core", co.Name).Fatal("New core failed")
		}
		err = c.Caller().Start(co.DataPath, co.Config)
		if err != nil {
			log.WithError(err).WithField("core", co.Name).Fatal("Start core failed")
		}
		cores[co.Name] = c
	}
	log.Info("Done.")

	log.Info("Init panel plugin...")
	// new panels
	panels = make(map[string]*plugin.Client[panel.Panel], len(c.Panel))
	for _, p := range c.Panel {
		pn, err := plugin.NewClient[panel.Panel](&plugin.Config{
			Type: plugin.PanelType,
			Cmd:  exec.Command(p.Path),
		})
		if err != nil {
			log.WithError(err).WithField("panel", p.Name).Fatal("New panel failed")
		}
		panels[p.Name] = pn
	}
	log.Info("Done.")

	log.Info("Init acme...")
	// new acme
	acmes = make(map[string]*acme.Acme)
	for _, a := range c.Acme {
		ac, err := acme.NewAcme(&a)
		if err != nil {
			log.WithError(err).Fatal("New acme failed")
		}
		acmes[a.Name] = ac
	}
	log.Info("Done.")

	log.Info("Starting...")
	// new node
	triggers = make([]*trigger.Trigger, 0, len(c.Node))
	handlers = make([]*handler.Handler, len(c.Node))
	for _, nd := range c.Node {
		var co core.Core
		var pl panel.Panel
		var ac *acme.Acme
		if c, ok := cores[nd.Options.Core]; ok {
			co = c.Caller()
		} else {
			log.WithField("core", nd.Options.Core).Fatal("Couldn't find core")
		}
		if p, ok := panels[nd.Options.Panel]; ok {
			pl = p.Caller()
		} else {
			log.WithField("panel", nd.Options.Panel).Fatal("Couldn't find panel")
		}
		if a, ok := acmes[nd.Options.Acme]; ok {
			ac = a
		} else {
			log.WithField("acme", nd.Options.Acme).Fatal("Couldn't find acme")
		}

		h := handler.New(co, pl, nd.Name, ac, log.WithFields(
			map[string]interface{}{
				"node":    nd.Name,
				"service": "handler",
			},
		), &nd.Options)
		handlers = append(handlers, h)
		tr, err := trigger.NewTrigger(log.WithFields(
			map[string]interface{}{
				"node":    nd.Name,
				"service": "trigger",
			},
		), &nd.Trigger, h, pl, &nd.Remote)
		if err != nil {
			log.WithError(err).Fatal("New trigger failed")
		}
		triggers = append(triggers, tr)
		err = tr.Start()
		if err != nil {
			log.WithError(err).Fatal("Start trigger failed")
		}
	}
	log.Info("Started.")

	runtime.GC()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sig
}
