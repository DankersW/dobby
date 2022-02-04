package wsn_terminal

import (
	"bufio"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
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
	/*
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
	*/

	config := &serial.Config{
		Name:        "/dev/ttyACM0",
		Baud:        9600,
		ReadTimeout: 1,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for {
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			log.Info(scanner.Text()) // Println will add back the final '\n'
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		if count == 5 {
			log.Info("c 20")
			count = 0
		}
		time.Sleep(1 * time.Second) // TODO: use a go-routine for this
		count = count + 1
	}

	/*

	 */

}

func read() {
}

func write() {
}

// Serial: https://pkg.go.dev/github.com/tarm/serial
//https://stackoverflow.com/questions/50088669/golang-reading-from-serial
