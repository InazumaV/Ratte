package conf

import (
	"fmt"
	"github.com/goccy/go-json"
)

type Core struct {
	Name     string          `json:"Name,omitempty"`
	Path     string          `json:"Path,omitempty"`
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
