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

// TODO: get the log level from config file so that prod docker only prints Warning msgs

func main() {
	config := config.Get()

	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stage := "test"
	if stage == "test" {

		testTopic := make(kafka.MsgCallback)

		brokers := []string{"localhost:29092"}
		topics := map[string]kafka.MsgCallback{
			"test": testTopic,
		}
		consumer, err := kafka.NewConsumer(brokers, topics)
		if err != nil {
			log.Errorf("Failed to setup kafka consumer, %s", err.Error())
			return
		}
		producer, err := kafka.NewProducer(brokers)
		if err != nil {
			log.Errorf("Failed to setup kafka producer, %s", err.Error())
			return
		}
		log.Info("Kafka good")
		go consumer.Serve(mainCtx)

		publish := time.NewTicker(time.Duration(10) * time.Second)
		close := time.NewTicker(time.Duration(25) * time.Second)

		for {
			select {
			case <-publish.C:
				producer.Send("test", []byte("hi a msg"))
			case <-close.C:
				log.Info("closing consumer")
				cancel()
				consumer.Close()
				publish.Stop()
				break
			case msg := <-testTopic:
				log.Infof("receveid msg on test topic: %s", string(msg))
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

	mainCtx.Done()

}
