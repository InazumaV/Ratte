package handler

import (
	"github.com/InazumaV/Ratte/acme"
	"github.com/InazumaV/Ratte/conf"
	"github.com/sirupsen/logrus"
	"sync/atomic"
)
import "github.com/InazumaV/Ratte-Interface/core"
import "github.com/InazumaV/Ratte-Interface/panel"

type Handler struct {
	c         core.Core
	p         panel.Panel
	nodeName  string
	acme      *acme.Acme
	l         *logrus.Entry
	userList  []panel.UserInfo
	userHash  map[string]struct{}
	nodeAdded atomic.Bool
	needTls   bool
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
