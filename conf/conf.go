package conf

import (
	"fmt"
	trim "github.com/InazumaV/Ratte/common/json"
	"github.com/InazumaV/Ratte/common/watcher"
	"net/http"
	"os"

	"github.com/goccy/go-json"
)

type Conf struct {
	// internal fields
	path             string
	watcherHandle    EventHandler
	errorHandler     ErrorHandler
	configWatcher    watcher.Watcher
	coreDataWatchers map[int]watcher.Watcher

	// config fields
	Log     Log     `json:"Log,omitempty"`
	Watcher Watcher `json:"Watcher,omitempty"`
	Plugin  Plugins `json:"Plugin,omitempty"` // Only accept from file path
	Acme    []ACME  `json:"Acme,omitempty"`
	Node    []Node  `json:"Node,omitempty"`
}

func New(path string) *Conf {
	return &Conf{
		path: path,
		Watcher: Watcher{
			WatchLocalConfig:  true,
			WatchRemoteConfig: true,
		},
		Log:    newLog(),
		Acme:   make([]ACME, 0),
		Plugin: Plugins{},
		Node:   make([]Node, 0),
	}
}

func (c *Conf) Load(data []byte) error {
	if len(data) > 0 {
		err := json.Unmarshal(data, c)
		if err != nil {
			return fmt.Errorf("decode json error: %w", err)
		}
		return nil
	}
	if IsHttpUrl(c.path) {
		rsp, err := http.Get(c.path)
		if err != nil {
			return err
		}
		defer rsp.Body.Close()
		err = json.NewDecoder(trim.NewTrimNodeReader(rsp.Body)).Decode(&c)
		if err != nil {
			return fmt.Errorf("decode json error: %w", err)
		}
	} else {
		f, err := os.Open(c.path)
		if err != nil {
			return err
		}
		defer f.Close()
		err = json.NewDecoder(trim.NewTrimNodeReader(f)).Decode(&c)
		if err != nil {
			return fmt.Errorf("decode json error: %w", err)
		}
	}
	return nil
}
