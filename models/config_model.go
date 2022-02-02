package models

type Config struct {
	Wsn struct {
		Usb struct {
			Port string `yaml:"port"`
		}
	}
}
