package handler

import (
	"Ratte/common/maps"
	"fmt"
	"github.com/Yuzuki616/Ratte-Interface/core"
	"github.com/Yuzuki616/Ratte-Interface/panel"
	"github.com/Yuzuki616/Ratte-Interface/params"
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
	err := h.c.AddNode(&core.AddNodeParams{
		NodeInfo: core.NodeInfo{
			CommonNodeInfo: params.CommonNodeInfo{
				Type:        n.Type,
				VMess:       n.VMess,
				Shadowsocks: n.Shadowsocks,
				Trojan:      n.Trojan,
				Hysteria:    n.Hysteria,
				Other:       n.Other,
				ExpandParams: params.ExpandParams{
					OtherOptions: maps.Merge(n.OtherOptions, h.Expand),
					CustomData:   n.CustomData,
				},
			},
			TlsOptions: core.TlsOptions{
				CertPath: h.Cert.CertPath,
				KeyPath:  h.Cert.KeyPath,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("add node error: %w", err)
	}
	if h.nodeAdded.Load() {
		h.nodeAdded.Store(true)
	}
	return nil
}
