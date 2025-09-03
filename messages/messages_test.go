package phocus_messages

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
)

func Run(cmd *exec.Cmd) {
	err := cmd.Run()
	fmt.Printf("Couldn't run cmd: %v\n", err)
}

func StartCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	go Run(cmd)
	return cmd
}

func TerminateCmd(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		if runtime.GOOS == "windows" {
			cmd.Process.Signal(os.Kill)
		} else {
			cmd.Process.Signal(os.Interrupt)
		}
	} else {
		panic(fmt.Sprintf("command isn't running: %v", cmd))
	}
}

func TestMesages(t *testing.T) {
	cmd := StartCmd("socat", "-d", "-d", "PTY,link=./messages1,raw,echo=0,crnl", "PTY,link=./messages2,raw,echo=0,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(51 * time.Millisecond)

	t.Run("TestInterpretWriteErrors", func(t *testing.T) {
		var client mqtt.Client

		// setup virtual port
		serialPort2 := phocus_serial.Port{
			Path:    "./messages2",
			Baud:    9600,
			Retries: 1,
		}
		commonPort2, err := serialPort2.Open()
		assert.NoError(t, err)

		err = commonPort2.Close()
		assert.NoError(t, err)

		qpgsnresponse, err := Interpret(client, commonPort2, Message{uuid.New(), "QPGS1", ""}, 0*time.Second)
		assert.EqualError(t, err, "serial port is nil on write")
		assert.Nil(t, qpgsnresponse)

		qpgsnresponse, err = Interpret(client, commonPort2, Message{uuid.New(), "QPGS2", ""}, 0*time.Second)
		assert.EqualError(t, err, "serial port is nil on write")
		assert.Nil(t, qpgsnresponse)

		qpgsnresponse, err = Interpret(client, commonPort2, Message{uuid.New(), "QID", ""}, 0*time.Second)
		assert.EqualError(t, err, "serial port is nil on write")
		assert.Nil(t, qpgsnresponse)

		qpgsnresponse, err = Interpret(client, commonPort2, Message{uuid.New(), "SOMETHING_ELSE", ""}, 0*time.Second)
		assert.EqualError(t, err, "serial port is nil on write")
		assert.Nil(t, qpgsnresponse)
	})

	t.Run("TestInterpretReadErrors", func(t *testing.T) {
		var client mqtt.Client
		// setup virtual port
		serialPort2 := phocus_serial.Port{
			Path:    "./messages2",
			Baud:    9600,
			Retries: 1,
		}
		commonPort2, err := serialPort2.Open()
		assert.NoError(t, err)

		defer func() {
			err := commonPort2.Close()
			assert.NoError(t, err)
		}()

		qpgsnresponse, err := Interpret(client, commonPort2, Message{uuid.New(), "QPGS1", ""}, 0*time.Second)
		assert.EqualError(t, err, "read returned nothing")
		assert.Nil(t, qpgsnresponse)

		qpgsnresponse, err = Interpret(client, commonPort2, Message{uuid.New(), "QPGS2", ""}, 0*time.Second)
		assert.EqualError(t, err, "read returned nothing")
		assert.Nil(t, qpgsnresponse)

		qpgsnresponse, err = Interpret(client, commonPort2, Message{uuid.New(), "QID", ""}, 0*time.Second)
		assert.EqualError(t, err, "read returned nothing")
		assert.Nil(t, qpgsnresponse)

		qpgsnresponse, err = Interpret(client, commonPort2, Message{uuid.New(), "SOMETHING_ELSE", ""}, 0*time.Second)
		assert.EqualError(t, err, "read returned nothing")
		assert.Nil(t, qpgsnresponse)
	})

	t.Run("TestInterpret", func(t *testing.T) {
		var client mqtt.Client
		serialPort1 := phocus_serial.Port{
			Path:    "./messages1",
			Baud:    9600,
			Retries: 1,
		}
		commonPort1, err := serialPort1.Open()
		assert.NoError(t, err)

		serialPort2 := phocus_serial.Port{
			Path:    "./messages2",
			Baud:    9600,
			Retries: 1,
		}
		commonPort2, err := serialPort2.Open()
		assert.NoError(t, err)

		defer func() {
			err := commonPort1.Close()
			assert.NoError(t, err)
			err = commonPort2.Close()
			assert.NoError(t, err)
		}()

		// commonPort2.Read = func(port serial.Port, timeout time.Duration) (string, error) {
		// 	return "1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\xf2\x2d\r", nil
		// }
		written, err := commonPort1.Write("1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006") // \xf2\x2d\r
		assert.Equal(t, 134, written)
		assert.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		qpgsnresponse, err := Interpret(client, commonPort2, Message{uuid.New(), "QPGS1", ""}, 0*time.Second)
		assert.EqualError(t, err, "client not defined in send")
		assert.Equal(t, QPGSnResponse{InverterNumber: 1,
			OtherUnits:                          true,
			SerialNumber:                        "92932004102443",
			OperationMode:                       "Off-grid",
			FaultCode:                           "",
			ACInputVoltage:                      "237.0",
			ACInputFrequency:                    "50.01",
			ACOutputVoltage:                     "000.0",
			ACOutputFrequency:                   "00.00",
			ACOutputApparentPower:               "0483",
			ACOutputActivePower:                 "0387",
			PercentageOfNominalOutputPower:      "009",
			BatteryVoltage:                      "51.1",
			BatteryChargingCurrent:              "000",
			BatteryStateOfCharge:                "069",
			PVInputVoltage:                      "020.4",
			TotalChargingCurrent:                "000",
			TotalACOutputApparentPower:          "00942",
			TotalACOutputActivePower:            "00792",
			TotalPercentageOfNominalOutputPower: "007",
			InverterStatus: InverterStatus{MPPT: "off",
				ACCharging:    "off",
				SolarCharging: "off",
				BatteryStatus: "Battery voltage normal",
				ACInput:       "connected",
				ACOutput:      "on",
				Reserved:      "0"},
			ACOutputMode:                 "Parallel output",
			BatteryChargerSourcePriority: "Solar first",
			MaxChargingCurrentSet:        "060",
			MaxChargingCurrentPossible:   "080",
			MaxACChargingCurrentSet:      "10",
			PVInputCurrent:               "00.0",
			BatteryDischargeCurrent:      "006",
			Checksum:                     "0xf22d"},
			*qpgsnresponse,
		)

		// commonPort2.Read = func(port serial.Port, timeout time.Duration) (string, error) {
		// 	return "1 92932004102453 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\x9f\x50\r", nil
		// }
		written, err = commonPort1.Write("1 92932004102453 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006") // \xf2\x2d\r
		assert.Equal(t, 134, written)
		assert.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		qpgsnresponse, err = Interpret(client, commonPort2, Message{uuid.New(), "QPGS2", ""}, 0*time.Second)
		assert.EqualError(t, err, "client not defined in send")
		assert.Equal(t, QPGSnResponse{InverterNumber: 2,
			OtherUnits:                          true,
			SerialNumber:                        "92932004102453",
			OperationMode:                       "Off-grid",
			FaultCode:                           "",
			ACInputVoltage:                      "237.0",
			ACInputFrequency:                    "50.01",
			ACOutputVoltage:                     "000.0",
			ACOutputFrequency:                   "00.00",
			ACOutputApparentPower:               "0483",
			ACOutputActivePower:                 "0387",
			PercentageOfNominalOutputPower:      "009",
			BatteryVoltage:                      "51.1",
			BatteryChargingCurrent:              "000",
			BatteryStateOfCharge:                "069",
			PVInputVoltage:                      "020.4",
			TotalChargingCurrent:                "000",
			TotalACOutputApparentPower:          "00942",
			TotalACOutputActivePower:            "00792",
			TotalPercentageOfNominalOutputPower: "007",
			InverterStatus: InverterStatus{MPPT: "off",
				ACCharging:    "off",
				SolarCharging: "off",
				BatteryStatus: "Battery voltage normal",
				ACInput:       "connected",
				ACOutput:      "on",
				Reserved:      "0"},
			ACOutputMode:                 "Parallel output",
			BatteryChargerSourcePriority: "Solar first",
			MaxChargingCurrentSet:        "060",
			MaxChargingCurrentPossible:   "080",
			MaxACChargingCurrentSet:      "10",
			PVInputCurrent:               "00.0",
			BatteryDischargeCurrent:      "006",
			Checksum:                     "0x9f50"},
			*qpgsnresponse,
		)

		// commonPort2.Read = func(port serial.Port, timeout time.Duration) (string, error) {
		// 	return "92932004102453\xa7\x4a\r", nil
		// }
		written, err = commonPort1.Write("92932004102453") // \xa7\x4a\r
		assert.Equal(t, 17, written)
		assert.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		qpgsnresponse, err = Interpret(client, commonPort2, Message{uuid.New(), "QID", ""}, 0*time.Second)
		assert.EqualError(t, err, "client not defined in send")
		assert.Nil(t, qpgsnresponse)

		// commonPort2.Read = func(port serial.Port, timeout time.Duration) (string, error) {
		// 	return "SOME_RESPONSE\xb2\xb2\r", nil
		// }
		written, err = commonPort1.Write("SOME_RESPONSE") // \xb2\xb2\r
		assert.Equal(t, 16, written)
		assert.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		qpgsnresponse, err = Interpret(client, commonPort2, Message{uuid.New(), "SOME_MESSAGE", ""}, 0*time.Second)
		assert.EqualError(t, err, "client not defined in send")
		assert.Nil(t, qpgsnresponse)
	})
}
