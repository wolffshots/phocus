//go:build linux || darwin

package phocus_messages

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	phocus_crc "github.com/wolffshots/phocus/v2/crc"
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
	"go.bug.st/serial"
)

func TestQPGSn(t *testing.T) {
	cmd := StartCmd("socat", "PTY,link=./qpgsn1,raw,echo=1,crnl", "PTY,link=./qpgsn2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(10 * time.Millisecond)

	t.Run("TestSendQPGSn", func(t *testing.T) {
		// setup virtual port
		port1, err := phocus_serial.Setup("./qpgsn1", 2400, 1)
		assert.NoError(t, err)

		// valid write to virtual port
		written, err := SendQPGSn(port1, nil)
		assert.Equal(t, 8, written)
		assert.NoError(t, err)

		// valid write to virtual port
		written, err = SendQPGSn(port1, 1)
		assert.Equal(t, 8, written)
		assert.NoError(t, err)

		// invalid write to virtual port with string payload
		written, err = SendQPGSn(port1, "1")
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("qpgsn does not support string payloads"), err)

		port1.Port.Close()
		port1.Port = nil

		// invalid write
		written, err = SendQPGSn(port1, nil)
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("port is nil on write"), err)

		// invalid write
		written, err = SendQPGSn(port1, "1")
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("qpgsn does not support string payloads"), err)

		// invalid write
		written, err = SendQPGSn(port1, 1)
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("port is nil on write"), err)
	})

	t.Run("TestReceiveQPGSn", func(t *testing.T) {
		// setup virtual port
		port1, err := phocus_serial.Setup("./qpgsn1", 2400, 1)
		assert.NoError(t, err)

		// valid read from virtual port
		// should time out
		response, err := ReceiveQPGSn(port1, 0*time.Millisecond, 0)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("read returned nothing"), err)

		port1.Port.Close()
		port1.Port = nil

		// invalid read
		response, err = ReceiveQPGSn(port1, 10*time.Millisecond, 1)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("port is nil on read"), err)

		port1.Read = func(port serial.Port, timeout time.Duration) (string, error) {
			return "some response\xea\xac\r", nil
		}

		// valid read from virtual port
		// should respond
		response, err = ReceiveQPGSn(port1, 10*time.Millisecond, 0)
		assert.Equal(t, "some response\xea\xac\r", response)
		assert.NoError(t, err)

		port1.Read = func(port serial.Port, timeout time.Duration) (string, error) {
			return "", errors.New("some error")
		}

		// valid read from virtual port
		// should respond with err
		response, err = ReceiveQPGSn(port1, 0*time.Millisecond, 0)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("some error"), err)
	})

	t.Run("TestVerifyQPGSn", func(t *testing.T) {
		// invalid length QPGSn
		response, err := VerifyQPGSn("", 1)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("response not long enough: "), err)

		// invalid length QPGSn
		response, err = VerifyQPGSn("\r", 0)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("response not long enough: \r"), err)

		// invalid length QPGSn
		response, err = VerifyQPGSn("1\r", 2)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("response not long enough: 1\r"), err)

		// invalid length QPGSn
		response, err = VerifyQPGSn("QI\r", 2)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("invalid response from QPGS2: CRC should have been 0 but was 5149"), err)

		// invalid crc QPGSn
		response, err = VerifyQPGSn("(92932004102453\x2d\x2b\r", 1)
		assert.Equal(t, "", response)
		assert.Equal(t, errors.New("invalid response from QPGS1: CRC should have been 1d1b but was 2d2b"), err)

		// valid crc QPGSn
		response, err = VerifyQPGSn("(92932004102453\x1d\x1b\r", 1)
		assert.Equal(t, "(92932004102453\x1d\x1b\r", response)
		assert.NoError(t, err)
	})

	t.Run("TestInterpretQPGSn", func(t *testing.T) {
		// test grabbed input
		input := "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\x06\x6e\r"
		want := &QPGSnResponse{5, true, "92932004102443", "Off-grid", "", "237.0", "50.01", "000.0", "00.00", "0483", "0387", "009", "51.1", "000", "069", "020.4", "000", "00942", "00792", "007", InverterStatus{"off", "off", "off", "Battery voltage normal", "connected", "on", "0"}, "Parallel output", "Solar first", "060", "080", "10", "00.0", "006", fmt.Sprintf("0x%02x%02x", 0x06, 0x6e)}
		actual, err := InterpretQPGSn(input, 5)
		assert.NoError(t, err)
		assert.Equal(t, want, actual)

		assert.Equal(t, false, phocus_crc.Verify(input[1:]))
		assert.Equal(t, true, phocus_crc.Verify(input))
		assert.Equal(t, uint16(0xf22d), phocus_crc.Checksum(input[1:len(input)-3]))

		// test empty input
		input = ""
		want = (*QPGSnResponse)(nil)
		actual, err = InterpretQPGSn(input, 0)
		assert.Equal(t, errors.New("can't create a response from an empty string"), err)
		assert.Equal(t, want, actual)

		input = "(1 929320483 0387 942 00792 007 00000010 1 1 060 080 10 00.0 006\x06\x6e\r"
		actual, err = InterpretQPGSn(input, 5)
		assert.Equal(t, (*QPGSnResponse)(nil), actual)
		assert.Equal(t, errors.New("input for QPGSnResponse was 14 but should have been 27"), err)

		input = "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00010 1 1 060 080 10 00.0 006\x06\x6e\r"
		actual, err = InterpretQPGSn(input, 5)
		assert.Equal(t, (*QPGSnResponse)(nil), actual)
		assert.Equal(t, errors.New("inverter status buffer should have been 8 but was 5"), err)
	})

	t.Run("TestEncodeQPGSn", func(t *testing.T) {
		jsonResponse := EncodeQPGSn(nil)
		assert.Equal(t, "null", jsonResponse)

		input := "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\x06\x6e\r"
		actual, err := InterpretQPGSn(input, 1)
		assert.NoError(t, err)

		want := "{\"InverterNumber\":1,\"OtherUnits\":true,\"SerialNumber\":\"92932004102443\",\"OperationMode\":\"Off-grid\",\"FaultCode\":\"\",\"ACInputVoltage\":\"237.0\",\"ACInputFrequency\":\"50.01\",\"ACOutputVoltage\":\"000.0\",\"ACOutputFrequency\":\"00.00\",\"ACOutputApparentPower\":\"0483\",\"ACOutputActivePower\":\"0387\",\"PercentageOfNominalOutputPower\":\"009\",\"BatteryVoltage\":\"51.1\",\"BatteryChargingCurrent\":\"000\",\"BatteryStateOfCharge\":\"069\",\"PVInputVoltage\":\"020.4\",\"TotalChargingCurrent\":\"000\",\"TotalACOutputApparentPower\":\"00942\",\"TotalACOutputActivePower\":\"00792\",\"TotalPercentageOfNominalOutputPower\":\"007\",\"InverterStatus\":{\"MPPT\":\"off\",\"ACCharging\":\"off\",\"SolarCharging\":\"off\",\"BatteryStatus\":\"Battery voltage normal\",\"ACInput\":\"connected\",\"ACOutput\":\"on\",\"Reserved\":\"0\"},\"ACOutputMode\":\"Parallel output\",\"BatteryChargerSourcePriority\":\"Solar first\",\"MaxChargingCurrentSet\":\"060\",\"MaxChargingCurrentPossible\":\"080\",\"MaxACChargingCurrentSet\":\"10\",\"PVInputCurrent\":\"00.0\",\"BatteryDischargeCurrent\":\"006\",\"Checksum\":\"0x066e\"}"
		jsonResponse = EncodeQPGSn(actual)
		assert.Equal(t, want, jsonResponse)
	})
}
