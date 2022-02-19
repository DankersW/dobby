package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type consumer struct {
	ctx    context.Context
	client sarama.ConsumerGroup
	topics []string
}

type Consumer interface {
	Serve()
}

func NewConsumer(brokers []string, groupId string, topics []string, exit chan bool) (Consumer, error) {
	config := consumerConfig()

	ctx, _ := context.WithCancel(context.Background())

	client, err := sarama.NewConsumerGroup(brokers, groupId, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	consumer := &consumer{
		ctx:    ctx,
		client: client,
		topics: topics,
	}

	log.Println("Sarama consumer up and running!...")

	return consumer, nil
}

func (c *consumer) Serve() {
	for {
		if err := c.client.Consume(c.ctx, c.topics, c); err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if c.ctx.Err() != nil {
			return
		}
	}
}

func (c *consumer) Close() {
	log.Info("Close")
	//TODO: cancel context
	//err := c.conn.Close()
	//log.Infof("Close error: %s", err.Error())
}

/*
func (c *Consumer) Serve() {
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
*/
func consumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	version, err := sarama.ParseKafkaVersion("2.1.1") // FIXME: get real value here
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	config.Version = version
	return config
}

////

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Info("Kafka consumer is setup")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}
