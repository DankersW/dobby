package wsn_terminal

import (
	"bufio"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type uart struct {
	stream *serial.Port
	mutex  sync.Mutex
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
	cmds := []string{"", "shell echo off", "shell colors off", "clear", "thread monitor decode off"}
	for _, cmd := range cmds {
		u.write(cmd)
		time.Sleep(1 * time.Millisecond)
	}
	time.Sleep(10 * time.Millisecond)
	log.Info("Serial setup done")
}

func (u *uart) close() error {
	return u.stream.Close()
}

func (u *uart) read() string {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	scanner := bufio.NewScanner(u.stream)
	for scanner.Scan() {
		return scanner.Text()
	}
	if scanner.Err() != nil {
		log.Errorf("Received an error while reading from UART, %s", scanner.Err())
	}
	return ""
}

func (u *uart) write(cmd string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
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
