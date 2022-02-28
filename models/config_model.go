package models

type Config struct {
	Wsn struct {
		Usb struct {
			Port string `yaml:"port"`
		}
	}
	Kafka KafkaConfig
	Log   struct {
		Level string `yaml:"level"`
	}
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}
