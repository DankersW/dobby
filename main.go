package main

import (
	"time"

	"github.com/DankersW/dobby/config"
	"github.com/DankersW/dobby/kafka"
	"github.com/DankersW/dobby/wsn_terminal"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := config.Get()

	stage := "test"

	if stage == "test" {
		go kafka.NewConsumer()
		for {
			kafka.NewProducer()
			time.Sleep(2 * time.Second)
		}
	} else {
		term, err := wsn_terminal.New(config.Wsn.Usb.Port)
		if err != nil {
			log.Errorf("Terminal failed to setup: %s", err.Error())
		}
		term.Start()
	}

}
