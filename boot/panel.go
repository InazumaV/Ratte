package boot

import (
	"fmt"
	"github.com/InazumaV/Ratte-Interface/panel"
	"github.com/InazumaV/Ratte/conf"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

type PanelGroup map[string]*panel.PluginClient

func (p *PanelGroup) Get(name string) *panel.PluginClient {
	if pl, ok := (*p)[name]; ok {
		return pl
	}
	return nil
}

func (p *PanelGroup) Close() error {
	for _, pl := range *p {
		err := pl.Close()
		if err != nil {
			log.WithError(err).Warn()
		}
	}
	return nil
}

func initPanels(panelsP []conf.Plugin) (PanelGroup, error) {
	panels := make(PanelGroup, len(panelsP))
	for _, p := range panelsP {
		pn, err := panel.NewClient(nil, exec.Command(p.Path))
		if err != nil {
			return nil, fmt.Errorf("new panel error: %w", err)
		}
		panels[p.Name] = pn
	}
	return panels, nil
}
