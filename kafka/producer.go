package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type KafkaTxQueue struct {
	Topic string
	Data  []byte
}

type producer struct {
	syncProducer sarama.SyncProducer
	txQueue      chan KafkaTxQueue
}

type Producer interface {
	Send(topic string, data []byte) error
	Serve()
	Shutdown() error
}

func NewProducer(brokers []string, queue chan KafkaTxQueue) (Producer, error) {
	config := producerConfig()
	syncPoducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka broker, %s", err)
	}

	producer := &producer{
		syncProducer: syncPoducer,
		txQueue:      queue,
	}
	return producer, nil
}

func producerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	return config
}

func (p *producer) Shutdown() error {
	return p.syncProducer.Close()
}

func (p *producer) Send(topic string, data []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}
	_, _, err := p.syncProducer.SendMessage(msg)
	return err
}

func (p *producer) Serve() {
	for msg := range p.txQueue {
		if err := p.Send(msg.Topic, msg.Data); err != nil {
			log.Errorf("Failed to send msg on topic %q, %s", msg.Topic, err.Error())
		}
	}
}
