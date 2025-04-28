package main

import (
	"fmt"
	"github.com/InazumaV/Ratte-Interface/core"
	"github.com/InazumaV/Ratte-Interface/panel"
	"github.com/InazumaV/Ratte/acme"
	"github.com/InazumaV/Ratte/conf"
	"github.com/InazumaV/Ratte/handler"
	"github.com/InazumaV/Ratte/trigger"
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
	cores    map[string]*core.PluginClient
	panels   map[string]*panel.PluginClient
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
	err = startCores(c.Plugin.Core, c.Core)
	if err != nil {
		log.WithError(err).Fatal("Init core plugin failed")
	}
	log.Info("Done.")

	log.Info("Init panel plugin...")
	err = startPanel(c.Plugin.Panel)
	if err != nil {
		log.WithError(err).Fatal("Init panel plugin failed")
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
	err = startTriggerAndHandler(c.Node)
	if err != nil {
		log.WithError(err).Fatal("Start trigger and handler failed")
	}
	log.Info("Started.")

	c.SetErrorHandler(func(err error) {
		log.WithFields(map[string]interface{}{
			"error":    err,
			"Services": "ConfWatcher",
		}).Error("")
	})
	c.SetEventHandler(func(event uint, target ...string) {
		l := log.WithFields(map[string]interface{}{
			"Service": "ConfWatcher",
		})
		l.Info("Config changed, restart...")
		err := closeTriggerAndHandler()
		if err != nil {
			log.WithError(err).Fatal("Close trigger and handler failed")
		}
		err = closeCore()
		if err != nil {
			log.WithError(err).Fatal("Close core failed")
		}
		err = closePanel()
		if err != nil {
			log.WithError(err).Fatal("Close panel failed")
		}
		err = startCores(c.Plugin.Core, c.Core)
		if err != nil {
			log.WithError(err).Fatal("Start core failed")
		}
		err = startPanel(c.Plugin.Panel)
		if err != nil {
			log.WithError(err).Fatal("Start panel failed")
		}
		err = startTriggerAndHandler(c.Node)
		if err != nil {
			log.WithError(err).Fatal("Start trigger and handler failed")
		}
		l.Info("Restart Done.")
	})
	err = c.Watch()
	if err != nil {
		log.WithError(err).Fatal("Watch config failed")
	}

	runtime.GC()
	sig := make(chan os.Signal, 0)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sig
	// clear
	log.Info("Shutting down...")
	err = closeTriggerAndHandler()
	if err != nil {
		log.WithError(err).Fatal("Close trigger and handler failed")
	}
	err = closeCore()
	if err != nil {
		log.WithError(err).Fatal("Close core failed")
	}
	err = closePanel()
	if err != nil {
		log.WithError(err).Fatal("Close panel failed")
	}
	log.Info("Done.")
}

func startCores(coresP []conf.Plugin, cc []conf.Core) error {
	// new cores
	cores = make(map[string]*core.PluginClient, len(coresP))
	for _, ccV := range cc {
		for _, co := range coresP {
			if co.Name != ccV.Type {
				continue
			}
			c, err := core.NewClient(nil, exec.Command(co.Path))
			if err != nil {
				return fmt.Errorf("new core error: %w", err)
			}
			err = c.Start(ccV.DataPath, ccV.Config)
			if err != nil {
				return fmt.Errorf("start core error: %w", err)
			}
			cores[co.Name] = c
			break
		}
	}
	return nil
}

func closeCore() error {
	for _, c := range cores {
		err := c.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func startPanel(panelsP []conf.Plugin) error {
	panels = make(map[string]*panel.PluginClient, len(panelsP))
	for _, p := range panelsP {
		pn, err := panel.NewClient(nil, exec.Command(p.Path))
		if err != nil {
			return fmt.Errorf("new panel error: %w", err)
		}
		panels[p.Name] = pn
	}
	return nil
}

func closePanel() error {
	for _, p := range panels {
		err := p.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func startTriggerAndHandler(c []conf.Node) error {
	triggers = make([]*trigger.Trigger, 0, len(c))
	handlers = make([]*handler.Handler, 0, len(c))
	for _, nd := range c {
		var co core.Core
		if c, ok := cores[nd.Options.Core]; ok {
			co = c
		} else {
			return fmt.Errorf("unknown core name: %s", nd.Options.Core)
		}
		var pl panel.Panel
		if p, ok := panels[nd.Options.Panel]; ok {
			pl = p
		} else {
			return fmt.Errorf("")
		}
		var ac *acme.Acme
		if len(acmes) != 0 {
			if a, ok := acmes[nd.Options.Acme]; ok {
				ac = a
			} else {
				return fmt.Errorf("unknown acme name: %s", nd.Options.Acme)
			}
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
			return fmt.Errorf("new trigger error: %w", err)
		}
		triggers = append(triggers, tr)
		err = tr.Start()
		if err != nil {
			return fmt.Errorf("start trigger error: %w", err)
		}
	}
	return nil
}

func closeTriggerAndHandler() error {
	for _, t := range triggers {
		err := t.Close()
		if err != nil {
			return fmt.Errorf("close trigger error: %w", err)
		}
	}
	for _, h := range handlers {
		err := h.Close()
		if err != nil {
			return fmt.Errorf("close handler error: %w", err)
		}
	}
	return nil
}
