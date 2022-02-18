package kafka

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type consumer struct {
	exit chan bool
	conn sarama.Consumer
}

type Consumer interface {
	Serve()
	Close()
}

func NewConsumer(brokers []string, topics string, exit chan bool) (Consumer, error) {
	config := consumerConfig()

	conn, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	broker, err := conn.ConsumePartition(topics, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	c := &consumer{
		conn: conn,
		exit: exit,
	}
	return c, nil
}

func (c *consumer) Serve() {
	log.Info("Serve")
}

func (c *consumer) Close() {
	log.Info("Close")
	err := c.conn.Close()
	log.Infof("Close error: %s", err.Error())
}

func consumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return config
}
