package wsn_terminal

import (
	"bufio"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type uart struct {
	stream *serial.Port
}

// TODO: READ and write lock

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

// FIXME: sometimes this function really messes up
func (u *uart) read() (string, error) {
	scanner := bufio.NewScanner(u.stream)
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", scanner.Err()
}

// TODO: READ and write lock

func (u *uart) write(cmd string) {
	writer := bufio.NewWriter(u.stream)
	writer.Reset(u.stream)
	cmd += "\n"
	if _, err := writer.WriteString(cmd); err != nil {
		log.Errorf("failed to write to serial port, %s", err.Error())
	}
	if err := writer.Flush(); err != nil {
		log.Errorf("failed to flush, %s", err.Error())
	}
}
