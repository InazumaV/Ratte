package watcher

import (
	"fmt"
	"path"

	"github.com/fsnotify/fsnotify"
)

type LocalWatcher struct {
	dir          string
	filenames    []string
	handler      EventHandler
	errorHandler ErrorHandler
	watcher      *fsnotify.Watcher
	close        chan struct{}
}

func NewLocalWatcher(dir string, filenames []string) *LocalWatcher {
	return &LocalWatcher{
		dir:       dir,
		filenames: filenames,
		close:     make(chan struct{}),
	}
}

func (w *LocalWatcher) SetEventHandler(handler EventHandler) {
	w.handler = handler
}
func (w *LocalWatcher) SetErrorHandler(handler ErrorHandler) {
	w.errorHandler = handler
}

func (w *LocalWatcher) handle(e fsnotify.Event) error {
	if (!e.Has(fsnotify.Write)) && (!e.Has(fsnotify.Create)) {
		return nil
	}
	name := path.Base(e.Name)
	file := ""
	for _, filename := range w.filenames {
		ok, _ := path.Match(filename, name)
		if ok {
			file = filename
		}
	}
	if len(file) == 0 {
		return nil
	}
	err := w.handler(file)
	if err != nil {
		return err
	}
	return nil
}

func (w *LocalWatcher) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("new watcher error: %s", err)
	}
	go func() {
		defer watcher.Close()
		for {
			select {
			case e := <-watcher.Events:
				err := w.handle(e)
				if err != nil {
					w.errorHandler(err)
				}
			case err := <-watcher.Errors:
				if err != nil {
					w.errorHandler(err)
				}
			case <-w.close:
				return
			}
		}
	}()
	return watcher.Add(w.dir)
}

func (w *LocalWatcher) Close() error {
	close(w.close)
	return w.watcher.Close()
}
