package boot

import (
	"fmt"
	"github.com/InazumaV/Ratte-Interface/core"
	"github.com/InazumaV/Ratte-Interface/panel"
	"github.com/InazumaV/Ratte/conf"
	"github.com/InazumaV/Ratte/handler"
	"github.com/InazumaV/Ratte/trigger"
	log "github.com/sirupsen/logrus"
)

type NodeGroup struct {
	t []*trigger.Trigger
	h []*handler.Handler
}

func (n NodeGroup) Start() error {
	for _, t := range n.t {
		err := t.Start()
		if err != nil {
			return fmt.Errorf("start trigger error: %w", err)
		}
	}
	return nil
}

func (n NodeGroup) Close() error {
	for _, t := range n.t {
		err := t.Close()
		if err != nil {
			log.WithError(err).Errorln("Close trigger error")
		}
	}
	for _, h := range n.h {
		err := h.Close()
		if err != nil {
			log.WithError(err).Errorln("Close handler error")
		}
	}
	return nil
}

func initNode(
	n []conf.Node,
	acme AcmeGroup,
	cores CoreGroup,
	panels PanelGroup) (*NodeGroup, error) {
	triggers := make([]*trigger.Trigger, 0, len(n))
	handlers := make([]*handler.Handler, 0, len(n))
	for _, nd := range n {
		var co core.Core
		if c, ok := cores[nd.Options.Core]; ok {
			co = c
		} else {
			return nil, fmt.Errorf("unknown core name: %s", nd.Options.Core)
		}
		var pl panel.Panel
		if p, ok := panels[nd.Options.Panel]; ok {
			pl = p
		} else {
			return nil, fmt.Errorf("")
		}
		ac, e := acme[nd.Options.Acme]
		if !e {
			return nil, fmt.Errorf("unknown acme name: %s", nd.Options.Acme)
		}
		h := handler.New(co, pl, nd.Name, ac, log.WithFields(
			map[string]interface{}{
				"node":    nd.Name,
				"service": "handler",
			},
		), &nd.Options)
		handlers = append(handlers, h)
		tr, err := trigger.New(log.WithFields(
			map[string]interface{}{
				"node":    nd.Name,
				"service": "trigger",
			},
		), &nd.Trigger, h, pl, &nd.Remote)
		if err != nil {
			return nil, fmt.Errorf("new trigger error: %w", err)
		}
		triggers = append(triggers, tr)
		err = tr.Start()
		if err != nil {
			return nil, fmt.Errorf("start trigger error: %w", err)
		}
	}
	return &NodeGroup{
			t: triggers,
			h: handlers,
		},
		nil
}
