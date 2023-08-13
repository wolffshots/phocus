//go:build linux || darwin

package phocus_messages

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	phocus_crc "github.com/wolffshots/phocus/v2/crc"
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
)

func TestSendQID(t *testing.T) {
	// invalid write
	written, err := SendQID(phocus_serial.Port{Port: nil, Path: ""}, nil)
	assert.Equal(t, -1, written)
	assert.Equal(t, errors.New("port is nil on write"), err)

	// start virtual port
	cmd := StartCmd("socat", "PTY,link=./com1,raw,echo=1,crnl", "PTY,link=./com2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(51 * time.Millisecond)

	// setup virtual port
	port1, err := phocus_serial.Setup("./com1", 2400, 1)
	assert.NoError(t, err)
	defer port1.Port.Close()

	// valid write to virtual port
	written, err = SendQID(port1, nil)
	assert.Equal(t, 6, written)
	assert.NoError(t, err)
}

func TestReceiveQID(t *testing.T) {
	// invalid read
	response, err := ReceiveQID(phocus_serial.Port{Port: nil, Path: ""}, 10*time.Millisecond)
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("port is nil on read"), err)

	// start virtual port
	cmd := StartCmd("socat", "PTY,link=./com1,raw,echo=1,crnl", "PTY,link=./com2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(51 * time.Millisecond)

	// setup virtual port
	port1, err := phocus_serial.Setup("./com1", 2400, 1)
	assert.NoError(t, err)
	defer port1.Port.Close()

	// valid read from virtual port
	// should time out
	response, err = ReceiveQID(port1, 0*time.Millisecond)
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("read returned nothing"), err)
}

func TestVerifyQID(t *testing.T) {
	// invalid length qid
	response, err := VerifyQID("")
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("response not long enough: "), err)

	// invalid length qid
	response, err = VerifyQID("\r")
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("response not long enough: \r"), err)

	// invalid length qid
	response, err = VerifyQID("1\r")
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("response not long enough: 1\r"), err)

	// invalid length qid
	response, err = VerifyQID("QI\r")
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("invalid response from QID: CRC should have been 0 but was 5149"), err)

	// invalid crc qid
	response, err = VerifyQID("(92932004102453\x2d\x2b\r")
	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("invalid response from QID: CRC should have been 1d1b but was 2d2b"), err)

	// valid crc qid
	response, err = VerifyQID("(92932004102453\x1d\x1b\r")
	assert.Equal(t, "(92932004102453\x1d\x1b\r", response)
	assert.NoError(t, err)
}

func TestInterpretQID(t *testing.T) {
	// test grabbed input
	input := "(92932004102443\x2e\x2a\r"
	want := &QIDResponse{"92932004102443"}
	actual, err := InterpretQID(input)
	assert.NoError(t, err)
	assert.Equal(t, want, actual)

	assert.Equal(t, false, phocus_crc.Verify(input[1:]))
	assert.Equal(t, true, phocus_crc.Verify(input))
	assert.Equal(t, uint16(0x947b), phocus_crc.Checksum(input[1:len(input)-3]))

	// test empty input
	input = ""
	want = (*QIDResponse)(nil)
	actual, err = InterpretQID(input)
	assert.Equal(t, errors.New("can't create a response from an empty string"), err)
	assert.Equal(t, want, actual)

	// test short input
	input = "(9\x1d\x1b\r"
	want = (*QIDResponse)(nil)
	actual, err = InterpretQID(input)
	assert.Equal(t, errors.New("response is malformed or shorter than expected"), err)
	assert.Equal(t, want, actual)
}

func TestEncodeQID(t *testing.T) {
	jsonResponse := EncodeQID(nil)
	assert.Equal(t, "null", jsonResponse)

	actual, err := InterpretQID("(92932004102443\x2e\x2a\r")
	assert.NoError(t, err)

	jsonResponse = EncodeQID(actual)
	assert.Equal(t, "{\"SerialNumber\":\"92932004102443\"}", jsonResponse)
}
