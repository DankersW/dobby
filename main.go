package main

import (
	"fmt"
	"log"

	"github.com/DankersW/dobby/home-automation-ipc/generated/go/wsn"
	proto "github.com/golang/protobuf/proto"
)

func main() {
	fmt.Println("hi")
	data := []byte{10, 3, 83, 48, 49, 21, 0, 0, 234, 65}
	sensorData := &wsn.SensorData{}
	if err := proto.Unmarshal(data, sensorData); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	fmt.Printf("%v", sensorData)
}
