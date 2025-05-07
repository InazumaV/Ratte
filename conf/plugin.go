package conf

import (
	"fmt"
	"github.com/goccy/go-json"
)

type CorePlugin struct {
	DataPath string          `json:"DataPath,omitempty"`
	Config   json.RawMessage `json:"Config,omitempty"`
	Plugin
}
type _core CorePlugin

func (c *CorePlugin) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, (*_core)(c))
	if err != nil {
		return fmt.Errorf("failed to unmarshal core: %v", err)
	}
	if len(c.Config) == 0 {
		c.Config = data
	}

	err = json.Unmarshal(c.Config, &c.Plugin)
	if err != nil {
		return fmt.Errorf("failed to unmarshal core config: %v", err)
	}
	return nil
}

type Plugins struct {
	Core  []CorePlugin `json:"core"`
	Panel []Plugin     `json:"panel"`
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
