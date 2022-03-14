package wsn_terminal

import (
	"strings"

	"github.com/DankersW/dobby/home-automation-ipc/generated/go/wsn"
	"github.com/DankersW/dobby/kafka"
	"github.com/DankersW/dobby/models"
	log "github.com/sirupsen/logrus"
)

type wsnNodeMsg struct {
	breed int
	data  []byte
	ip    string
}

func (wt *wsnTerminal) listen() {
	rawData := wt.serial.read()
	if rawData == "" || strings.Contains(rawData, "uart:~$ ") {
		return
	}
	log.Debugf("Received: %q", msgCleanup(rawData))
	msg := wt.parseMsg(rawData)
	if len(msg.data) == 0 {
		return
	}
	go wt.msgHandler(msg)
}

func (wt *wsnTerminal) msgHandler(msg wsnNodeMsg) {
	switch msg.breed {
	case int(wsn.MessageType_SENSOR_DATA):
		wt.sensorDataHandler(msg.data)
	default:
		log.Warnf("Received an uknown type %q", msg.breed)
	}
}

func (wt *wsnTerminal) sensorDataHandler(rawData []byte) {
	telemetryData, err := models.TransformWsnSensorDataToIpcSensorDataTelemetry(rawData)
	if err != nil {
		log.Warn(err)
	}
	txItem := kafka.KafkaTxQueue{
		Topic: "wsn.sensor-data.telemetry", // TODO: Add this to a config file
		Data:  telemetryData,
	}
	wt.ipcTxQueue <- txItem
}
