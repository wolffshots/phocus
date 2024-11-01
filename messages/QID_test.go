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

func TestQID(t *testing.T) {
	cmd := StartCmd("socat", "-d", "-d", "PTY,link=./qid1,raw,echo=0,crnl", "PTY,link=./qid2,raw,echo=0,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(51 * time.Millisecond)

	t.Run("TestSendQID", func(t *testing.T) {
		// setup virtual port
		serialPort := phocus_serial.Port{
			Path:    "./qid1",
			Baud:    9600,
			Retries: 1,
		}
		commonPort1, err := serialPort.Open()
		assert.NoError(t, err)

		// valid write to virtual port
		written, err := SendQID(commonPort1, nil)
		assert.Equal(t, 6, written)
		assert.NoError(t, err)

		commonPort1.Close()

		// invalid write
		written, err = SendQID(commonPort1, nil)
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("serial port is nil on write"), err)
	})

	t.Run("TestReceiveQID", func(t *testing.T) {
		// setup virtual port
		serialPort1 := phocus_serial.Port{
			Path:    "./qid1",
			Baud:    9600,
			Retries: 1,
		}
		commonPort1, err := serialPort1.Open()
		assert.NoError(t, err)
		defer commonPort1.Close()

		serialPort2 := phocus_serial.Port{
			Path:    "./qid2",
			Baud:    9600,
			Retries: 1,
		}
		commonPort2, err := serialPort2.Open()
		assert.NoError(t, err)

		_, _ = commonPort2.Read(10 * time.Millisecond) // make sure it is empty when writing

		// valid read from virtual port
		// should time out
		response, err := ReceiveQID(commonPort2, 0*time.Millisecond)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("read returned nothing"), err)

		err = commonPort2.Close()
		assert.NoError(t, err)

		// invalid read
		response, err = ReceiveQID(commonPort2, 10*time.Millisecond)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("serial port is nil on read"), err)

		// reopen port
		commonPort2, err = serialPort2.Open()
		assert.NoError(t, err)
		defer commonPort2.Close()

		// valid read from virtual port
		// should respond
		written, err := SendQID(commonPort1, nil) // 3
		assert.Equal(t, 6, written)               // 3 + 2 (crc) + 1 (cr)
		assert.NoError(t, err)
		response, err = ReceiveQID(commonPort2, 10*time.Millisecond)
		assert.Equal(t, "QID\xd6\xea\r", response)
		assert.NoError(t, err)
	})

	t.Run("TestVerifyQID", func(t *testing.T) {
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
	})

	t.Run("TestInterpretQID", func(t *testing.T) {
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
	})

	t.Run("TestEncodeQID", func(t *testing.T) {
		jsonResponse := EncodeQID(nil)
		assert.Equal(t, "null", jsonResponse)

		actual, err := InterpretQID("(92932004102443\x2e\x2a\r")
		assert.NoError(t, err)

		jsonResponse = EncodeQID(actual)
		assert.Equal(t, "{\"SerialNumber\":\"92932004102443\"}", jsonResponse)
	})
}
