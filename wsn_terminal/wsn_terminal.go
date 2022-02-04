package wsn_terminal

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type wsnTerminal struct {
	usbPort    string
	serialTerm *uart
	// TODO: handlers to call on a certain type of message
}

type WsnTerminal interface {
	Start()
	Close()
}

func New(port string) WsnTerminal {
	wt := &wsnTerminal{
		usbPort: port,
	}
	return wt
}

func (wt *wsnTerminal) Start() {

	serialTerm, err := newUartConnection(wt.usbPort)
	if err != nil {
		log.Errorf("Failed to open Serial connection to WSN gateway, %s", err.Error())
		return
	}
	wt.serialTerm = serialTerm
	serialTerm.setup()

	count := 0
	for {
		data, err := serialTerm.read()
		if err != nil {
			log.Errorf("Failed to read from serial, %s", err)
			continue
		}
		log.Infof("Data: %s", data)
		time.Sleep(1 * time.Second)

		if count == 5 {
			serialTerm.write("thread multi_light toggle")
			count = 0
		}
		count++
	}
}

func (wt *wsnTerminal) Close() {
	if err := wt.serialTerm.close(); err != nil {
		log.Errorf("failed to close serial terminal, %s", err)
	}
}
