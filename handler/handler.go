package handler

import (
	"Ratte/acme"
	"Ratte/conf"
	"github.com/sirupsen/logrus"
	"sync/atomic"
)
import "github.com/Yuzuki616/Ratte-Interface/core"
import "github.com/Yuzuki616/Ratte-Interface/panel"

type Handler struct {
	c         core.Core
	p         panel.Panel
	nodeName  string
	acme      *acme.Acme
	l         *logrus.Entry
	userList  []panel.UserInfo
	userHash  map[string]struct{}
	nodeAdded atomic.Bool
	*conf.Options
}

func New(
	c core.Core,
	p panel.Panel,
	nodeName string,
	ac *acme.Acme,
	l *logrus.Entry,
	opts *conf.Options) *Handler {
	return &Handler{
		c:        c,
		p:        p,
		nodeName: nodeName,
		userList: make([]panel.UserInfo, 0),
		userHash: make(map[string]struct{}),
		acme:     ac,
		l:        l,
		Options:  opts,
	}
}

func (h *Handler) Close() error {
	if h.nodeAdded.Load() {
		err := h.execHookCmd(h.Hook.BeforeDelNode, h.nodeName)
		if err != nil {
			h.l.WithError(err).Warn("Exec before del node hook failed")
		}
		defer func() {
			err = h.execHookCmd(h.Hook.AfterDelNode, h.nodeName)
			if err != nil {
				h.l.WithError(err).Warn("Exec after del node hook failed")
			}
		}()
		return h.c.DelNode(h.nodeName)
	}
	return nil
}
