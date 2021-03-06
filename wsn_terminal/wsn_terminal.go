package wsn_terminal

import (
	"fmt"
	"time"

	"github.com/DankersW/dobby/kafka"
	log "github.com/sirupsen/logrus"
)

const (
	READ_INTERVAL = 100 // MILI_SECONDS
)

type wsnTerminal struct {
	serial     *uart
	quit       chan bool
	ipcTxQueue chan kafka.KafkaTxQueue
}

type WsnTerminal interface {
	Start()
	Close()
}

// TODO: from the handlers push data (topic, msg) onto a channel, then from kafka producer activly wait to transmit data

func New(port string, ipcTxQueue chan kafka.KafkaTxQueue) (WsnTerminal, error) {
	serialTerm, err := newUartConnection(port)
	if err != nil {
		return nil, fmt.Errorf("failed to open Serial connection to WSN gateway, %s", err.Error())
	}

	wt := &wsnTerminal{
		serial:     serialTerm,
		quit:       make(chan bool),
		ipcTxQueue: ipcTxQueue,
	}
	return wt, nil
}

func (wt *wsnTerminal) Start() {
	wt.serial.setup()
	wt.serial.write("thread monitor temp on")

	// TODO: add all times into a array or something to handle all timer related things
	read := time.NewTicker(time.Duration(READ_INTERVAL) * time.Millisecond)
	toggle := time.NewTicker(time.Duration(3) * time.Second)
	for {
		select {
		case <-read.C:
			wt.listen()
		case <-toggle.C:
			wt.serial.write("thread multi_light toggle")
		case <-wt.quit:
			read.Stop()
			toggle.Stop()
			return
			// Lwait for events that needs to be transmitted
		}
	}
}

func (wt *wsnTerminal) Close() {
	wt.quit <- true
	if err := wt.serial.close(); err != nil {
		log.Errorf("failed to close serial terminal, %s", err)
	}
}
