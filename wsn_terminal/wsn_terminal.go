package wsn_terminal

import (
	"strings"
	"time"
	"unicode"

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
}

func (wt *wsnTerminal) Close() {
	wt.quit <- true
	if err := wt.serial.close(); err != nil {
		log.Errorf("failed to close serial terminal, %s", err)
	}
}

func (wt *wsnTerminal) listen() {
	rawData, err := wt.serial.read()
	if err != nil {
		log.Errorf("Received an error while reading from UART, %s", err)
	}
	if rawData == "" || strings.Contains(rawData, "uart:~$ ") {
		return
	}
	data := cleanup(rawData)
	log.Infof("Received: %q", data)
}

func cleanup(raw string) string {
	log.Debugf("raw |%q|", raw)
	clean := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, raw)
	if strings.HasPrefix(clean, "[8D[J") {
		clean = strings.Replace(clean, "[8D[J", "", 1)
	}
	return clean
}
