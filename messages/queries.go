package messages

import (
	"log"
	"strings"
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

var FaultCodes = map[int]FaultCode{
	1:  "Fan locked while inverter off",
	2:  "Over-temperature",
	3:  "Battery voltage too high",
	4:  "Battery voltage too low",
	5:  "AC output short-circuit",
	6:  "AC output voltage too high",
	7:  "AC output overload",
	8:  "Internal bus voltage too high",
	9:  "Internal bus soft-start failed",
	10: "PV over-current",
	11: "PV over-voltage",
	12: "Internal DC converter over-current",
	13: "Battery discharge over-current",
	51: "Over-current",
	52: "Internal bus voltage too low",
	53: "Inverter soft-start failed",
	55: "DC over-voltage at AC output",
	57: "Current sensor failed",
	58: "AC Output voltage too low",
	60: "Reverse-current protection active",
	71: "Firmware version inconsistent",
	72: "Current sharing fault",
	80: "CAN communication fault",
	81: "Host loss",
	82: "Synchronization loss",
	83: "Battery voltage detected inconsistent",
	84: "AC in. voltage/frequency inconsistent",
	85: "AC output current imbalance",
	86: "AC output mode inconsistent",
}

type Status string

var Statuses = map[int]Status{
	1: "on",
	0: "off",
}

type BatteryStatus string

var BatteryStatuses = map[int]BatteryStatus{
	03: "Battery charging and discharging disabled by battery attached to BMS port of unit",
	02: "Battery disconnected",
	01: "Battery voltage low",
	00: "Battery voltage normal",
}

type GridAvailability string

var GridAvailabilities = map[int]GridAvailability{
	1: "disconnected",
	0: "connected",
}

type Reserved string

type InverterStatus struct {
	MPPT          Status
	ACCharging    Status
	BatteryStatus BatteryStatus // 2 bits
	ACInput       GridAvailability
	ACOutput      Status
	Reserved      Reserved
}

type ACOutputMode string

var ACOutputModes = map[int]ACOutputMode{
	0: "Single Any-Grid unit",
	1: "Parallel output",
	2: "Phase 1 of 3-phase output",
	3: "Phase 2 of 3-phase output",
	4: "Phase 3 of 3-phase output",
}

type BatteryChargerSourcePriority string

var BatteryChargerSourcePriorities = map[int]BatteryChargerSourcePriority{
	1: "Solar first",
	2: "Solar and Utility",
	3: "Solar only",
}

type QPGSnResponse struct {
	// (A BBBBBBBBBBBBBB C DD EEE.E FF.FF GGG.G HH.HH IIII JJJJ KKK LL.L MMM NNN OOO.O PPP QQQQQ RRRRR SSS b7b6b5b4b3b2b1b0 T U VVV WWW XX YY.Y ZZZ<CRC><cr>
	// start byte
	OtherUnits                          bool
	SerialNumber                        string
	OperationMode                       OperationMode
	FaultCode                           FaultCode
	ACInputVoltage                      float64
	ACInputFrequency                    float64
	ACOutputVoltage                     float64
	ACOutputFrequency                   float64
	ACOutputApparentPower               int
	ACOutputActivePower                 int
	PercentageOfNominalOutputPower      uint16
	BatteryVoltage                      float64
	BatteryChargingCurrent              int
	BatteryStateOfCharge                uint16
	PVInputVoltage                      float64
	TotalChargingCurrent                uint16
	TotalACOutputApparentPower          int
	TotalACOutputActivePower            int
	TotalPercentageOfNominalOutputPower uint16
	InverterStatus                      InverterStatus
	ACOutputMode                        ACOutputMode
	BatteryChargerSourcePriority        BatteryChargerSourcePriority
	MaxChargingCurrentSet               int
	MaxChargingCurrentPossible          int
	MaxACChargingCurrentSet             int
	PVInputCurrent                      float64
	BatteryDischargeCurrent             int
	Checksum                            uint16
}

func NewQPGSnResponse(input string) QPGSnResponse {
	buffer := strings.Split(input, " ")
    buffer[0] = strings.Trim(buffer[0], "(") // strip start byte
	log.Printf("%v\n", buffer)
    wantedLength := 27
    if len(buffer) != wantedLength {
        log.Fatalf("QPGS string should have been %d but was %d\n", wantedLength, len(buffer))
    }
	return QPGSnResponse{
		true,
		"",
		OperationModes["P"],
		FaultCodes[5],
		0.0,
		0.0,
		0.0,
		0.0,
		5,
		5,
		6,
		0.0,
		5,
		6,
		0.0,
		5,
		5,
		5,
		5,
		InverterStatus{
			Statuses[1],
			Statuses[1],
			BatteryStatuses[01], // 2 bits
			GridAvailabilities[1],
			Statuses[1],
			"",
		},
		ACOutputModes[1],
		BatteryChargerSourcePriorities[1],
		5,
		5,
		5,
		0.0,
		5,
		5,
	}

}
