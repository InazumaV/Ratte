package acme

import (
	"fmt"
	"github.com/InazumaV/Ratte/conf"
	"path"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
)

type Acme struct {
	client *lego.Client
	c      *conf.ACME
}

func NewAcme(c *conf.ACME) (*Acme, error) {
	user, err := NewLegoUser(
		path.Join(c.Storage, "user", fmt.Sprintf("user-%s.json", c.Email)),
		c.Email)
	if err != nil {
		return nil, fmt.Errorf("create user error: %s", err)
	}
	lc := lego.NewConfig(user)
	//c.CADirURL = "http://192.168.99.100:4000/directory"
	lc.Certificate.KeyType = certcrypto.RSA2048
	client, err := lego.NewClient(lc)
	if err != nil {
		return nil, err
	}
	l := Acme{
		client: client,
		c:      c,
	}
	err = l.SetProvider()
	if err != nil {
		return nil, fmt.Errorf("set provider error: %s", err)
	}
	return &l, nil
}
