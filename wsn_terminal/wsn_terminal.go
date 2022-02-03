package wsn_terminal

import (
	"io"

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

	for {
		var buff [100]byte
		l, err := port.Read(buff[:])
		if err == io.EOF {
			log.Error("EOF")
			break
		}
		log.Info(l)
		log.Infof("%s", buff)
	}

}

func read() {
}

func write() {
}

// Serial: https://pkg.go.dev/github.com/tarm/serial
//https://stackoverflow.com/questions/50088669/golang-reading-from-serial
