package handler

import "fmt"

func (h *Handler) RenewCertHandle() error {
	if !h.needTls {
		return nil
	}
	err := h.acme.RenewCert(h.Options.Cert.CertPath, h.Options.Cert.KeyPath, h.Options.Cert.Domain)
	if err != nil {
		return fmt.Errorf("renew cert error: %w", err)
	}
	return nil
}
