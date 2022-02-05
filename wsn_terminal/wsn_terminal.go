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
	wt.serial.write("thread monitor temp on")

	read := time.NewTicker(time.Duration(READ_INTERVAL) * time.Second)
	toggle := time.NewTicker(time.Duration(10) * time.Second)
	for {
		select {
		case <-read.C:
			wt.listen()
		case <-toggle.C:
			wt.serial.write("thread multi_light toggle")
		case <-wt.quit:
			read.Stop()
			return
			// Lwait for events that needs to be transmitted
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

func (wt *wsnTerminal) listen() {
	data, err := wt.serial.read()
	if err != nil {
		log.Errorf("Received an error while reading from UART, %s", err)
	}
	if data == "" || data == "uart:~$ " {
		return
	}
	if data == "uart:~$ " {
		log.Error(data)
	}
	sample := "[01:28:35.582,580] <inf> ot_coap: SensorData | 10 3 83 48 49 21 0 0 216 65 | 10 | fdde:ad00:beef:0:900a:1515:876a:92a4"
	log.Infof("Data: |%s| len: %d, sample: %d", data, len(data), len(sample))
}
