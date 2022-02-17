package kafka

import "github.com/Shopify/sarama"

func consumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return config
}
