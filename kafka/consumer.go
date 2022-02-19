package kafka

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type consumer struct {
	exit           chan bool
	conn           sarama.Consumer
	topicConsumers map[string]sarama.PartitionConsumer
}

type Consumer_ interface {
	Serve()
	Close()
}

func NewConsumer(brokers []string, topics []string, exit chan bool) (Consumer_, error) {
	config := consumerConfig()

	conn, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	tcs := make(map[string]sarama.PartitionConsumer)
	/*
		for _, topic := range topics {
			tc, err := conn.ConsumePartition(topic, 0, sarama.OffsetOldest)
			if err != nil {
				return nil, err
			}
			tcs[topic] = tc
		}
	*/

	c := &consumer{
		conn:           conn,
		exit:           exit,
		topicConsumers: tcs,
	}
	return c, nil
}

func (c *consumer) Serve() {
	log.Info("Serve")

	for topic, consumer := range c.topicConsumers {
		log.Infof("started listening on topic %q", topic)
		go c.listen(consumer)
	}
}

func (c *consumer) listen(consumer sarama.PartitionConsumer) {
	for {
		select {
		case err := <-consumer.Errors():
			log.Errorf("received error on topic %q, err: %s", err.Topic, err.Err)
		case msg := <-consumer.Messages():
			log.Infof("Received message on topic %q: %s", string(msg.Topic), string(msg.Value))
		case <-c.exit:
			log.Info("Closing the consumer")
			return
		}
	}
}

func (c *consumer) Close() {
	log.Info("Close")
	err := c.conn.Close()
	log.Infof("Close error: %s", err.Error())
}

func consumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	//config.ClientID = "Dobby-consumers"
	//config.Consumer.Offsets.Initial = sarama.OffsetNewest
	return config
}
