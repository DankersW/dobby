package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/DankersW/dobby/config"
	"github.com/DankersW/dobby/kafka"
	"github.com/DankersW/dobby/wsn_terminal"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := config.Get()

	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stage := "test"
	if stage == "test" {
		brokers := []string{"localhost:29092"}
		topics := []string{"test"}
		consumer, err := kafka.NewConsumer(brokers, topics)
		if err != nil {
			log.Panicf("failed to setup kafka consumer, %s", err.Error())
		}
		go consumer.Serve(mainCtx)

		publish := time.NewTicker(time.Duration(10) * time.Second)
		close := time.NewTicker(time.Duration(25) * time.Second)

		for {
			select {
			case <-publish.C:
				kafka.NewProducer()
			case <-close.C:
				log.Info("closing consumer")
				cancel()
				consumer.Close()
				publish.Stop()
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	//consumer.Close()
	mainCtx.Done()

}
