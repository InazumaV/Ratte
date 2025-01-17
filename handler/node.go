package handler

import (
	"fmt"
	"github.com/Yuzuki616/Ratte-Interface/core"
	"github.com/Yuzuki616/Ratte-Interface/panel"
)

func (h *Handler) PullNodeHandle(n *panel.NodeInfo) error {
	if h.nodeAdded.Load() {
		err := h.c.DelNode(h.nodeName)
		if err != nil {
			return fmt.Errorf("del node error: %w", err)
		}
	} else {
		err := h.acme.CreateCert(h.Cert.CertPath, h.Cert.KeyPath, h.Cert.Domain)
		if err != nil {
			return fmt.Errorf("create cert error: %w", err)
		}
	}
	var protocol, port string
	switch n.Type {
	case "vmess":
		protocol = "vmess"
		port = n.VMess.Port
	case "vless":
		protocol = "vless"
		port = n.VLess.Port
	case "shadowsocks":
		protocol = "shadowsocks"
		port = n.Shadowsocks.Port
	case "trojan":
		protocol = "trojan"
		port = n.Trojan.Port
	case "other":
		protocol = "other"
		port = n.Other.Port
	}
	err := h.execHookCmd(h.Hook.BeforeAddNode, h.nodeName, protocol, port)
	if err != nil {
		h.l.WithError(err).Error("Exec before add node hook failed")
	}
	err = h.c.AddNode(&core.AddNodeParams{
		NodeInfo: (*core.NodeInfo)(n),
		TlsOptions: core.TlsOptions{
			CertPath: h.Cert.CertPath,
			KeyPath:  h.Cert.KeyPath,
		},
	})
	if err != nil {
		return fmt.Errorf("add node error: %w", err)
	}
	err = h.execHookCmd(h.Hook.AfterAddNode, h.nodeName, protocol, port)
	if err != nil {
		h.l.WithError(err).Warn("Exec after add node hook failed")
	}
	if h.nodeAdded.Load() {
		h.nodeAdded.Store(true)
	}
	return nil
}
