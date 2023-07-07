package phocus_messages

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	phocus_crc "github.com/wolffshots/phocus/v2/crc"
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
)

func TestSendGeneric(t *testing.T) {
	// invalid write
	written, err := SendGeneric(phocus_serial.Port{Port: nil, Path: ""}, "GENERIC", nil)
	assert.Equal(t, -1, written)
	assert.Equal(t, errors.New("port is nil on write"), err)
	written, err = SendGeneric(phocus_serial.Port{Port: nil, Path: ""}, "GENERIC", 1)
	assert.Equal(t, -1, written)
	assert.Equal(t, errors.New("port is nil on write"), err)
	written, err = SendGeneric(phocus_serial.Port{Port: nil, Path: ""}, "GENERIC", "1")
	assert.Equal(t, -1, written)
	assert.Equal(t, errors.New("port is nil on write"), err)

	// start virtual port
	cmd := StartCmd("socat", "PTY,link=./com1,raw,echo=1,crnl", "PTY,link=./com2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(200 * time.Millisecond)

	// setup virtual port
	port1, err := phocus_serial.Setup("./com1", 2400)
	defer port1.Port.Close()
	assert.NoError(t, err)

	// valid write to virtual port
	written, err = SendGeneric(port1, "GENERIC", nil)
	assert.Equal(t, 10, written)
	assert.NoError(t, err)
	written, err = SendGeneric(port1, "GENERIC", 1)
	assert.Equal(t, 11, written)
	assert.NoError(t, err)
	written, err = SendGeneric(port1, "GENERIC", "1")
	assert.Equal(t, 11, written)
	assert.NoError(t, err)
}

func TestReceiveGeneric(t *testing.T) {
	// invalid read
	response, err := ReceiveGeneric(phocus_serial.Port{Port: nil, Path: ""}, "GENERIC", 10*time.Millisecond)
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("port is nil on read"), err)

	// start virtual port
	cmd := StartCmd("socat", "PTY,link=./com1,raw,echo=1,crnl", "PTY,link=./com2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(200 * time.Millisecond)

	// setup virtual port
	port1, err := phocus_serial.Setup("./com1", 2400)
	defer port1.Port.Close()
	assert.NoError(t, err)

	// valid read from virtual port
	// should time out
	response, err = ReceiveGeneric(port1, "GENERIC", 0*time.Millisecond)
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("read timed out"), err)
}

func TestVerifyGeneric(t *testing.T) {
	// invalid length generic
	response, err := VerifyGeneric("", "GENERIC")
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("response not long enough: "), err)

	// invalid length generic
	response, err = VerifyGeneric("\r", "GENERIC")
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("response not long enough: \r"), err)

	// invalid length generic
	response, err = VerifyGeneric("1\r", "GENERIC")
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("response not long enough: 1\r"), err)

	// invalid length generic
	response, err = VerifyGeneric("QI\r", "GENERIC")
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("invalid response from GENERIC: CRC should have been 0 but was 5149"), err)

	// invalid crc Generic
	response, err = VerifyGeneric("(92932004102453\x2d\x2b\r", "GENERIC")
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("invalid response from GENERIC: CRC should have been 1d1b but was 2d2b"), err)

	// valid crc Generic
	response, err = VerifyGeneric("(92932004102453\x1d\x1b\r", "GENERIC")
	assert.Equal(t, "(92932004102453\x1d\x1b\r", response)
	assert.NoError(t, err)
}

func TestInterpretGeneric(t *testing.T) {
	// test grabbed input
	input := "(92932004102443\x2e\x2a\r"
	want := &GenericResponse{"92932004102443"}
	actual, err := InterpretGeneric(input)
	assert.NoError(t, err)
	assert.Equal(t, want, actual)

	assert.Equal(t, false, phocus_crc.Verify(input[1:]))
	assert.Equal(t, true, phocus_crc.Verify(input))
	assert.Equal(t, uint16(0x2e2a), phocus_crc.Checksum(input[:len(input)-3]))

	// test grabbed input
	input = "(ACK\x94\x7b\r"
	want = &GenericResponse{"ACK"}
	actual, err = InterpretGeneric(input)
	assert.NoError(t, err)
	assert.Equal(t, want, actual)

	// test grabbed input
	input = "(NAK\x94\x7b\r"
	want = &GenericResponse{"NAK"}
	actual, err = InterpretGeneric(input)
	assert.NoError(t, err)
	assert.Equal(t, want, actual)

	// test empty input
	input = ""
	want = (*GenericResponse)(nil)
	actual, err = InterpretGeneric(input)
	assert.Equal(t, errors.New("can't create a response from an empty string"), err)
	assert.Equal(t, want, actual)
}

func TestEncodeGeneric(t *testing.T) {
	jsonResponse := EncodeGeneric(nil)
	assert.Equal(t, "null", jsonResponse)

	actual, err := InterpretGeneric("(ACK\x2e\x2a\r")
	assert.NoError(t, err)

	jsonResponse = EncodeGeneric(actual)
	assert.Equal(t, "{\"Result\":\"ACK\"}", jsonResponse)

	actual, err = InterpretGeneric("(NAK\x2e\x2a\r")
	assert.NoError(t, err)

	jsonResponse = EncodeGeneric(actual)
	assert.Equal(t, "{\"Result\":\"NAK\"}", jsonResponse)
}
