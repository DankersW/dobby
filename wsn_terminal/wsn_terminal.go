package wsn_terminal

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	READ_INTERVAL = 100 // MILI_SECONDS
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
