package phocus_messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	comms "github.com/wolffshots/phocus/v2/comms"
	crc "github.com/wolffshots/phocus/v2/crc"
	mqtt "github.com/wolffshots/phocus/v2/mqtt"
)

type GenericResponse struct {
	Result string
}

func SendGeneric(port comms.Port, command string, payload interface{}) (int, error) {
	switch payload.(type) {
	case int:
		written, err := port.Write(fmt.Sprintf("%s%d", command, payload))
		if err != nil {
			return -1, err
		} else {
			fmt.Printf("Wrote %s of %d bytes\n", command, written)
			return written, nil
		}
	case string:
		written, err := port.Write(fmt.Sprintf("%s%s", command, payload))
		if err != nil {
			return -1, err
		} else {
			fmt.Printf("Wrote %s of %d bytes\n", command, written)
			return written, nil
		}
	default:
		written, err := port.Write(command)
		if err != nil {
			return -1, err
		} else {
			fmt.Printf("Wrote %s of %d bytes\n", command, written)
			return written, nil
		}
	}
}

func ReceiveGeneric(port comms.Port, command string, timeout time.Duration) (string, error) {
	// read from port
	response, err := port.Read(timeout)
	log.Printf("%s\n", response)
	// verify
	if err != nil || response == "" {
		log.Printf("Failed to read with: %v\n", err)
		return "", err
	} else {
		return VerifyGeneric(response, command)
	}
}

func VerifyGeneric(response string, command string) (string, error) {
	if crc.Verify(response) {
		return response, nil
	} else {
		if len(response) < 3 {
			return "", fmt.Errorf("response not long enough: %s", response)
		}
		actual := response[len(response)-3 : len(response)-1] // 2 bytes of crc
		remainder := response[:len(response)-3]               // actual response
		wanted := crc.Checksum(remainder)                     // response calculated on response data
		message := fmt.Sprintf("invalid response from %s: CRC should have been %x but was %x", command, wanted, actual)
		log.Println(message)
		return "", errors.New(message)
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

func EncodeGeneric(response *GenericResponse) string {
	jsonGenericResponse, _ := json.Marshal(response) // err ignored because it can't fail with this input
	return string(jsonGenericResponse)
}

func PublishGeneric(client mqtt.Client, response *GenericResponse, command string) error {
	jsonResponse := EncodeGeneric(response)
	err := mqtt.Send(client, "phocus/stats/generic", 0, true, jsonResponse, 10)
	if err != nil {
		log.Printf("MQTT send of %s failed with: %v\ntype of thing sent was: %T", command, err, jsonResponse)
	} else {
		log.Printf("Sent to MQTT:\n%s\n", jsonResponse)
	}
	return err
}
