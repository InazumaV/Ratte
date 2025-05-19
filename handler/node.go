package handler

import (
	"fmt"
	"github.com/InazumaV/Ratte-Interface/core"
	"github.com/InazumaV/Ratte-Interface/panel"
	"github.com/InazumaV/Ratte/common/maps"
	"github.com/InazumaV/Ratte/common/number"
	"strconv"
)

func (h *Handler) PullNodeHandle(n *panel.NodeInfo) error {
	if h.nodeAdded.Load() {
		err := h.c.DelNode(h.nodeName)
		if err != nil {
			return fmt.Errorf("del node error: %w", err)
		}
	} else {
		if n.Security == "" {
			h.needTls = true
			err := h.acme.CreateCert(h.Cert.CertPath, h.Cert.KeyPath, h.Cert.Domain)
			if err != nil {
				return fmt.Errorf("create cert error: %w", err)
			}
		}
	}

	if h.Hook.BeforeAddNode != "" {
		err := h.execHookCmd(h.Hook.BeforeAddNode, h.nodeName, n.Type, strconv.Itoa(n.Port))
		if err != nil {
			h.l.WithError(err).Error("Exec before add node hook failed")
		}
	}

	ni := (*core.NodeInfo)(n)
	ni.Options = maps.Merge[string, any](ni.Options, h.Options.Expand)
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
		err = h.execHookCmd(h.Hook.AfterAddNode, h.nodeName, ni.Type, strconv.Itoa(ni.Port))
		if err != nil {
			h.l.WithError(err).Warn("Exec after add node hook failed")
		}
	}
	if h.nodeAdded.Load() {
		h.nodeAdded.Store(true)
	}
	return nil
}
