package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/DankersW/dobby/config"
	"github.com/DankersW/dobby/kafka"
	"github.com/DankersW/dobby/wsn_terminal"
	log "github.com/sirupsen/logrus"
)

const (
	queue_size = 20
)

// TODO: get the log level from config file so that prod docker only prints Warning msgs

func main() {
	config := config.Get()

	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//kafka.Example()

	// TODO: get kafka from config
	brokers := []string{"localhost:29092"}

	log.Info("Starting IPC handlder")
	txQueue := make(chan kafka.KafkaTxQueue, queue_size)
	producer, err := kafka.NewProducer(brokers, txQueue)
	if err != nil {
		log.Fatalf("Failed to setup kafka producer, %s", err.Error())
		return
	}
	go producer.Serve()
	log.Info("Started IPC handler")

	log.Info("Starting WSN terminal CLI")
	term, err := wsn_terminal.New(config.Wsn.Usb.Port)
	if err != nil {
		log.Fatalf("WSN terminal CLI failed to setup: %s", err.Error())
	}
	term.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	producer.Shutdown()
	mainCtx.Done()

}
