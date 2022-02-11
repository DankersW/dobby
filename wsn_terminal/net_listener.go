package wsn_terminal

import (
	"strings"

	"github.com/DankersW/dobby/home-automation-ipc/generated/go/wsn"
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
	rawData := wt.serial.read()
	if rawData == "" || strings.Contains(rawData, "uart:~$ ") {
		return
	}
	log.Infof("Received: %q", msgCleanup(rawData))
	msg := wt.parseMsg(rawData)
	if len(msg.data) == 0 {
		log.Debug("empty message")
		return
	}
	go wt.msgHandler(msg)
}

func (wt *wsnTerminal) msgHandler(msg wsnNodeMsg) {
	switch msg.breed {
	case int(wsn.MessageType_SENSOR_DATA):
		log.Info("handling Sensor")
	default:
		log.Info("not sup")
	}
}
