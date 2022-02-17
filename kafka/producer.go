package kafka

import "github.com/Shopify/sarama"

func producerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return config
}
