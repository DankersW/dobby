package wsn_terminal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgCleanup(t *testing.T) {
	raw := "\x1b[8D\x1b[J[00:06:24.337,646] <inf> ot_coap: SensorData | 10 3 83 48 49 21 0 0 232 65 | 10 | fdde:ad00:beef:0:900a:1515:876a:92a4"
	expected := "[00:06:24.337,646] <inf> ot_coap: SensorData | 10 3 83 48 49 21 0 0 232 65 | 10 | fdde:ad00:beef:0:900a:1515:876a:92a4"
	actual := cleanup(raw)
	assert.Equal(t, expected, actual)
}
