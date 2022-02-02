package wsnterminal

import (
	"github.com/pkg/term"
	log "github.com/sirupsen/logrus"
)

type wsnTerminal struct {
	port string
	// TODO: handlers to call on a certain type of message
}

type WsnTerminal interface {
	Start()
}

func New(port string) WsnTerminal {
	wt := &wsnTerminal{
		port: port,
	}
	return wt
}

func (wt *wsnTerminal) Start() {
	log.Info("starting")
	port, err := term.Open("/dev/ttyACM0")
	log.Info(err)
	log.Info(port)
	var buff [100]byte
	l, err := port.Read(buff[:])
	log.Info(err)
	log.Info(l)
	log.Infof("%s", buff)
}

func read() {
}

func write() {
}

// USe go usb: https://pkg.go.dev/github.com/karalabe/usb#section-documentation
