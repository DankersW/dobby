package main

import (
	"github.com/DankersW/dobby/config"
	"github.com/DankersW/dobby/wsn_terminal"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := config.Get()

	term, err := wsn_terminal.New(config.Wsn.Usb.Port)
	if err != nil {
		log.Errorf("Terminal failed to setup: %s", err.Error())
	}
	term.Start()
}
