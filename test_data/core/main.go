package main

import (
	"github.com/InazumaV/Ratte-Core-Xray"
	"github.com/InazumaV/Ratte-Interface/core"
	log "github.com/sirupsen/logrus"
)

func main() {
	c, err := core.NewServer(nil, xray.NewXray())
	if err != nil {
		log.Fatalln(err)
	}
	if err = c.Run(); err != nil {
		log.Fatalln(err)
	}
}
