package models

import (
	"fmt"
	"time"

	"github.com/DankersW/dobby/home-automation-ipc/generated/go/ipc"
	"github.com/DankersW/dobby/home-automation-ipc/generated/go/wsn"
	"github.com/golang/protobuf/proto"
)

func TransformWsnSensorDataToIpcSensorDataTelemetry(data []byte) ([]byte, error) {
	sensorData := &wsn.SensorData{}
	if err := proto.Unmarshal(data, sensorData); err != nil {
		return nil, fmt.Errorf("failed to parse SensorData msg, %s", err.Error())
	}

	var temp float32 = 0.0
	if sensorData.Temperature != 0 {
		temp = float32(sensorData.Temperature / 10)
	}

	var humi float32 = 0.0
	if sensorData.Humidity != 0 {
		humi = float32(sensorData.Humidity / 10)
	}

	telemetryData := ipc.WsnSensorDataTelemetry{
		Timestamp:   uint32(time.Now().Unix()),
		SensorId:    sensorData.SensorId,
		Temperature: temp,
		Humidity:    humi,
	}
	serialData, err := proto.Marshal(&telemetryData)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize IpcSensorDataTelemetry msg, %s", err.Error())
	}
	return serialData, nil
}
