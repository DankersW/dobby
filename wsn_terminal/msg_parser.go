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

func (wt *wsnTerminal) parseMsg(rawMsg string) wsnNodeMsg {
	msg := msgCleanup(rawMsg)
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
		breed: parseMsgBreed(msgItems[INDEX_MSG_TYPE]),
		data:  str2byteSlice(msgItems[INDEX_DATA]),
		ip:    msgItems[INDEX_IP],
	}
}

func msgCleanup(raw string) string {
	clean := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, raw)

	for _, prefix := range []string{"[8D[J", "[21D[J", "[49D[J", "[25D[J"} {
		if strings.HasPrefix(clean, prefix) {
			clean = strings.Replace(clean, prefix, "", 1)
		}
	}
	return clean
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

func parseMsgBreed(str string) int {
	val, err := strconv.Atoi(strings.Trim(str, " "))
	if err != nil {
		log.Warnf("Failed to parse message type to int value, %s", err.Error())
	}
	return val
}
