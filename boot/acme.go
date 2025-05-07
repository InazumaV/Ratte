package boot

import (
	"github.com/InazumaV/Ratte/acme"
	"github.com/InazumaV/Ratte/conf"
	log "github.com/sirupsen/logrus"
)

type AcmeGroup map[string]*acme.Acme

func initAcme(a []conf.ACME) AcmeGroup {
	acmes := make(AcmeGroup, len(a))
	for _, a := range a {
		ac, err := acme.NewAcme(&a)
		if err != nil {
			log.WithError(err).Fatal("New acme failed")
		}
		acmes[a.Name] = ac
	}
	return acmes
}
