package main

import (
	"github.com/InazumaV/Ratte-Interface/panel"
	v2b "github.com/InazumaV/Ratte-Panel-V2b"
	log "github.com/sirupsen/logrus"
)

func main() {
	c, err := panel.NewServer(nil, v2b.NewPanel())
	if err != nil {
		log.Fatalln(err)
	}
	if err = c.Run(); err != nil {
		log.Fatalln(err)
	}
}
