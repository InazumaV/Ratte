package watcher

type EventHandler func(filename string) error
type ErrorHandler func(err error)

type Watcher interface {
	SetEventHandler(handler EventHandler)
	SetErrorHandler(handler ErrorHandler)
	Watch() error
	Close() error
}
