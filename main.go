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
		exit := make(chan bool)
		brokers := []string{"localhost:29092"}
		topics := []string{"test"}
		groupId := "dobby"
		consumer, err := kafka.NewConsumer(brokers, groupId, topics, exit)
		if err != nil {
			log.Error(err)
		}
		go consumer.Serve()
		/*
			brokers := []string{"localhost:29092"}
			topics := []string{"test"}
			exit := make(chan bool)
			consumer, err := kafka.NewConsumer(brokers, topics, exit)
			if err != nil {
				log.Panic(err)
			}
			go consumer.Serve()
		*/
		publish := time.NewTicker(time.Duration(10) * time.Second)
		close := time.NewTicker(time.Duration(250) * time.Second)

		for {
			select {
			case <-publish.C:
				kafka.NewProducer()
			case <-close.C:
				log.Info("closing all consumers")
				exit <- true
				break
			}
		}
	} else {
		term, err := wsn_terminal.New(config.Wsn.Usb.Port)
		if err != nil {
			log.Errorf("Terminal failed to setup: %s", err.Error())
		}
		term.Start()
	}

}
