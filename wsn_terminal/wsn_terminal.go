package wsn_terminal

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type wsnTerminal struct {
	serial *uart
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
	}
	return wt, nil
}

func (wt *wsnTerminal) Start() {

	wt.serial.setup()

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
	if err := wt.serial.close(); err != nil {
		log.Errorf("failed to close serial terminal, %s", err)
	}
}
