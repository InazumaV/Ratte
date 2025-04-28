package conf

import (
	"fmt"
	"github.com/goccy/go-json"
)

type Core struct {
	Type     string          `json:"type"`
	DataPath string          `json:"DataPath,omitempty"`
	Config   json.RawMessage `json:"Config,omitempty"`
}
type _core Core

func (c *Core) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, (*_core)(c))
	if err != nil {
		return fmt.Errorf("failed to unmarshal core: %v", err)
	}
	if len(c.Config) == 0 {
		c.Config = data
	}
	return nil
}

type Plugins struct {
	Core  []Plugin `json:"core"`
	Panel []Plugin `json:"panel"`
}

type Plugin struct {
	Name string `json:"Name,omitempty"`
	Path string `json:"Path,omitempty"`
}

type _plugins Plugins

func (p *Plugins) UnmarshalJSON(data []byte) error {
	var path string
	err := json.Unmarshal(data, path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, (*_plugins)(p))
	if err != nil {
		return fmt.Errorf("failed to unmarshal plugin: %v", err)
	}
	return nil
}
