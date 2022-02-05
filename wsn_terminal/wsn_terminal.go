package wsn_terminal

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	READ_INTERVAL = 1 // SECONDS
)

type wsnTerminal struct {
	serial           *uart
	quitPeriodicWork chan bool
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
		serial:           serialTerm,
		quitPeriodicWork: make(chan bool),
	}
	return wt, nil
}

func (wt *wsnTerminal) Start() {
	wt.serial.setup()

	go wt.periodicRead(READ_INTERVAL)

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
	}
}

func (wt *wsnTerminal) Close() {
	wt.quitPeriodicWork <- true
	if err := wt.serial.close(); err != nil {
		log.Errorf("failed to close serial terminal, %s", err)
	}
}

func (wt *wsnTerminal) periodicRead(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Info("hi")
		case <-wt.quitPeriodicWork:
			ticker.Stop()
			return
		}
	}
}
