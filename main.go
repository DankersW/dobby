package main

import (
	"fmt"

	"github.com/DankersW/dobby/config"
	"github.com/DankersW/dobby/home-automation-ipc/generated/go/wsn"
	"github.com/DankersW/dobby/wsn_terminal"
	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("hi")
	data := []byte{10, 3, 83, 48, 49, 21, 0, 0, 234, 65}
	log.Info(data)
	sensorData := &wsn.SensorData{}
	if err := proto.Unmarshal(data, sensorData); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	fmt.Printf("%v\n", sensorData)

	config := config.Get()

	term, err := wsn_terminal.New(config.Wsn.Usb.Port)
	if err != nil {
		log.Errorf("terminal failed to setup: %s", err.Error())
	}
	term.Start()
}
