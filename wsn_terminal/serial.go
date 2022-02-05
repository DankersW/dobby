package wsn_terminal

import (
	"bufio"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type uart struct {
	stream *serial.Port
}

func newUartConnection(device string) (*uart, error) {
	config := &serial.Config{
		Name:        device,
		Baud:        9600,
		ReadTimeout: 1,
		Size:        8,
	}
	stream, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}
	return &uart{stream: stream}, nil
}

func (u *uart) setup() {
	u.write("")
	u.write("shell echo off")
	u.write("shell colors off")
	u.write("clear")
}

func (u *uart) close() error {
	return u.stream.Close()
}

func (u *uart) read() (string, error) {
	scanner := bufio.NewScanner(u.stream)
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	log.Infof("err: %s", scanner.Err())
	return "", scanner.Err()
}

func (u *uart) write(cmd string) {
	writer := bufio.NewWriter(u.stream)
	writer.Reset(u.stream)
	cmd += "\n"
	if _, err := writer.WriteString(cmd); err != nil {
		log.Error("failed to write to serial port, %s", err)
	}
	if err := writer.Flush(); err != nil {
		log.Error("failed to flush, %s", err)
	}
}