package conf

import (
	"Ratte/common/watcher"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path"
)

const (
	ConfigFileChangedEvent   = 0
	CoreDataPathChangedEvent = 1
)

type EventHandler func(event uint, target ...string)
type ErrorHandler func(err error)

type Watcher struct {
	WatchLocalConfig  bool `json:"WatchLocalConfig,omitempty"`
	WatchRemoteConfig bool `json:"WatchRemoteConfig,omitempty"`
	WatchCoreDataPath bool `json:"WatchCoreDataPath,omitempty"`
	RemoteInterval    uint `json:"Interval,omitempty"`
}

func (c *Conf) SetEventHandler(w EventHandler) {
	c.watcherHandle = w
}

func (c *Conf) SetErrorHandler(w ErrorHandler) {
	c.errorHandler = w
}

func (c *Conf) Watch() error {
	if c.watcherHandle == nil {
		return errors.New("no watch handler")
	}
	if c.errorHandler == nil {
		c.errorHandler = func(err error) {
			log.WithField("service", "conf_watcher").Error(err)
		}
	}
	if IsHttpUrl(c.path) {
		if c.Watcher.WatchRemoteConfig {
			w := watcher.NewHTTPWatcher(c.path, c.Watcher.RemoteInterval)
			c.configWatcher = w
		}
	} else {
		if !c.Watcher.WatchLocalConfig {
			w := watcher.NewLocalWatcher(path.Dir(c.path), []string{path.Base(c.path)})
			c.configWatcher = w
		}
	}
	if c.Watcher.WatchLocalConfig || c.Watcher.WatchRemoteConfig {
		c.configWatcher.SetErrorHandler(watcher.ErrorHandler(c.errorHandler))
		c.configWatcher.SetEventHandler(func(_ string) error {
			c.watcherHandle(ConfigFileChangedEvent)
			return nil
		})
		err := c.configWatcher.Watch()
		if err != nil {
			return fmt.Errorf("watch config err:%w", err)
		}
	}
	if !c.Watcher.WatchCoreDataPath {
		return nil
	}

	watchers := make(map[int]*watcher.LocalWatcher, len(c.Core))
	for i, co := range c.Core {
		w := watcher.NewLocalWatcher(co.DataPath, []string{"*"})
		w.SetErrorHandler(watcher.ErrorHandler(c.errorHandler))
		w.SetEventHandler(func(_ string) error {
			c.watcherHandle(CoreDataPathChangedEvent, c.Core[i].Name)
			return nil
		})
		err := w.Watch()
		if err != nil {
			return fmt.Errorf("watch core %s err:%w", co.Name, err)
		}
		watchers[i] = w
	}
	return nil
}
