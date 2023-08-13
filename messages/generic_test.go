//go:build linux || darwin

package phocus_messages

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	phocus_crc "github.com/wolffshots/phocus/v2/crc"
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
	"go.bug.st/serial"
)

func TestGeneric(t *testing.T) {
	cmd := StartCmd("socat", "PTY,link=./generic1,raw,echo=1,crnl", "PTY,link=./generic2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(10 * time.Millisecond)

	t.Run("TestSendGeneric", func(t *testing.T) {
		// start virtual port
		time.Sleep(51 * time.Millisecond)

		// setup virtual port
		port1, err := phocus_serial.Setup("./generic1", 2400, 1)
		assert.NoError(t, err)

		// valid write to virtual port
		written, err := SendGeneric(port1, "GENERIC", nil)
		assert.Equal(t, 10, written)
		assert.NoError(t, err)
		written, err = SendGeneric(port1, "GENERIC", 1)
		assert.Equal(t, 11, written)
		assert.NoError(t, err)
		written, err = SendGeneric(port1, "GENERIC", "1")
		assert.Equal(t, 11, written)
		assert.NoError(t, err)

		port1.Port.Close()
		port1.Port = nil

		// invalid write
		written, err = SendGeneric(port1, "GENERIC", nil)
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("port is nil on write"), err)
		written, err = SendGeneric(port1, "GENERIC", 1)
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("port is nil on write"), err)
		written, err = SendGeneric(port1, "GENERIC", "1")
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("port is nil on write"), err)
	})

	t.Run("TestReceiveGeneric", func(t *testing.T) {
		// start virtual port
		time.Sleep(51 * time.Millisecond)

		// setup virtual port
		port1, err := phocus_serial.Setup("./generic1", 2400, 1)
		assert.NoError(t, err)

		// valid read from virtual port
		// should time out
		response, err := ReceiveGeneric(port1, "GENERIC", 0*time.Millisecond)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("read returned nothing"), err)

		port1.Port.Close()
		port1.Port = nil

		// invalid read
		response, err = ReceiveGeneric(port1, "GENERIC", 10*time.Millisecond)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("port is nil on read"), err)

		port1.Read = func(port serial.Port, timeout time.Duration) (string, error) {
			return "some response\xea\xac\r", nil
		}

		// valid read from virtual port
		// should respond
		response, err = ReceiveGeneric(port1, "some message", 10*time.Millisecond)
		assert.Equal(t, "some response\xea\xac\r", response)
		assert.NoError(t, err)

		port1.Read = func(port serial.Port, timeout time.Duration) (string, error) {
			return "", errors.New("some error")
		}

		// valid read from virtual port
		// should respond with err
		response, err = ReceiveGeneric(port1, "some message", 0*time.Millisecond)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("some error"), err)
	})

	t.Run("TestVerifyGeneric", func(t *testing.T) {
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
	})

	t.Run("TestInterpretGeneric", func(t *testing.T) {
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
	})

	t.Run("TestEncodeGeneric", func(t *testing.T) {
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
	})
}
