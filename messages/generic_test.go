//go:build linux || darwin

package phocus_messages

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	crc "github.com/wolffshots/phocus/v2/crc"
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
)

func TestGeneric(t *testing.T) {
	cmd := StartCmd("socat", "-d", "-d", "PTY,link=./generic1,raw,echo=0,crnl", "PTY,link=./generic2,raw,echo=0,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(51 * time.Millisecond)

	t.Run("TestSendGeneric", func(t *testing.T) {
		// setup virtual port
		serialPort := phocus_serial.Port{
			Path:    "./generic1",
			Baud:    9600,
			Retries: 1,
		}
		commonPort1, err := serialPort.Open()
		assert.NoError(t, err)

		// valid write to virtual port
		written, err := SendGeneric(commonPort1, "GENERIC", nil)
		assert.Equal(t, 10, written)
		assert.NoError(t, err)
		written, err = SendGeneric(commonPort1, "GENERIC", 1)
		assert.Equal(t, 11, written)
		assert.NoError(t, err)
		written, err = SendGeneric(commonPort1, "GENERIC", "1")
		assert.Equal(t, 11, written)
		assert.NoError(t, err)

		err = commonPort1.Close()
		assert.NoError(t, err)

		// invalid write
		written, err = SendGeneric(commonPort1, "GENERIC", nil)
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("serial port is nil on write"), err)
		written, err = SendGeneric(commonPort1, "GENERIC", 1)
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("serial port is nil on write"), err)
		written, err = SendGeneric(commonPort1, "GENERIC", "1")
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("serial port is nil on write"), err)
	})

	t.Run("TestReceiveGeneric", func(t *testing.T) {
		// start virtual port
		time.Sleep(51 * time.Millisecond)

		// setup virtual ports
		serialPort1 := phocus_serial.Port{
			Path:    "./generic1",
			Baud:    9600,
			Retries: 1,
		}
		commonPort1, err := serialPort1.Open()
		assert.NoError(t, err)
		defer commonPort1.Close()

		serialPort2 := phocus_serial.Port{
			Path:    "./generic2",
			Baud:    9600,
			Retries: 1,
		}
		commonPort2, err := serialPort2.Open()
		assert.NoError(t, err)

		_, _ = commonPort2.Read(10 * time.Millisecond) // make sure it is empty when writing

		// valid read from virtual port
		// should time out
		response, err := ReceiveGeneric(commonPort2, "GENERIC", 0*time.Millisecond)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("read returned nothing"), err)

		err = commonPort2.Close()
		assert.NoError(t, err)

		// invalid read
		response, err = ReceiveGeneric(commonPort2, "GENERIC", 10*time.Millisecond)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("serial port is nil on read"), err)

		commonPort2, err = serialPort2.Open()
		assert.NoError(t, err)
		defer commonPort2.Close()

		// valid read from virtual port
		// should respond
		written, err := SendGeneric(commonPort1, "some message", nil) // 12
		assert.Equal(t, 15, written)                                  // 12 + 2 (crc) + 1 (cr)
		assert.NoError(t, err)
		response, err = ReceiveGeneric(commonPort2, "some message", 10*time.Millisecond)
		assert.Equal(t, "some message\xbe\x0f\r", response)
		assert.NoError(t, err)
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

		assert.Equal(t, false, crc.Verify(input[1:]))
		assert.Equal(t, true, crc.Verify(input))
		assert.Equal(t, uint16(0x2e2a), crc.Checksum(input[:len(input)-3]))

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
