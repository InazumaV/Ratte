package acme

import (
	"fmt"
	"github.com/InazumaV/Ratte/common/file"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/providers/dns"
	"os"
	"path"
	"time"
)

func checkPath(p string) error {
	if !file.IsExist(path.Dir(p)) {
		err := os.MkdirAll(path.Dir(p), 0755)
		if err != nil {
			return fmt.Errorf("create dir error: %s", err)
		}
	}
	return nil
}

func (l *Acme) SetProvider() error {
	switch l.c.Provider {
	case "http":
		err := l.client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80"))
		if err != nil {
			return err
		}
	case "dns":
		for k, v := range l.c.DNSEnv {
			os.Setenv(k, v)
		}
		p, err := dns.NewDNSChallengeProviderByName(l.c.Provider)
		if err != nil {
			return fmt.Errorf("create dns challenge provider error: %s", err)
		}
		err = l.client.Challenge.SetDNS01Provider(p)
		if err != nil {
			return fmt.Errorf("set dns provider error: %s", err)
		}
	default:
		return fmt.Errorf("unsupported provider %s", l.c.Provider)
	}
	return nil
}

func (l *Acme) CreateCert(certPath, keyPath, domain string) (err error) {
	if certPath == "" || keyPath == "" {
		return fmt.Errorf("cert file path or key file path not exist")
	}
	if file.IsExist(certPath) && file.IsExist(keyPath) {
		return l.RenewCert(certPath, keyPath, domain)
	}
	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}
	certificates, err := l.client.Certificate.Obtain(request)
	if err != nil {
		return fmt.Errorf("obtain certificate error: %s", err)
	}
	err = l.writeCert(certPath, keyPath, certificates)
	return nil
}

func (l *Acme) RenewCert(certPath, keyPath, domain string) error {
	file, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("read cert file error: %s", err)
	}
	if e, err := l.CheckCert(file); !e {
		return nil
	} else if err != nil {
		return fmt.Errorf("check cert error: %s", err)
	}
	res, err := l.client.Certificate.Renew(certificate.Resource{
		Domain:      domain,
		Certificate: file,
	}, true, false, "")
	if err != nil {
		return err
	}
	err = l.writeCert(certPath, keyPath, res)
	return nil
}

func (l *Acme) CheckCert(file []byte) (bool, error) {
	cert, err := certcrypto.ParsePEMCertificate(file)
	if err != nil {
		return false, err
	}
	notAfter := int(time.Until(cert.NotAfter).Hours() / 24.0)
	if notAfter > 30 {
		return false, nil
	}
	return true, nil
}
func (l *Acme) writeCert(cert, key string, certificates *certificate.Resource) error {
	err := checkPath(cert)
	if err != nil {
		return fmt.Errorf("check path error: %s", err)
	}
	err = os.WriteFile(cert, certificates.Certificate, 0644)
	if err != nil {
		return err
	}
	err = checkPath(key)
	if err != nil {
		return fmt.Errorf("check path error: %s", err)
	}
	err = os.WriteFile(key, certificates.PrivateKey, 0644)
	if err != nil {
		return err
	}
	return nil
}
