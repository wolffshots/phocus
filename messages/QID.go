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

type QIDResponse struct {
	SerialNumber string
}

func SendQID(port comms.Port, payload interface{}) (int, error) {
	written, err := port.Write("QID")
	if err != nil {
		return -1, err
	} else {
		fmt.Printf("Wrote QID of %d bytes\n", written)
		return written, nil
	}
}

func ReceiveQID(port comms.Port, timeout time.Duration) (string, error) {
	response, err := port.Read(timeout)
	log.Printf("%s\n", response)
	if err != nil || response == "" {
		log.Printf("Failed to read with: %v\n", err)
		return "", err
	} else {
		return VerifyQID(response)
	}
}

func VerifyQID(response string) (string, error) {
	if crc.Verify(response) {
		// return
		// TODO add check for length
		log.Printf("Serial number queried: %s\n", response)
		return response, nil
	} else {
		if len(response) < 3 {
			return "", fmt.Errorf("response not long enough: %s", response)
		}
		actual := response[len(response)-3 : len(response)-1] // 2 bytes of crc
		remainder := response[:len(response)-3]               // actual response
		wanted := crc.Checksum(remainder)                     // response calculated on response data
		message := fmt.Sprintf("invalid response from QID: CRC should have been %x but was %x", wanted, actual)
		log.Println(message)
		return "", errors.New(message)
	}
}

func InterpretQID(response string) (*QIDResponse, error) {
	if response == "" {
		return nil, errors.New("can't create a response from an empty string")
	} else if len(response) < 6 {
		return nil, errors.New("response is malformed or shorter than expected")
	}
	serialNumber := strings.Trim(response[:len(response)-3], "(")
	return &QIDResponse{
		SerialNumber: serialNumber,
	}, nil
}

func EncodeQID(response *QIDResponse) string {
	jsonQIDResponse, _ := json.Marshal(response) // err ignored because it can't fail with this input
	return string(jsonQIDResponse)
}

func PublishQID(client mqtt.Client, response *QIDResponse) error {
	jsonResponse := EncodeQID(response)
	err := mqtt.Send(client, "phocus/stats/qid", 0, true, jsonResponse, 10)
	if err != nil {
		log.Printf("MQTT send of %s failed with: %v\ntype of thing sent was: %T", "QID", err, jsonResponse)
	} else {
		log.Printf("Sent to MQTT:\n%s\n", jsonResponse)
	}
	return err
}
