package models

type Config struct {
	Wsn struct {
		Usb struct {
			Port string `yaml:"port"`
		}
	}
	Kafka KafkaConfig
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}
