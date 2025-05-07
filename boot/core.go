package boot

import (
	"fmt"
	"github.com/InazumaV/Ratte-Interface/core"
	"github.com/InazumaV/Ratte/conf"
	"os/exec"
)

type CoreGroup map[string]*core.PluginClient

func (c CoreGroup) Get(name string) *core.PluginClient {
	if co, ok := c[name]; ok {
		return co
	}
	return nil
}

func (c CoreGroup) Close() {
	for _, co := range c {
		co.Close()
	}
}

func initCores(cc []conf.CorePlugin) (CoreGroup, error) {
	cores := make(CoreGroup, len(cc))
	for _, co := range cc {
		c, err := core.NewClient(nil, exec.Command(co.Path))
		if err != nil {
			return nil, fmt.Errorf("new core error: %w", err)
		}
		err = c.Start(co.DataPath, co.Config)
		if err != nil {
			return nil, fmt.Errorf("start core error: %w", err)
		}
		cores[co.Name] = c
		break
	}
	return cores, nil
}
