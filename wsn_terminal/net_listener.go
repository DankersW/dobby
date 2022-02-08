package wsn_terminal

import (
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"
)

const (
	HEADER_CHARS  = ": "
	DATA_SEP_CHAR = "|"
	INDEX_MSG_ID  = 0
	INDEX_DATA    = 1
	INDEX_LENGHT  = 2
	INDEX_IP      = 3
)

type wsnNodeData struct {
	msgId  int
	data   []byte
	lenght int
	ip     string
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
	wt.parseMsg(data)
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

func (wt *wsnTerminal) parseMsg(msg string) {
	headerIndex := strings.Index(msg, HEADER_CHARS)
	if headerIndex < 0 {
		log.Warnf("Received a bad msg: %q")
		return
	}
	data := msg[headerIndex+len(HEADER_CHARS):]
	msgItems := strings.Split(data, DATA_SEP_CHAR)
	if len(msgItems) <= 1 || len(msgItems) > 4 {
		log.Warnf("No WSN Node data: %q", strings.Join(msgItems, ""))
		return
	}

	log.Infof("Parsed, ID: %s \t data: %s \t lenght: %s \t ip: %s", msgItems[INDEX_MSG_ID], msgItems[INDEX_DATA], msgItems[INDEX_LENGHT], msgItems[INDEX_IP])
}

func str2byteSlice(str string) []byte {

}
