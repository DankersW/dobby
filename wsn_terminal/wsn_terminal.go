package wsn_terminal

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	READ_INTERVAL = 1 // SECONDS
)

type wsnTerminal struct {
	serial *uart
	quit   chan bool
	// TODO: handlers to call on a certain type of message
}

type WsnTerminal interface {
	Start()
	Close()
}

func New(port string) (WsnTerminal, error) {
	serialTerm, err := newUartConnection(port)
	if err != nil {
		log.Errorf("Failed to open Serial connection to WSN gateway, %s", err.Error())
		return nil, err
	}

	wt := &wsnTerminal{
		serial: serialTerm,
		quit:   make(chan bool),
	}
	return wt, nil
}

func (wt *wsnTerminal) Start() {
	wt.serial.setup()

	read := time.NewTicker(time.Duration(READ_INTERVAL) * time.Second)
	toggle := time.NewTicker(time.Duration(10) * time.Second)
	for {
		select {
		case <-read.C:
			log.Info("Reading")
		case <-toggle.C:
			log.Info("toggle")
		case <-wt.quit:
			read.Stop()
			return
		}
	}

	count := 0
	for {
		data, err := wt.serial.read()
		if err != nil {
			log.Errorf("Failed to read from serial, %s", err)
			continue
		}
		log.Infof("Data: %s", data)
		time.Sleep(1 * time.Second)

		if count == 5 {
			wt.serial.write("thread multi_light toggle")
			count = 0
		}
		count++

		// Here we will wait for an event to come in and we handle it aka write it
	}
}

func (wt *wsnTerminal) Close() {
	wt.quit <- true
	if err := wt.serial.close(); err != nil {
		log.Errorf("failed to close serial terminal, %s", err)
	}
}
