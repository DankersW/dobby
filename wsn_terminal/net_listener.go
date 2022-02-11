package wsn_terminal

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

type wsnNodeMsg struct {
	breed int
	data  []byte
	ip    string
}

// TODO: proto decoding
// TODO: message decoding
// TODO: Use a channel to pass data inbetween someone that needs it and the listner

func (wt *wsnTerminal) listen() {
	rawData, err := wt.serial.read()
	if err != nil {
		log.Errorf("Received an error while reading from UART, %s", err)
	}
	if rawData == "" || strings.Contains(rawData, "uart:~$ ") {
		return
	}
	log.Infof("Received: %q", msgCleanup(rawData))
	msg := wt.parseMsg(rawData)
	if len(msg.data) == 0 {
		log.Debug("empty message")
	}
	log.Infof("Parsed: %v", msg)
}
