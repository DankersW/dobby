package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

func NewProducer() {
	config := producerConfig()

	brokers := []string{"localhost:29092"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Errorf("connection error, %s", err)
		return
	}

	defer producer.Close()

	topic := "test"
	message := fmt.Sprintf("hello- %s", time.Now())
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Error(err)
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
}
