package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/DankersW/dobby/config"
	"github.com/DankersW/dobby/kafka"
	"github.com/DankersW/dobby/models"
	"github.com/DankersW/dobby/wsn_terminal"
	log "github.com/sirupsen/logrus"
)

const (
	queue_size = 20
)

var conf models.Config

func init() {
	conf = config.Get()
	logLevel, err := log.ParseLevel(conf.Log.Level)
	if err != nil {
		log.Errorf("Failed to set minumimum log level, %s", err)
	}
	log.SetLevel(logLevel)
}

func main() {
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	txQueue := make(chan kafka.KafkaTxQueue, queue_size)

	log.Info("Starting IPC handlder")
	producer, err := kafka.NewProducer(conf.Kafka.Brokers, txQueue)
	if err != nil {
		log.Fatalf("Failed to setup kafka producer, %s", err.Error())
		return
	}
	go producer.Serve()
	log.Info("Started IPC handler")

	log.Info("Starting WSN terminal CLI")
	term, err := wsn_terminal.New(conf.Wsn.Usb.Port, txQueue)
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
