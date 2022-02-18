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
		brokers := []string{"localhost:29092"}
		topics := []string{"test"}
		exit := make(chan bool)
		consumer, err := kafka.NewConsumer(brokers, topics, exit)
		if err != nil {
			log.Panic(err)
		}
		go consumer.Serve()
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
