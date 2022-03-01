package wsn_terminal

import (
	"strings"

	"github.com/DankersW/dobby/home-automation-ipc/generated/go/wsn"
	"github.com/DankersW/dobby/kafka"
	proto "github.com/golang/protobuf/proto"
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

func (wt *wsnTerminal) sensorDataHandler(data []byte) {
	sensorData := &wsn.SensorData{}
	if err := proto.Unmarshal(data, sensorData); err != nil {
		log.Warnf("Failed to parse SensorData msg, %s", err.Error())
	}
	log.Debugf("%v\n", sensorData)
	// TODO: parse recieved data to a differnt data format
	txItem := kafka.KafkaTxQueue{
		Topic: "wsn.sensor-data.telemetry",
		Data:  data,
	}
	wt.ipcTxQueue <- txItem
}