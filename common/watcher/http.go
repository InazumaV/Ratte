package watcher

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPWatcher struct {
	hash         [32]byte
	url          string
	interval     uint
	handler      EventHandler
	errorHandler ErrorHandler
	close        chan struct{}
}

func NewHTTPWatcher(url string, interval uint) *HTTPWatcher {
	return &HTTPWatcher{
		url:      url,
		interval: interval,
	}
}

func (w *HTTPWatcher) handle() error {
	rsp, err := http.Get(w.url)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer rsp.Body.Close()
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return fmt.Errorf("read body error: %w", err)
	}
	h := sha256.Sum256(b)
	if bytes.Equal(w.hash[:], h[:]) {
		return nil
	}
	w.hash = h
	err = w.handler(w.url)
	if err != nil {
		return fmt.Errorf("handle error: %w", err)
	}
	return nil
}

func (w *HTTPWatcher) SetEventHandler(handler EventHandler) {
	w.handler = handler
}

func (w *HTTPWatcher) SetErrorHandler(handler ErrorHandler) {
	w.errorHandler = handler
}

func (w *HTTPWatcher) Watch() error {
	go func() {
		for range time.Tick(time.Duration(w.interval) * time.Second) {
			select {
			case <-w.close:
				return
			default:
			}
			w.errorHandler(w.handle())
		}
	}()
	return nil
}

func (w *HTTPWatcher) Close() error {
	close(w.close)
	return nil
}
