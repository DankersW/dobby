package kafka

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func Example() {
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Topic setup for consumer callback
	testTopicClb := make(MsgCallback)
	topics := map[string]MsgCallback{
		"test": testTopicClb,
	}
	brokers := []string{"localhost:29092"}

	// Producer
	producer, err := NewProducer(brokers)
	if err != nil {
		log.Fatalf("Failed to setup kafka producer, %s", err.Error())
		return
	}

	// Consumer
	consumer, err := NewConsumer(brokers, topics)
	if err != nil {
		log.Fatalf("Failed to setup kafka consumer, %s", err.Error())
	}
	go consumer.Serve(mainCtx)

	// Main loop
	running := true
	publish := time.NewTicker(time.Duration(3) * time.Second)
	close := time.NewTicker(time.Duration(10) * time.Second)
	for running {
		select {
		case <-publish.C:
			producer.Send("test", []byte("hi a msg"))
		case <-close.C:
			log.Info("stopping")
			running = false
			publish.Stop()
			close.Stop()
		case msg := <-testTopicClb:
			log.Infof("receveid msg | test: %s", string(msg))
		}
	}
	cancel()
	producer.Shutdown()
	consumer.Close()
	log.Info("apps closed")

	mainCtx.Done()
	log.Info("all done!")
}
