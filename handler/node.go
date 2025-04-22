package handler

import (
	"fmt"
	"github.com/InazumaV/Ratte-Interface/core"
	"github.com/InazumaV/Ratte-Interface/panel"
	"github.com/InazumaV/Ratte/common/maps"
	"github.com/InazumaV/Ratte/common/number"
)

func (h *Handler) PullNodeHandle(n *panel.NodeInfo) error {
	if h.nodeAdded.Load() {
		err := h.c.DelNode(h.nodeName)
		if err != nil {
			return fmt.Errorf("del node error: %w", err)
		}
	} else {
		if n.TlsType() != panel.NoTls {
			h.needTls = true
			err := h.acme.CreateCert(h.Cert.CertPath, h.Cert.KeyPath, h.Cert.Domain)
			if err != nil {
				return fmt.Errorf("create cert error: %w", err)
			}
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
	if h.Hook.BeforeAddNode != "" {
		err := h.execHookCmd(h.Hook.BeforeAddNode, h.nodeName, protocol, port)
		if err != nil {
			h.l.WithError(err).Error("Exec before add node hook failed")
		}
	}

	ni := (*core.NodeInfo)(n)
	ni.OtherOptions = maps.Merge[string, any](ni.OtherOptions, h.Options.Expand)
	ni.Limit.IPLimit = number.SelectBigger(ni.Limit.IPLimit, h.Limit.IPLimit)
	ni.Limit.SpeedLimit = number.SelectBigger(ni.Limit.SpeedLimit, uint64(h.Limit.SpeedLimit))
	err := h.c.AddNode(&core.AddNodeParams{
		Name:     h.nodeName,
		NodeInfo: ni,
		TlsOptions: core.TlsOptions{
			CertPath: h.Cert.CertPath,
			KeyPath:  h.Cert.KeyPath,
		},
	})
	if err != nil {
		return fmt.Errorf("add node error: %w", err)
	}
	if h.Hook.AfterAddNode != "" {
		err = h.execHookCmd(h.Hook.AfterAddNode, h.nodeName, protocol, port)
		if err != nil {
			h.l.WithError(err).Warn("Exec after add node hook failed")
		}
	}
	if h.nodeAdded.Load() {
		h.nodeAdded.Store(true)
	}
	return nil
}
