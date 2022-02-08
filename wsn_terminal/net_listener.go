package wsn_terminal

import (
	"strconv"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"
)

const (
	HEADER_CHARS   = ": "
	DATA_SEP_CHAR  = "|"
	INDEX_MSG_TYPE = 0
	INDEX_DATA     = 1
	INDEX_IP       = 2
)

type wsnNodeMsg struct {
	breed int
	data  []byte
	ip    string
}

// TODO: proto decoding
// TODO: message decoding
// TODO: Use a channel to pass data inbetween someone that needs it and the listner

func (wt *wsnTerminal) listen() {
	rawData, err := wt.serial.read()
	if err != nil {
		log.Errorf("Received an error while reading from UART, %s", err)
	}
	if rawData == "" || strings.Contains(rawData, "uart:~$ ") {
		return
	}
	data := msgCleanup(rawData)
	log.Infof("Received: %q", data)
	msg := wt.parseMsg(data)
	if len(msg.data) == 0 {
		log.Debug("empty message")
	}
	log.Infof("Parsed: %v", msg)
}

// Move these to own file I tihnk

func msgCleanup(raw string) string {
	clean := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, raw)

	for _, prefix := range []string{"[8D[J", "[21D[J"} {
		if strings.HasPrefix(clean, prefix) {
			clean = strings.Replace(clean, prefix, "", 1)
		}
	}
	return clean
}

func (wt *wsnTerminal) parseMsg(msg string) wsnNodeMsg {
	headerIndex := strings.Index(msg, HEADER_CHARS)
	if headerIndex < 0 {
		log.Warnf("Received a bad msg: %q")
		return wsnNodeMsg{}
	}
	data := msg[headerIndex+len(HEADER_CHARS):]
	msgItems := strings.Split(data, DATA_SEP_CHAR)
	if len(msgItems) <= 1 || len(msgItems) > 4 {
		log.Warnf("No WSN Node data: %q", strings.Join(msgItems, ""))
		return wsnNodeMsg{}
	}
	return wsnNodeMsg{
		breed: 0,
		data:  str2byteSlice(msgItems[INDEX_DATA]),
		ip:    msgItems[INDEX_IP+1],
	}
}

func str2byteSlice(str string) []byte {
	b := make([]byte, 0)
	items := strings.Split(strings.Trim(str, " "), " ")
	for _, item := range items {
		i, _ := strconv.Atoi(item)
		b = append(b, byte(i))
	}
	return b
}
