package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

func NewConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokers := []string{"localhost:29092"}
	conn, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Errorf("connection error, %s", err)
		return
	}

	// Calling ConsumePartition. It will open one connection per broker
	// and share it for all partitions that live on it.
	topic := "test"
	consumer, err := conn.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Consumer started ")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Count how many message processed
	msgCount := 0

	// Get signal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))
			case <-sigchan:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")

	if err := conn.Close(); err != nil {
		panic(err)
	}
}
