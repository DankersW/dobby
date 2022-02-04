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

	serialTerm, err := newUartConnection("/dev/ttyACM1")
	if err != nil {
		log.Errorf("Failed to open Serial connection to WSN gateway, %s", err.Error())
		return
	}
	serialTerm.setup()

	count := 0

	for {
		data, err := serialTerm.read()
		if err != nil {
			log.Errorf("Failed to read from serial, %s", err)
			continue
		}
		log.Infof("Data: %s", data)
		time.Sleep(1 * time.Second)

		if count == 5 {
			serialTerm.write("thread multi_light toggle")
			count = 0
		}
		count++
	}

	// FIXME: REMOVE

	config := &serial.Config{
		Name:        "/dev/ttyACM1",
		Baud:        9600,
		ReadTimeout: 1,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	//count := 0
	for {
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			log.Info(scanner.Text()) // Println will add back the final '\n'
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		if count == 3 {
			log.Info("c 20")

			writer := bufio.NewWriter(stream)
			writer.Reset(stream)
			writer.WriteString("thread multi_light toggle\n")
			writer.Flush()

			count = 0
		}
		time.Sleep(1 * time.Second) // TODO: use a go-routine for this
		count = count + 1
	}

	/*

	 */

}

// Serial: https://pkg.go.dev/github.com/tarm/serial
//https://stackoverflow.com/questions/50088669/golang-reading-from-serial
