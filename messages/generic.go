package phocus_messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	phocus_crc "github.com/wolffshots/phocus/v2/crc"
	phocus_mqtt "github.com/wolffshots/phocus/v2/mqtt"
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
)

type GenericResponse struct {
	Result string
}

func SendGeneric(port phocus_serial.Port, command string, payload interface{}) (int, error) {
	switch payload.(type) {
	case int:
		written, err := port.Write(phocus_crc.Encode(fmt.Sprintf("%s%d", command, payload)))
		if err != nil {
			return -1, err
		} else {
			return written, nil
		}
	case string:
		written, err := port.Write(phocus_crc.Encode(fmt.Sprintf("%s%s", command, payload)))
		if err != nil {
			return -1, err
		} else {
			return written, nil
		}
	default:
		return -1, errors.New("qpgsn does not support string payloads")
	}
}

func ReceiveGeneric(port phocus_serial.Port, command string, timeout time.Duration) (string, error) {
	// read from port
	response, err := port.Read(timeout)
	log.Println(response)
	// verify
	if err != nil || response == "" {
		log.Printf("Failed to read from serial with: %v\n", err)
		return "", err
	} else {
		if phocus_crc.Verify(response) {
			// return
			// TODO check for success
			return response, nil
		} else {
			return "", fmt.Errorf("invalid CRC for %s", command)
		}
	}
}

func InterpretGeneric(response string) (*GenericResponse, error) {
	if response == "" {
		return nil, errors.New("can't create a response from an empty string")
	}
	result := strings.Trim(response[:len(response)-3], "(")
	return &GenericResponse{
		Result: result,
	}, nil
}

func PublishGeneric(response *GenericResponse, command string) error {
	jsonGenericResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to parse response to json with :%v", err)
	}
	err = phocus_mqtt.Send("phocus/stats/generic", 0, false, string(jsonGenericResponse), 10)
	if err != nil {
		log.Fatalf("MQTT send of %s failed with: %v\ntype of thing sent was: %T", command, err, jsonGenericResponse)
	}
	log.Printf("Sent to MQTT:\n%s\n", jsonGenericResponse)
	return err
}
