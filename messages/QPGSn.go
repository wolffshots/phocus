package phocus_messages

import (
	"encoding/json" // encoding to json for mqtt
	"errors"        // creating custom err messages
	"fmt"           // string formatting
	"log"           // logging to std out
	"strings"       // string manipulation
	"time"          // sleeping

	crc "github.com/wolffshots/phocus/v2/crc"   // checksum calculations
	mqtt "github.com/wolffshots/phocus/v2/mqtt" // comms with mqtt broker
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
)

type OperationMode string

var OperationModes = map[string]OperationMode{
	"P": "Powered on",
	"S": "Stand-By",
	"L": "Grid",
	"B": "Off-grid",
	"F": "Fault",
	"D": "Shutdown",
}

type FaultCode string

var FaultCodes = map[string]FaultCode{
	"1":  "Fan locked while inverter off",
	"2":  "Over-temperature",
	"3":  "Battery voltage too high",
	"4":  "Battery voltage too low",
	"5":  "AC output short-circuit",
	"6":  "AC output voltage too high",
	"7":  "AC output overload",
	"8":  "Internal bus voltage too high",
	"9":  "Internal bus soft-start failed",
	"10": "PV over-current",
	"11": "PV over-voltage",
	"12": "Internal DC converter over-current",
	"13": "Battery discharge over-current",
	"51": "Over-current",
	"52": "Internal bus voltage too low",
	"53": "Inverter soft-start failed",
	"55": "DC over-voltage at AC output",
	"57": "Current sensor failed",
	"58": "AC Output voltage too low",
	"60": "Reverse-current protection active",
	"71": "Firmware version inconsistent",
	"72": "Current sharing fault",
	"80": "CAN communication fault",
	"81": "Host loss",
	"82": "Synchronization loss",
	"83": "Battery voltage detected inconsistent",
	"84": "AC in. voltage/frequency inconsistent",
	"85": "AC output current imbalance",
	"86": "AC output mode inconsistent",
}

type Status string

var Statuses = map[string]Status{
	"1": "on",
	"0": "off",
}

type BatteryStatus string

var BatteryStatuses = map[string]BatteryStatus{
	"03": "Battery charging and discharging disabled by battery attached to BMS port of unit",
	"02": "Battery disconnected",
	"01": "Battery voltage low",
	"00": "Battery voltage normal",
}

type GridAvailability string

var GridAvailabilities = map[string]GridAvailability{
	"1": "disconnected",
	"0": "connected",
}

type Reserved string

type InverterStatus struct {
	MPPT          Status
	ACCharging    Status
	SolarCharging Status
	BatteryStatus BatteryStatus // 2 bits
	ACInput       GridAvailability
	ACOutput      Status
	Reserved      Reserved
}

type ACOutputMode string

var ACOutputModes = map[string]ACOutputMode{
	"0": "Single Any-Grid unit",
	"1": "Parallel output",
	"2": "Phase 1 of 3-phase output",
	"3": "Phase 2 of 3-phase output",
	"4": "Phase 3 of 3-phase output",
}

type BatteryChargerSourcePriority string

var BatteryChargerSourcePriorities = map[string]BatteryChargerSourcePriority{
	"1": "Solar first",
	"2": "Solar and Utility",
	"3": "Solar only",
}

type QPGSnResponse struct {
	// (A BBBBBBBBBBBBBB C DD EEE.E FF.FF GGG.G HH.HH IIII JJJJ KKK LL.L MMM NNN OOO.O PPP QQQQQ RRRRR SSS b7b6b5b4b3b2b1b0 T U VVV WWW XX YY.Y ZZZ<CRC><cr>
	InverterNumber                      int
	OtherUnits                          bool
	SerialNumber                        string
	OperationMode                       OperationMode
	FaultCode                           FaultCode
	ACInputVoltage                      string
	ACInputFrequency                    string
	ACOutputVoltage                     string
	ACOutputFrequency                   string
	ACOutputApparentPower               string
	ACOutputActivePower                 string
	PercentageOfNominalOutputPower      string
	BatteryVoltage                      string
	BatteryChargingCurrent              string
	BatteryStateOfCharge                string
	PVInputVoltage                      string
	TotalChargingCurrent                string
	TotalACOutputApparentPower          string
	TotalACOutputActivePower            string
	TotalPercentageOfNominalOutputPower string
	InverterStatus                      InverterStatus
	ACOutputMode                        ACOutputMode
	BatteryChargerSourcePriority        BatteryChargerSourcePriority
	MaxChargingCurrentSet               string
	MaxChargingCurrentPossible          string
	MaxACChargingCurrentSet             string
	PVInputCurrent                      string
	BatteryDischargeCurrent             string
	Checksum                            string
}

func NewQPGSnResponse(input string, inverterNum int) (*QPGSnResponse, error) {
	if input == "" {
		return nil, errors.New("can't create a response from an empty string")
	}
	buffer := strings.Split(input[:len(input)-3], " ")
	buffer[0] = strings.Trim(buffer[0], "(") // strip start byte
	checksum := input[len(input)-3 : len(input)-1]
	log.Printf("Buffer: %v\n", buffer)
	log.Printf("Checksum: %x\n", checksum)
	wantedLength := 27
	if len(buffer) != wantedLength {
		return nil, fmt.Errorf("input for QPGSnResponse was %v but should have been %v", len(buffer), wantedLength)
	}

	inverterStatusBuffer := strings.Split(buffer[19], "")
	wantedLength = 8
	if len(inverterStatusBuffer) != wantedLength {
		return nil, fmt.Errorf("inverter status buffer should have been %d but was %d", wantedLength, len(inverterStatusBuffer))
	}
	return &QPGSnResponse{
		InverterNumber:                      inverterNum,
		OtherUnits:                          buffer[0] == "1" || buffer[0] == "(1",
		SerialNumber:                        buffer[1],
		OperationMode:                       OperationModes[buffer[2]],
		FaultCode:                           FaultCodes[buffer[3]],
		ACInputVoltage:                      buffer[4],
		ACInputFrequency:                    buffer[5],
		ACOutputVoltage:                     buffer[6],
		ACOutputFrequency:                   buffer[7],
		ACOutputApparentPower:               buffer[8],
		ACOutputActivePower:                 buffer[9],
		PercentageOfNominalOutputPower:      buffer[10],
		BatteryVoltage:                      buffer[11],
		BatteryChargingCurrent:              buffer[12],
		BatteryStateOfCharge:                buffer[13],
		PVInputVoltage:                      buffer[14],
		TotalChargingCurrent:                buffer[15],
		TotalACOutputApparentPower:          buffer[16],
		TotalACOutputActivePower:            buffer[17],
		TotalPercentageOfNominalOutputPower: buffer[18],
		InverterStatus: InverterStatus{
			MPPT:          Statuses[inverterStatusBuffer[0]],
			ACCharging:    Statuses[inverterStatusBuffer[1]],
			SolarCharging: Statuses[inverterStatusBuffer[2]],
			BatteryStatus: BatteryStatuses[inverterStatusBuffer[3]+inverterStatusBuffer[4]], // 2 bits
			ACInput:       GridAvailabilities[inverterStatusBuffer[5]],
			ACOutput:      Statuses[inverterStatusBuffer[6]],
			Reserved:      Reserved(inverterStatusBuffer[7]),
		},
		ACOutputMode:                 ACOutputModes[buffer[20]],
		BatteryChargerSourcePriority: BatteryChargerSourcePriorities[buffer[21]],
		MaxChargingCurrentSet:        buffer[22],
		MaxChargingCurrentPossible:   buffer[23],
		MaxACChargingCurrentSet:      buffer[24],
		PVInputCurrent:               buffer[25],
		BatteryDischargeCurrent:      buffer[26],
		Checksum:                     fmt.Sprintf("0x%x", checksum),
	}, nil

}

// HandleQPGS writes the query to the inverter and
// reads the response, deserialises it into a response
// object and sends it to MQTT
func HandleQPGS(port phocus_serial.Port, inverterNum int) (*QPGSnResponse, error) {
	query := fmt.Sprintf("QPGS%d", inverterNum)
	log.Println(query)
	bytes, err := port.Write(query)
	log.Printf("Sent %v bytes\n", bytes)
	if err != nil {
		log.Printf("Failed to write to serial with: %v\n", err)
		return nil, err
	}
	response, err := port.Read(2 * time.Second)
	if err != nil || response == "" {
		log.Printf("Failed to read from serial with: %v\n", err)
		return nil, err
	}
	valid := crc.Verify(response)
	if valid {
		QPGSResponse, err := NewQPGSnResponse(response, inverterNum)
		if err != nil || QPGSResponse == nil {
			log.Fatalf("Failed to create response with :%v", err)
		}
		jsonQPGSResponse, err := json.Marshal(QPGSResponse)
		if err != nil {
			log.Fatalf("Failed to parse response to json with :%v", err)
		}
		err = mqtt.Send(fmt.Sprintf("phocus/stats/qpgs%d", inverterNum), 0, false, string(jsonQPGSResponse), 10)
		if err != nil {
			log.Fatalf("MQTT send of %s failed with: %v\ntype of thing sent was: %T", query, err, jsonQPGSResponse)
		}
		log.Printf("Sent to MQTT:\n%s\n", jsonQPGSResponse)
		return QPGSResponse, err
	} else {
		actual := response[len(response)-3 : len(response)-1]
		remainder := response[:len(response)-3]
		wanted := crc.Checksum(remainder)
		message := fmt.Sprintf("invalid response from QPGS%d: CRC should have been %x but was %x", inverterNum, wanted, actual)
		log.Println(message)
		err = errors.New(message)
	}
	return nil, err
}
