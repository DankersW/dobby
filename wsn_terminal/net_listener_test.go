package wsn_terminal

import (
	"testing"

	"gotest.tools/assert"
)

func TestMsgCleanup(t *testing.T) {
	msgs := map[string]string{
		"a": "a",
		"[00:06:24.337,646] <inf> ot_coap: SensorData | 10 3 83 48 49 21 0 0 232 65 | 10 | fdde:ad00:beef:0:900a:1515:876a:92a4": "\x1b[8D\x1b[J[00:06:24.337,646] <inf> ot_coap: SensorData | 10 3 83 48 49 21 0 0 232 65 | 10 | fdde:ad00:beef:0:900a:1515:876a:92a4",
		"[01:04:43.933,959] <dbg> ot_coap.temp_monitor_set_state: Transmitted msg with return code 21":                           "\x1b[8D\x1b[J[01:04:43.933,959] <dbg> ot_coap.temp_monitor_set_state: Transmitted msg with return code 21",
	}
	for key, val := range msgs {
		expected := key
		actual := msgCleanup(val)
		assert.Equal(t, expected, actual)
	}
}
